package watch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/tree"
)

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func rr(id string, st rules.Status, sev config.Severity, msg string) rules.RuleResult {
	return rules.RuleResult{RuleID: id, Status: st, Severity: sev, Message: msg}
}

func resultOf(list ...rules.RuleResult) *rules.ExecutionResult {
	out := &rules.ExecutionResult{Results: list}
	rules.Recount(out)
	return out
}

// --- Fingerprint -------------------------------------------------------------

// A fingerprint that changes for an unchanged tree makes watch re-check forever.
func TestFingerprintIsStableForAnUnchangedTree(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "go.mod"), "module x\n")
	write(t, filepath.Join(root, "cmd", "main.go"), "package main\n")

	first, count := scan(t, root)
	second, _ := scan(t, root)

	if first != second {
		t.Errorf("fingerprint changed for an untouched tree: %q then %q", first, second)
	}
	if count != 2 {
		t.Errorf("file count = %d, want 2", count)
	}
}

func TestFingerprintDetectsContentRenameAndRemoval(t *testing.T) {
	t.Run("edit", func(t *testing.T) {
		root := t.TempDir()
		write(t, filepath.Join(root, "a.txt"), "one\n")
		before, _ := scan(t, root)

		// Same byte count, different content: only mtime can catch this.
		write(t, filepath.Join(root, "a.txt"), "two\n")
		if after, _ := scan(t, root); after == before {
			t.Error("an edit that preserves file size was not detected")
		}
	})

	t.Run("rename", func(t *testing.T) {
		root := t.TempDir()
		write(t, filepath.Join(root, "a.txt"), "x\n")
		before, _ := scan(t, root)
		if err := os.Rename(filepath.Join(root, "a.txt"), filepath.Join(root, "b.txt")); err != nil {
			t.Fatal(err)
		}
		if after, _ := scan(t, root); after == before {
			t.Error("a rename was not detected")
		}
	})

	t.Run("remove", func(t *testing.T) {
		root := t.TempDir()
		write(t, filepath.Join(root, "a.txt"), "x\n")
		before, _ := scan(t, root)
		if err := os.Remove(filepath.Join(root, "a.txt")); err != nil {
			t.Fatal(err)
		}
		if after, _ := scan(t, root); after == before {
			t.Error("a removal was not detected")
		}
	})
}

// .git, node_modules and vendor change on every commit or install. Watching them
// burns a re-check each time for a verdict that cannot change.
func TestFingerprintIgnoresVendoredAndGitDirectories(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "main.go"), "package main\n")
	before, count := scan(t, root)
	if count != 1 {
		t.Errorf("file count = %d, want 1", count)
	}

	write(t, filepath.Join(root, "vendor", "lib.go"), "package lib\n")
	write(t, filepath.Join(root, "node_modules", "dep.js"), "1\n")
	write(t, filepath.Join(root, ".git", "HEAD"), "ref: main\n")

	if after, count := scan(t, root); after != before || count != 1 {
		t.Errorf("vendored or git content changed the fingerprint (%q, %d files)", after, count)
	}
}

// The whole point of routing through tree.Scan: a user's own ignore list prunes
// the change detection too, not just the check.
func TestFingerprintHonoursConfiguredIgnorePatterns(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "main.go"), "package main\n")
	before, _ := scan(t, root)

	write(t, filepath.Join(root, "dist", "bundle.js"), "1\n")
	if after, count := scanWith(t, root, []string{"dist/"}); after != before || count != 1 {
		t.Errorf("an ignored directory changed the fingerprint (%q, %d files)", after, count)
	}
}

// .git is not negatable. Someone writing "!/.git/" has no idea it costs them a
// full re-check on every commit, for a verdict that cannot change.
func TestAlwaysIgnoredIsNotNegatableByTheUser(t *testing.T) {
	for _, reinclude := range []string{"!/.git/", "!.git/", "/.git", ".git/"} {
		t.Run(reinclude, func(t *testing.T) {
			root := t.TempDir()
			write(t, filepath.Join(root, "main.go"), "package main\n")
			before, _ := scan(t, root)

			write(t, filepath.Join(root, ".git", "index"), "bin")
			if after, count := scanWith(t, root, []string{reinclude}); after != before || count != 1 {
				t.Errorf("%q re-opened .git to the watcher (fingerprint %q, %d files)", reinclude, after, count)
			}
		})
	}
}

// --- Config hashing ----------------------------------------------------------

func TestConfigFileIsFingerprintedSeparately(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")
	cfg := filepath.Join(root, "psx.yml")
	write(t, cfg, "version: 1\n")

	if fileHash(cfg) == fileHash("") {
		t.Error("an existing config must not hash like a missing one")
	}
	if fileHash("") != "" {
		t.Error("no config file should hash to the empty string")
	}
	if fileHash(cfg) == "absent" {
		t.Error("an existing config hashed as absent")
	}
	if fileHash(filepath.Join(root, "nope.yml")) != "absent" {
		t.Error("a missing config should hash as absent")
	}

	write(t, cfg, "version: 1\nrules: {}\n")
	if fileHash(cfg) == fileHash(cfg+"\n") {
		t.Error("editing the config did not change its hash")
	}
}

