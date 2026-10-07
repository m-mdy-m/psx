package tree

import (
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// ignoreRule is one gitignore-style pattern.
type ignoreRule struct {
	pattern  string
	dirOnly  bool
	anchored bool
	negate   bool
}

type matcher struct {
	rules []ignoreRule
}

// newMatcher compiles ignore patterns into gitignore semantics:
// a trailing slash limits the rule to directories, a leading slash anchors it
// to the root, and a leading "!" re-includes a previously ignored path.
func newMatcher(patterns []string) *matcher {
	m := &matcher{}
	for _, raw := range patterns {
		p := strings.TrimSpace(strings.ReplaceAll(raw, "\\", "/"))
		if p == "" || strings.HasPrefix(p, "#") {
			continue
		}

		r := ignoreRule{}
		if strings.HasPrefix(p, "!") {
			r.negate = true
			p = p[1:]
		}
		if strings.HasSuffix(p, "/") {
			r.dirOnly = true
			p = strings.TrimSuffix(p, "/")
		}
		if strings.HasPrefix(p, "/") {
			r.anchored = true
			p = strings.TrimPrefix(p, "/")
		}
		if p == "" {
			continue
		}
		// A pattern with a slash in the middle is already anchored to the root.
		r.anchored = r.anchored || strings.Contains(p, "/")
		r.pattern = p
		m.rules = append(m.rules, r)
	}
	return m
}

func (m *matcher) match(rel string, isDir bool) bool {
	if m == nil || len(m.rules) == 0 {
		return false
	}

	ignored := false
	for _, r := range m.rules {
		if r.dirOnly && !isDir {
			continue
		}
		if r.matches(rel) {
			ignored = !r.negate
		}
	}
	return ignored
}

func (r ignoreRule) matches(rel string) bool {
	if r.anchored {
		return globMatch(r.pattern, rel)
	}
	// Unanchored rules match any path segment suffix, e.g. "dist" hits "a/b/dist".
	if globMatch(r.pattern, rel) {
		return true
	}
	for i := 0; i < len(rel); i++ {
		if rel[i] == '/' && globMatch(r.pattern, rel[i+1:]) {
			return true
		}
	}
	return false
}

// globMatch reports a doublestar match, treating a literal directory name as a
// prefix so "dist" also covers "dist/a.js". A glob tail such as "logs/*" does
// not match the directory "logs" itself, matching gitignore behaviour.
func globMatch(pattern, name string) bool {
	if ok, err := doublestar.Match(pattern, name); err == nil && ok {
		return true
	}
	if strings.ContainsAny(pattern, "*?[") {
		return false
	}
	base := strings.TrimSuffix(pattern, "/")
	return base != "" && strings.HasPrefix(name, base+"/")
}

func IgnoreMatcher(patterns []string) func(rel string, isDir bool) bool {
	m := newMatcher(patterns)
	return m.match
}

func NormalizePattern(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimPrefix(p, "./")
	if p == "" {
		return "."
	}
	return path.Clean(p)
}
