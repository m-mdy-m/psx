// Package watch re-runs a check as the project changes.
//
// It polls rather than using filesystem notifications, which keeps the
// implementation dependency-free and behaves identically on Linux, macOS,
// Windows, WSL and network shares. Polling cost is proportional to file count,
// not to the number of active rules.
package watch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/tree"
)

type State struct {
	// Fingerprint hashes path, size and modification time for every file.
	Fingerprint string
	// Count is the number of files observed.
	Count int
	// ConfigHash fingerprints the configuration file, so an edited config can
	// be told apart from an edited source file.
	ConfigHash string
}

// Scan builds a fingerprint for the project tree.
func Scan(root, configFile string, ignore []string) (State, error) {
	var sb strings.Builder
	count := 0

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if d.IsDir() {
			name := d.Name()
			if path != root && (name == ".git" || name == "node_modules" || name == "vendor") {
				return filepath.SkipDir
			}
			return nil
		}
		info, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		sb.WriteString(filepath.ToSlash(rel))
		sb.WriteByte(':')
		sb.WriteString(fmt.Sprint(info.Size()))
		sb.WriteByte(':')
		sb.WriteString(fmt.Sprint(info.ModTime().UnixNano()))
		sb.WriteByte('\n')
		count++
		return nil
	})
	if err != nil {
		return State{}, err
	}

	return State{
		Fingerprint: hash(sb.String()),
		Count:       count,
		ConfigHash:  fileHash(configFile),
	}, nil
}

// Changed reports whether the project or its configuration differs.
func (s State) Changed(prev State) bool {
	return s.Fingerprint != prev.Fingerprint || s.ConfigHash != prev.ConfigHash
}

// ConfigChanged reports whether only the configuration differs.
func (s State) ConfigChanged(prev State) bool {
	return s.ConfigHash != prev.ConfigHash
}

// Run watches the project until ctx is cancelled or stop is closed.
//
// runCheck is invoked on every change; onChange receives the new and previous
// results so the caller can render a diff rather than a full report.
func Run(ctx context.Context, opts Options, runCheck func() (*rules.ExecutionResult, error), onChange func(prev, cur *rules.ExecutionResult)) error {
	root := opts.Root

	current, err := Scan(root, opts.ConfigFile, opts.Ignore)
	if err != nil {
		return fmt.Errorf("initial scan: %w", err)
	}
	logger.Verbosef("watching %s (%d files)", root, current.Count)

	prevResult, err := runCheck()
	if err != nil {
		return err
	}
	if onChange != nil {
		onChange(nil, prevResult)
	}
	if opts.OnceClean && isClean(prevResult) {
		logger.Success("No problems found")
		return nil
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	debounce := time.NewTimer(opts.Debounce)
	defer debounce.Stop()
	if !debounce.Stop() {
		<-debounce.C
	}

	var pending bool
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-debounce.C:
			// Coalesce bursts of writes into a single re-check.
			pending = false
			next, err := Scan(root, opts.ConfigFile, opts.Ignore)
			if err != nil {
				logger.Verbosef("rescan failed: %v", err)
				continue
			}
			if !next.Changed(current) {
				continue
			}

			current = next
			result, err := runCheck()
			if err != nil {
				logger.Errorf("check failed: %v", err)
				continue
			}

			if onChange != nil {
				onChange(prevResult, result)
			}
			prevResult = result

			if opts.OnceClean && isClean(result) {
				logger.Success("No problems found")
				return nil
			}

		case <-ticker.C:
			if !pending {
				pending = true
				debounce.Reset(opts.Debounce)
			}
		}
	}
}

// Options configures a watch run.
type Options struct {
	Root       string
	ConfigFile string
	Ignore     []string
	Interval   time.Duration
	Debounce   time.Duration
	OnceClean  bool
}

func Defaults() Options {
	return Options{
		Interval: 1500 * time.Millisecond,
		Debounce: 250 * time.Millisecond,
	}
}

func isClean(res *rules.ExecutionResult) bool {
	return res != nil && res.Summary.Failed == 0
}

// Diff describes how a rule's status changed between two runs.
type Diff struct {
	RuleID   string
	From     rules.Status
	To       rules.Status
	Severity config.Severity
	Message  string
}

func Compare(prev, cur *rules.ExecutionResult) []Diff {
	before := map[string]rules.RuleResult{}
	if prev != nil {
		for _, r := range prev.Results {
			before[r.RuleID] = r
		}
	}

	var out []Diff
	for _, r := range cur.Results {
		old, existed := before[r.RuleID]
		if existed && old.Status == r.Status {
			continue
		}

		d := Diff{RuleID: r.RuleID, To: r.Status, Severity: r.Severity, Message: r.Message}
		if existed {
			d.From = old.Status
		} else {
			d.From = rules.StatusSkipped
		}
		out = append(out, d)
	}

	// Rules that disappeared, such as after a config reload.
	if prev != nil {
		now := map[string]bool{}
		for _, r := range cur.Results {
			now[r.RuleID] = true
		}
		for _, r := range prev.Results {
			if !now[r.RuleID] {
				out = append(out, Diff{RuleID: r.RuleID, From: r.Status, To: rules.StatusSkipped})
			}
		}
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].RuleID < out[j].RuleID })
	return out
}

func Snapshot(root string, ignore []string) (*tree.Snapshot, error) {
	return tree.Scan(root, ignore)
}

func fileHash(path string) string {
	if path == "" {
		return ""
	}
	info, err := os.Stat(path)
	if err != nil {
		return "absent"
	}
	return hash(fmt.Sprintf("%s:%d:%d", filepath.Base(path), info.Size(), info.ModTime().UnixNano()))
}

func hash(s string) string {
	const (
		offset = 14695981039346656037
		prime  = 1099511628211
	)
	var h uint64 = offset
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime
	}
	return fmt.Sprintf("%016x", h)
}