// --- Compare -----------------------------------------------------------------

func TestCompareReportsOnlyStatusChanges(t *testing.T) {
	prev := resultOf(
		rr("readme", rules.StatusFailed, config.SeverityError, "missing"),
		rr("license", rules.StatusPassed, config.SeverityWarning, ""),
		rr("gitignore", rules.StatusFailed, config.SeverityInfo, "missing"),
	)
	cur := resultOf(
		rr("readme", rules.StatusPassed, config.SeverityError, ""),
		rr("license", rules.StatusPassed, config.SeverityWarning, ""),
		rr("gitignore", rules.StatusFailed, config.SeverityInfo, "missing"),
	)

	got := Compare(prev, cur)
	if len(got) != 1 {
		t.Fatalf("got %d diffs, want 1: %+v", len(got), got)
	}
	if got[0].RuleID != "readme" {
		t.Errorf("RuleID = %q, want readme", got[0].RuleID)
	}
	if got[0].From != rules.StatusFailed || got[0].To != rules.StatusPassed {
		t.Errorf("got %s -> %s, want failed -> passed", got[0].From, got[0].To)
	}
}

func TestCompareTreatsNewRulesAsArriving(t *testing.T) {
	prev := resultOf(rr("readme", rules.StatusPassed, config.SeverityError, ""))
	cur := resultOf(
		rr("readme", rules.StatusPassed, config.SeverityError, ""),
		rr("tests_folder", rules.StatusFailed, config.SeverityError, "No tests found"),
	)

	got := Compare(prev, cur)
	if len(got) != 1 || got[0].RuleID != "tests_folder" {
		t.Fatalf("got %+v, want a single diff for tests_folder", got)
	}
	if got[0].From != rules.StatusSkipped {
		t.Errorf("a rule absent from the previous run should come From skipped, got %s", got[0].From)
	}
	if got[0].Message != "No tests found" {
		t.Errorf("Message = %q, want the current message so watch can print it", got[0].Message)
	}
}

func TestCompareReportsRulesThatVanished(t *testing.T) {
	prev := resultOf(
		rr("readme", rules.StatusPassed, config.SeverityError, ""),
		rr("docker_compose", rules.StatusFailed, config.SeverityInfo, "missing"),
	)
	cur := resultOf(rr("readme", rules.StatusPassed, config.SeverityError, ""))

	got := Compare(prev, cur)
	if len(got) != 1 {
		t.Fatalf("got %d diffs, want 1: %+v", len(got), got)
	}
	if got[0].RuleID != "docker_compose" || got[0].To != rules.StatusSkipped {
		t.Errorf("got %+v, want docker_compose -> skipped", got[0])
	}
}

func TestCompareIsSortedAndHandlesANilPrevious(t *testing.T) {
	cur := resultOf(
		rr("zed", rules.StatusFailed, config.SeverityInfo, ""),
		rr("adr", rules.StatusFailed, config.SeverityInfo, ""),
		rr("makefile", rules.StatusFailed, config.SeverityInfo, ""),
	)

	got := Compare(nil, cur)
	if len(got) != 3 {
		t.Fatalf("got %d diffs, want 3", len(got))
	}
	for i, want := range []string{"adr", "makefile", "zed"} {
		if got[i].RuleID != want {
			t.Errorf("position %d = %q, want %q", i, got[i].RuleID, want)
		}
	}
	if got[0].From != rules.StatusSkipped {
		t.Errorf("with no previous run every rule should read From skipped, got %s", got[0].From)
	}
}

// Compare is exported and takes pointers, so a nil current result must not panic
// a watcher that has just started.
func TestCompareSurvivesANilResult(t *testing.T) {
	if got := Compare(nil, nil); len(got) != 0 {
		t.Errorf("Compare(nil, nil) = %+v, want empty", got)
	}
}

// --- isClean -----------------------------------------------------------------

func TestIsCleanIgnoresSkippedRules(t *testing.T) {
	clean := resultOf(
		rr("readme", rules.StatusPassed, config.SeverityError, ""),
		rr("package_manager", rules.StatusSkipped, config.SeverityError, "not applicable"),
	)
	if !isClean(clean) {
		t.Error("a run whose only non-passing rule was skipped is clean")
	}

	dirty := resultOf(rr("readme", rules.StatusFailed, config.SeverityError, ""))
	if isClean(dirty) {
		t.Error("a run with a failure is not clean")
	}
	if isClean(nil) {
		t.Error("a nil result is not clean; a --once watcher would exit immediately")
	}
}

// --- Options -----------------------------------------------------------------

func TestDefaultsAreUsable(t *testing.T) {
	d := Defaults()
	if d.Interval <= 0 || d.Debounce <= 0 {
		t.Fatalf("Defaults() = %+v, want positive interval and debounce", d)
	}
	if d.Debounce >= d.Interval {
		t.Errorf("Debounce (%v) must be shorter than Interval (%v), otherwise a change is noticed late",
			d.Debounce, d.Interval)
	}
	if d.Interval > 10*time.Second {
		t.Errorf("default Interval %v is too long to feel responsive", d.Interval)
	}
}

