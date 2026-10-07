// Package tree provides an immutable snapshot of a project directory.
//
// One WalkDir builds the whole index; every rule then matches against memory
// instead of issuing its own stat or glob syscall.
package tree

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Entry is a single indexed path.
type Entry struct {
	Size  int64
	IsDir bool
}

// Snapshot is an immutable index of a project tree, keyed by slash-separated
// paths relative to the root.
type Snapshot struct {
	root    string
	entries map[string]Entry
	// children counts entries per directory, used for the non-empty test.
	children map[string]int
}

func (s *Snapshot) Root() string { return s.root }

func (s *Snapshot) HasDir(rel string) bool {
	e, ok := s.entries[clean(rel)]
	return ok && e.IsDir
}

func (s *Snapshot) HasFile(rel string) bool {
	e, ok := s.entries[clean(rel)]
	return ok && !e.IsDir && e.Size > 0
}

func (s *Snapshot) Exists(rel string) bool {
	_, ok := s.entries[clean(rel)]
	return ok
}

// NonEmpty reports whether rel is a directory holding at least one entry.
func (s *Snapshot) NonEmpty(rel string) bool {
	key := clean(rel)
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	if !e.IsDir {
		return e.Size > 0
	}
	return s.children[key] > 0
}

func (s *Snapshot) Size(rel string) int64 {
	return s.entries[clean(rel)].Size
}

// ReadFile returns the contents of rel.
func (s *Snapshot) ReadFile(rel string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.root, filepath.FromSlash(clean(rel))))
}

func (s *Snapshot) Paths() []string {
	out := make([]string, 0, len(s.entries))
	for p := range s.entries {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func (s *Snapshot) Match(pattern string) bool {
	if s == nil {
		return false
	}
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}

	// Directories are stored without a trailing separator.
	if p := strings.TrimSuffix(pattern, "/"); !strings.ContainsAny(p, "*?[") {
		if e, ok := s.entries[p]; ok && (p != "" || e.IsDir) {
			// A trailing slash demands a directory.
			return !strings.HasSuffix(pattern, "/") || e.IsDir
		}
		return false
	}
	return len(s.Glob(pattern)) > 0
}

func (s *Snapshot) Glob(pattern string) []string {
	if s == nil || pattern == "" {
		return nil
	}
	// doublestar needs the pattern to be slash-separated.
	pattern = filepath.ToSlash(pattern)

	var (
		out  []string
		seen = make(map[string]struct{})
	)
	for _, p := range s.candidates(pattern) {
		if _, dup := seen[p]; dup {
			continue
		}
		if ok, err := doublestar.Match(pattern, p); err == nil && ok {
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

func (s *Snapshot) candidates(pattern string) []string {
	prefix := literalPrefix(pattern)
	if prefix == "" {
		return s.Paths()
	}

	out := make([]string, 0, 64)
	// A dir path prefix also matches everything beneath it.
	if e, ok := s.entries[prefix]; ok && !e.IsDir {
		out = append(out, prefix)
	}
	for p, e := range s.entries {
		if e.IsDir {
			continue
		}
		if p == prefix || strings.HasPrefix(p, prefix+"/") {
			out = append(out, p)
		}
	}
	return out
}

func literalPrefix(pattern string) string {
	segs := strings.Split(pattern, "/")
	var fixed []string
	for _, s := range segs {
		if strings.ContainsAny(s, "*?[") {
			break
		}
		fixed = append(fixed, s)
	}
	if len(fixed) == len(segs) {
		return strings.Join(fixed, "/")
	}
	return strings.Join(fixed, "/")
}

func clean(p string) string {
	p = filepath.ToSlash(p)
	p = strings.TrimPrefix(p, "./")
	return strings.TrimSuffix(p, "/")
}

func Scan(root string, ignore []string) (*Snapshot, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	m := newMatcher(ignore)
	snap := &Snapshot{
		root:     abs,
		entries:  make(map[string]Entry, 256),
		children: make(map[string]int, 256),
	}

	// "." is the root itself, so a pattern like "*/x" can address top-level entries.
	snap.entries["."] = Entry{IsDir: true}

	err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Unreadable entries are skipped rather than aborting the whole scan.
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		rel, rerr := filepath.Rel(abs, path)
		if rerr != nil {
			return nil
		}
		rel = clean(rel)
		if rel == "" {
			return nil
		}

		isDir := d.IsDir()
		if m.match(rel, isDir) {
			if isDir {
				return fs.SkipDir
			}
			return nil
		}

		if isDir {
			snap.entries[rel] = Entry{IsDir: true}
			return nil
		}

		info, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		snap.entries[rel] = Entry{Size: info.Size()}
		snap.bump(rel)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// bump increments the child counter of every ancestor directory of rel.
func (s *Snapshot) bump(rel string) {
	for i := len(rel) - 1; i > 0; i-- {
		if rel[i] != '/' {
			continue
		}
		s.children[rel[:i]]++
	}
	s.children["."]++
}
