package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/tree"
)

// stage materialises a project layout and returns its root.
func stage(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// testConfig activates the named rules with their default severities.
func testConfig(t *testing.T, ids []string, ignore []string) *config.Config {
	t.Helper()
	cfg := &config.Config{Ignore: ignore, ActiveRules: map[string]*config.ActiveRule{}}
	meta := config.GetRulesMetadata()
	for _, id := range ids {
		m, ok := meta.Rules[id]
		if !ok {
			t.Fatalf("unknown rule id in test: %s", id)
		}
		cfg.ActiveRules[id] = &config.ActiveRule{ID: id, Metadata: m, Severity: m.DefaultSeverity}
	}
	return cfg
}

func newTestContext(root, projectType string) *Context {
	return &Context{ProjectPath: root, ProjectType: projectType}
}

// runCheck executes the engine against a staged layout.
func runCheck(t *testing.T, root, projectType string, ruleIDs []string, ignore []string) *ExecutionResult {
	t.Helper()
	cfg := testConfig(t, ruleIDs, ignore)
	snap, err := tree.Scan(root, ignore)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	res, err := ExecuteSnapshot(cfg, newTestContext(root, projectType), snap)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	return res
}

func resultByRule(res *ExecutionResult, id string) RuleResult {
	for _, r := range res.Results {
		if r.RuleID == id {
			return r
		}
	}
	return RuleResult{RuleID: id, Status: StatusSkipped, Message: "rule not executed"}
}

func TestGoProjectWithNestedTestsPasses(t *testing.T) {
	// BUG-04: filepath.Glob returned 0 matches for **/*_test.go, so this rule
	// failed on every Go project regardless of its tests.
	root := stage(t, map[string]string{
		"go.mod":                     "module demo\n",
		"internal/rules/eng_test.go": "package rules\n",
		"cmd/main.go":                "package main\n",
	})

	res := runCheck(t, root, "go", []string{"tests_folder"}, nil)
	got := resultByRule(res, "tests_folder")
	if got.Status != StatusPassed {
		t.Errorf("tests_folder = %s (%s), want passed; evidence=%q", got.Status, got.Message, got.Evidence)
	}
}

func TestGoProjectWithoutTestsFails(t *testing.T) {
	root := stage(t, map[string]string{"go.mod": "module demo\n"})

	res := runCheck(t, root, "go", []string{"tests_folder"}, nil)
	if got := resultByRule(res, "tests_folder"); got.Status != StatusFailed {
		t.Errorf("tests_folder = %s, want failed", got.Status)
	}
}

func TestSpecGlobsResolve(t *testing.T) {
	root := stage(t, map[string]string{
		"package.json":  "{}",
		"src/a.test.ts": "x",
		"src/b.spec.js": "x",
	})

	res := runCheck(t, root, "nodejs", []string{"tests_folder"}, nil)
	if got := resultByRule(res, "tests_folder"); got.Status != StatusPassed {
		t.Errorf("tests_folder = %s (%s), want passed", got.Status, got.Message)
	}
}

func TestEmptyDirectoryDoesNotSatisfyRule(t *testing.T) {
	root := stage(t, map[string]string{"README.md": "hi"})
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}

	res := runCheck(t, root, "generic", []string{"docs_folder"}, nil)
	if got := resultByRule(res, "docs_folder"); got.Status != StatusFailed {
		t.Errorf("docs_folder = %s, want failed for an empty dir", got.Status)
	}
}

func TestIgnoreIsAppliedDuringCheck(t *testing.T) {
	// BUG-09: ignore was validated but never applied, so vendored trees
	// satisfied tests_folder and src_folder.
	root := stage(t, map[string]string{
		"node_modules/pkg/a.test.js": "x",
		"vendor/lib/b_test.go":       "x",
	})

	res := runCheck(t, root, "nodejs", []string{"tests_folder", "src_folder"},
		[]string{"node_modules/", "vendor/"})
	if got := resultByRule(res, "tests_folder"); got.Status != StatusFailed {
		t.Errorf("tests_folder = %s, want failed when only ignored paths match", got.Status)
	}
	if got := resultByRule(res, "src_folder"); got.Status != StatusFailed {
		t.Errorf("src_folder = %s, want failed when only ignored paths match", got.Status)
	}
}

