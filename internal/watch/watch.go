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
	"time"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/tree"
)

// alwaysIgnored never contributes to a project's structure, and watching them
// wastes a re-check per git operation or per dependency install.
var alwaysIgnored = []string{".git/", "node_modules/", "vendor/"}

// watchIgnore layers the always-ignored directories over the user's patterns, so
// the walk that detects a change prunes exactly what the check itself prunes.
//
// The defaults go last because gitignore negation is last-match-wins. Someone who
// writes "!/.git/" would otherwise re-open .git to the watcher, and every commit
// would cost a full re-check for a verdict that cannot change.
func watchIgnore(user []string) []string {
	out := make([]string, 0, len(user)+len(alwaysIgnored))
	out = append(out, user...)
	return append(out, alwaysIgnored...)
}

// Run watches the project until ctx is cancelled.
//
// The tree is walked once per change: the snapshot that detects the change is the
// same one the check is evaluated against, so a file cannot change between the
// fingerprint and the verdict.
func Run(ctx context.Context, opts Options, runCheck func(*tree.Snapshot) (*rules.ExecutionResult, error), onChange func(prev, cur *rules.ExecutionResult)) error {
	// time.NewTicker panics on a non-positive interval, and a zero here would
	// mean an unpopulated Options rather than a deliberate request.
	if opts.Interval <= 0 {
		return fmt.Errorf("watch interval must be positive, got %v", opts.Interval)
	}
	if opts.Debounce < 0 {
		return fmt.Errorf("watch debounce must not be negative, got %v", opts.Debounce)
	}

	root := opts.Root

	snap, err := tree.Scan(root, watchIgnore(opts.Ignore))
	if err != nil {
		return fmt.Errorf("initial scan: %w", err)
	}
	current, files := snap.Fingerprint()
	logger.Verbosef("watching %s (%d files)", root, files)

	prevResult, err := runCheck(snap)
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

	configHash := fileHash(opts.ConfigFile)
	var pending bool
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-debounce.C:
			// Coalesce bursts of writes into a single re-check.
			pending = false
			next, err := tree.Scan(root, watchIgnore(opts.Ignore))
			if err != nil {
				logger.Verbosef("rescan failed: %v", err)
				continue
			}
			fingerprint, _ := next.Fingerprint()
			cfgHash := fileHash(opts.ConfigFile)
			if fingerprint == current && cfgHash == configHash {
				continue
			}
			current, configHash = fingerprint, cfgHash

			result, err := runCheck(next)
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
	if cur == nil {
		return nil
	}

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
