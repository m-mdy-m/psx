package command

import (
	"testing"

	"github.com/m-mdy-m/psx/internal/rules"
)

// watch registered --baseline but never used it, so a known-issue run looked like
// a fresh failure in the watch line and in the per-rule change feed.
func TestBaselineIsAppliedOnEveryWatchIteration(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.psx-baseline.txt"

	res := failed("readme", "tests_folder")
	if err := applyBaseline(res, path); err != nil {
		t.Fatalf("bootstrap baseline: %v", err)
	}

	// A later iteration sees the same two rules still failing.
	next := failed("readme", "tests_folder")
	if err := applyBaselineTo(next, path); err != nil {
		t.Fatalf("applyBaselineTo: %v", err)
	}

	if next.Summary.Failed != 0 {
		t.Errorf("recorded issues should not resurface in watch, got %d failures", next.Summary.Failed)
	}
	if next.Status != rules.OutcomePassed {
		t.Errorf("status should be clean while only recorded issues remain, got %s", next.Status)
	}
}

func TestBaselineStillLetsNewWatchIssuesThrough(t *testing.T) {
	path := t.TempDir() + "/baseline.txt"

	first := failed("readme")
	if err := applyBaseline(first, path); err != nil {
		t.Fatalf("bootstrap baseline: %v", err)
	}

	// The user then breaks something new while editing.
	second := failed("readme", "license")
	if err := applyBaselineTo(second, path); err != nil {
		t.Fatalf("applyBaselineTo: %v", err)
	}

	if second.Summary.Failed != 1 {
		t.Fatalf("the new failure should survive the baseline, got %d", second.Summary.Failed)
	}
	if got := second.Results[1]; got.RuleID != "license" || got.Status != rules.StatusFailed {
		t.Errorf("license should still be reported as failing, got %+v", got)
	}
}

// An unreadable baseline must be reported rather than silently ignored, which
// would hide every recorded issue from the watch feed.
func TestBaselineReadFailureIsReported(t *testing.T) {
	if err := applyBaselineTo(failed("readme"), t.TempDir()); err == nil {
		t.Error("a directory is not a readable baseline, want an error")
	}
}