func TestUnknownProjectTypeUsesGenericAndStillFails(t *testing.T) {
	// BUG-03/BUG-H: an unknown type resolved to no patterns and the rule was
	// reported as passed, hiding a real violation. It must evaluate the generic
	// tier instead.
	root := stage(t, map[string]string{})

	res := runCheck(t, root, "klingon", []string{"src_folder"}, nil)
	got := resultByRule(res, "src_folder")
	if got.Status != StatusFailed {
		t.Errorf("src_folder = %s, want failed for an unknown project type", got.Status)
	}
	if got.Evidence != "src/" {
		t.Errorf("Evidence = %q, want the generic pattern src/", got.Evidence)
	}
}

func TestUnknownProjectTypePassesWhenGenericPatternExists(t *testing.T) {
	root := stage(t, map[string]string{"src/index.ts": "x"})

	res := runCheck(t, root, "brandnewlang", []string{"src_folder"}, nil)
	if got := resultByRule(res, "src_folder"); got.Status != StatusPassed {
		t.Errorf("src_folder = %s, want passed via the generic fallback", got.Status)
	}
}

func TestRuleWithNoApplicableTierIsSkipped(t *testing.T) {
	// package_manager declares only language tiers, so for an unrelated type it
	// has nothing to check. It must be skipped, and skipped must not count as a
	// pass: borrowing Go's expectations here would be misleading.
	root := stage(t, map[string]string{})

	res := runCheck(t, root, "klingon", []string{"package_manager"}, nil)
	got := resultByRule(res, "package_manager")
	if got.Status != StatusSkipped {
		t.Errorf("package_manager = %s, want skipped for an unrelated type", got.Status)
	}

	// The same rule must still fail for a type it does cover.
	res = runCheck(t, root, "go", []string{"package_manager"}, nil)
	if got := resultByRule(res, "package_manager"); got.Status != StatusFailed {
		t.Errorf("package_manager = %s for go, want failed (no go.mod)", got.Status)
	}
}

func TestEvidenceNamesTheMatchedPattern(t *testing.T) {
	root := stage(t, map[string]string{"README": "hello"})

	res := runCheck(t, root, "generic", []string{"readme"}, nil)
	got := resultByRule(res, "readme")
	if got.Evidence != "README" {
		t.Errorf("Evidence = %q, want %q", got.Evidence, "README")
	}
}

func TestSummaryCountsSkippedSeparately(t *testing.T) {
	root := stage(t, map[string]string{"README.md": "hi"})

	res := runCheck(t, root, "generic", []string{"readme", "package_manager"}, nil)
	if res.Summary.Passed != 1 {
		t.Errorf("Passed = %d, want 1", res.Summary.Passed)
	}
	if res.Summary.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1", res.Summary.Skipped)
	}
	if res.Summary.Failed != 0 {
		t.Errorf("Failed = %d, want 0", res.Summary.Failed)
	}
	if res.Summary.Total != 2 {
		t.Errorf("Total = %d, want 2", res.Summary.Total)
	}
}

// TestSkippedRulesDoNotProduceWarnings guards the report contract: a rule that
// never ran must not make the run look dirty.
func TestSkippedRulesDoNotProduceWarnings(t *testing.T) {
	root := stage(t, map[string]string{})

	res := runCheck(t, root, "klingon", []string{"package_manager", "lockfile"}, nil)
	if res.Status != OutcomePassed {
		t.Errorf("Status = %s, want passed when every rule was skipped", res.Status)
	}
	if res.Summary.Warnings != 0 || res.Summary.Errors != 0 {
		t.Errorf("skipped rules leaked into severities: %+v", res.Summary)
	}
}

func TestResultsAreSortedDeterministically(t *testing.T) {
	root := stage(t, map[string]string{})

	for i := 0; i < 30; i++ {
		res := runCheck(t, root, "generic", []string{
			"readme", "license", "gitignore", "changelog", "editorconfig",
		}, nil)
		for j := 1; j < len(res.Results); j++ {
			if res.Results[j-1].RuleID > res.Results[j].RuleID {
				t.Fatalf("results not sorted: %q before %q",
					res.Results[j-1].RuleID, res.Results[j].RuleID)
			}
		}
	}
}