// --- Run ---------------------------------------------------------------------

func fastOptions(root string) Options {
	o := Defaults()
	o.Root = root
	o.Interval = 10 * time.Millisecond
	o.Debounce = time.Millisecond
	return o
}

// The snapshot that detects a change must be the one the check is evaluated
// against, otherwise a file can change between fingerprint and verdict.
func TestRunHandsTheDetectedSnapshotToTheCheck(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seen := make(chan *tree.Snapshot, 1)
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, fastOptions(root),
			func(snap *tree.Snapshot) (*rules.ExecutionResult, error) {
				select {
				case seen <- snap:
				default:
				}
				return resultOf(rr("readme", rules.StatusPassed, config.SeverityError, "")), nil
			},
			nil)
	}()

	select {
	case snap := <-seen:
		if snap == nil {
			t.Fatal("Run passed a nil snapshot to the check")
		}
		if !snap.HasFile("a.txt") {
			t.Error("the snapshot handed to the check does not contain the project files")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the check was never called")
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// A watcher must not spin: an unchanged tree produces no further check calls.
func TestRunDoesNotRecheckAnUnchangedTree(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	var calls int
	err := Run(ctx, fastOptions(root),
		func(*tree.Snapshot) (*rules.ExecutionResult, error) {
			calls++
			return resultOf(rr("readme", rules.StatusFailed, config.SeverityError, "missing")), nil
		},
		nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if calls != 1 {
		t.Errorf("the check ran %d times on an unchanging tree, want 1", calls)
	}
}

func TestRunRechecksAfterAChange(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		// Touch the file until the watcher notices.
		for range 40 {
			select {
			case <-ctx.Done():
				return
			default:
			}
			write(t, filepath.Join(root, "a.txt"), "changed\n")
			time.Sleep(20 * time.Millisecond)
		}
	}()

	calls := make(chan int, 64)
	err := Run(ctx, fastOptions(root),
		func(*tree.Snapshot) (*rules.ExecutionResult, error) {
			select {
			case calls <- 1:
			default:
			}
			return resultOf(rr("readme", rules.StatusFailed, config.SeverityError, "missing")), nil
		},
		nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(calls) < 2 {
		t.Errorf("the watcher ran %d checks, want at least 2 after a change", len(calls))
	}
}

// --once is how CI uses the watcher: exit as soon as the project is clean.
func TestRunOnceStopsOnACleanProject(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	opts := fastOptions(root)
	opts.OnceClean = true

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var calls int
	err := Run(ctx, opts,
		func(*tree.Snapshot) (*rules.ExecutionResult, error) {
			calls++
			return resultOf(rr("readme", rules.StatusPassed, config.SeverityError, "")), nil
		},
		nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if calls != 1 {
		t.Errorf("the check ran %d times, want 1: --once should exit on a clean project", calls)
	}
}

func TestRunOnceKeepsWatchingUntilClean(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	opts := fastOptions(root)
	opts.OnceClean = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var calls int
	go func() {
		// Break the project, then repair it.
		write(t, filepath.Join(root, "a.txt"), "broken\n")
		time.Sleep(150 * time.Millisecond)
		write(t, filepath.Join(root, "README.md"), "# ok\n")
	}()

	err := Run(ctx, opts,
		func(snap *tree.Snapshot) (*rules.ExecutionResult, error) {
			calls++
			if snap.HasFile("README.md") {
				return resultOf(rr("readme", rules.StatusPassed, config.SeverityError, "")), nil
			}
			return resultOf(rr("readme", rules.StatusFailed, config.SeverityError, "missing")), nil
		},
		nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if calls < 2 {
		t.Errorf("the check ran %d times, want at least 2: a failing project must keep the watcher alive", calls)
	}
}

func TestRunReportsAMissingRoot(t *testing.T) {
	opts := fastOptions(filepath.Join(t.TempDir(), "nope"))
	err := Run(context.Background(), opts,
		func(*tree.Snapshot) (*rules.ExecutionResult, error) {
			t.Error("the check must not run when the project directory does not exist")
			return resultOf(), nil
		},
		nil)
	if err == nil {
		t.Error("watching a directory that does not exist should report an error")
	}
}

func TestRunPropagatesAFirstCheckFailure(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "a.txt"), "x\n")

	want := os.ErrPermission
	err := Run(context.Background(), fastOptions(root),
		func(*tree.Snapshot) (*rules.ExecutionResult, error) { return nil, want },
		nil)
	if err == nil {
		t.Error("a failing first check should stop the watcher")
	}
}

// --- helpers -----------------------------------------------------------------

func scan(t *testing.T, root string) (string, int) {
	t.Helper()
	return scanWith(t, root, nil)
}

func scanWith(t *testing.T, root string, ignore []string) (string, int) {
	t.Helper()
	snap, err := tree.Scan(root, watchIgnore(ignore))
	if err != nil {
		t.Fatalf("tree.Scan(%s): %v", root, err)
	}
	return snap.Fingerprint()
}
