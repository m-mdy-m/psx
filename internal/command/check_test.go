package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/rules"
)

// failed builds a result where every listed rule failed at error severity.
func failed(ids ...string) *rules.ExecutionResult {
	res := &rules.ExecutionResult{}
	for _, id := range ids {
		res.Results = append(res.Results, rules.RuleResult{
			RuleID:   id,
			Status:   rules.StatusFailed,
			Severity: config.SeverityError,
			Message:  "missing",
		})
	}
	rules.Recount(res)
	return res
}

func statuses(res *rules.ExecutionResult) map[string]rules.Status {
	out := map[string]rules.Status{}
	for _, r := range res.Results {
		out[r.RuleID] = r.Status
	}
	return out
}

// Bootstrapping a baseline is the only way to adopt one on an existing repo, so a
// missing file has to be recorded rather than treated as an error.
func TestMissingBaselineFileIsRecordedNotFatal(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".psx-baseline.txt")
	res := failed("readme", "tests_folder")

	if err := applyBaseline(res, path); err != nil {
		t.Fatalf("applyBaseline on a missing file: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("baseline was not written: %v", err)
	}
	for _, want := range []string{"readme", "tests_folder"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("baseline file missing %q:\n%s", want, data)
		}
	}
	for id, st := range statuses(res) {
		if st != rules.StatusPassed {
			t.Errorf("%s should be forgiven by the baseline it just wrote, got %s", id, st)
		}
	}
}

func TestBaselineForgivesOnlyRecordedRules(t *testing.T) {
	path := filepath.Join(t.TempDir(), "baseline.txt")
	if err := os.WriteFile(path, []byte("readme\n# a comment\n\ntests_folder\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res := failed("readme", "tests_folder", "license")
	if err := applyBaseline(res, path); err != nil {
		t.Fatalf("applyBaseline: %v", err)
	}

	got := statuses(res)
	for _, id := range []string{"readme", "tests_folder"} {
		if got[id] != rules.StatusPassed {
			t.Errorf("%s is recorded, want passed; got %s", id, got[id])
		}
	}
	if got["license"] != rules.StatusFailed {
		t.Errorf("license is not recorded, want failed; got %s", got["license"])
	}
	if res.Summary.Errors != 1 {
		t.Errorf("summary should count only the unrecorded failure, got %d errors", res.Summary.Errors)
	}
	if res.Status != rules.OutcomeFailed {
		t.Errorf("status should stay failed while an unrecorded error remains, got %s", res.Status)
	}
}

func TestFullyBaselinedRunPasses(t *testing.T) {
	path := filepath.Join(t.TempDir(), "baseline.txt")
	if err := os.WriteFile(path, []byte("readme\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res := failed("readme")
	if err := applyBaseline(res, path); err != nil {
		t.Fatalf("applyBaseline: %v", err)
	}

	if res.Status != rules.OutcomePassed {
		t.Errorf("a fully forgiven run should pass, got %s", res.Status)
	}
	if res.Summary.Errors != 0 {
		t.Errorf("errors should be cleared, got %d", res.Summary.Errors)
	}
}

func TestBaselineAcceptsJSONArrays(t *testing.T) {
	path := filepath.Join(t.TempDir(), "baseline.json")
	if err := os.WriteFile(path, []byte(`["readme","license"]`), 0o644); err != nil {
		t.Fatal(err)
	}

	res := failed("readme", "license", "gitignore")
	if err := applyBaseline(res, path); err != nil {
		t.Fatalf("applyBaseline: %v", err)
	}

	if statuses(res)["gitignore"] != rules.StatusFailed {
		t.Error("gitignore is absent from the JSON baseline, want failed")
	}
}

func TestShouldFailHonoursTheThreshold(t *testing.T) {
	res := failed("readme", "gitignore") // both errors
	tests := []struct {
		failOn string
		want   bool
	}{
		{"error", true},
		{"warning", true},
		{"none", false},
	}
	for _, tc := range tests {
		if got := shouldFail(res, tc.failOn); got != tc.want {
			t.Errorf("shouldFail(%q) = %v, want %v", tc.failOn, got, tc.want)
		}
	}

	for i := range res.Results {
		res.Results[i].Severity = config.SeverityInfo
	}
	rules.Recount(res)
	if shouldFail(res, "warning") {
		t.Error("an info-only run must not fail at --fail-on warning")
	}
}
