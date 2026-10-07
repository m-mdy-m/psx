package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/rules"
)

// result builds a result covering passed, failed and skipped rules.
func result() *rules.ExecutionResult {
	return &rules.ExecutionResult{
		Context: &rules.Context{
			ProjectPath: "/tmp/proj",
			ProjectType: "go",
		},
		Results: []rules.RuleResult{
			{RuleID: "readme", Status: rules.StatusPassed, Severity: config.SeverityError, Message: "OK", Evidence: "README.md"},
			{RuleID: "license", Status: rules.StatusFailed, Severity: config.SeverityWarning,
				Message: "No LICENSE file found", FixHint: "psx fix --rule license", Evidence: "LICENSE"},
			{RuleID: "tests_folder", Status: rules.StatusFailed, Severity: config.SeverityError,
				Message: "No tests found", FixHint: "psx fix --rule tests_folder"},
			{RuleID: "docker_compose", Status: rules.StatusFailed, Severity: config.SeverityInfo,
				Message: "No docker-compose configuration found"},
			{RuleID: "package_manager", Status: rules.StatusSkipped, Severity: config.SeverityError,
				Message: "Not applicable to go projects"},
		},
		Summary: rules.Summary{Total: 5, Passed: 1, Failed: 3, Skipped: 1,
			Errors: 1, Warnings: 1, Info: 1},
		Status: rules.OutcomeFailed,
	}
}

func render(t *testing.T, opts Options) string {
	t.Helper()
	opts.NoColor = true
	opts.Version = "test"
	r, err := New(opts)
	if err != nil {
		t.Fatalf("New(%s): %v", opts.Format, err)
	}
	var buf bytes.Buffer
	if err := r.Render(&buf, result()); err != nil {
		t.Fatalf("Render(%s): %v", opts.Format, err)
	}
	return buf.String()
}

func TestEveryFormatRenders(t *testing.T) {
	formats := []string{
		flags.FormatTable, flags.FormatCompact, flags.FormatJSON, flags.FormatNDJSON,
		flags.FormatSARIF, flags.FormatGitHub, flags.FormatJUnit, flags.FormatMarkdown,
	}
	for _, f := range formats {
		t.Run(f, func(t *testing.T) {
			out := render(t, Options{Format: f, NoColor: true})
			if strings.TrimSpace(out) == "" {
				t.Fatalf("%s produced no output", f)
			}
		})
	}
}

func TestUnknownFormatIsRejected(t *testing.T) {
	if _, err := New(Options{Format: "yaml"}); err == nil {
		t.Fatal("expected an error for an unsupported format")
	}
}

// TestJSONIsPureAndValid is the guard for `psx check -o json | jq`.
func TestJSONIsPureAndValid(t *testing.T) {
	out := render(t, Options{Format: flags.FormatJSON})

	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	for _, key := range []string{"version", "status", "summary", "context", "results"} {
		if _, ok := doc[key]; !ok {
			t.Errorf("missing key %q in JSON report", key)
		}
	}
}

func TestJSONOmitsPassedRulesByDefault(t *testing.T) {
	var doc struct {
		Results []struct {
			RuleID string `json:"rule_id"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(render(t, Options{Format: flags.FormatJSON})), &doc); err != nil {
		t.Fatal(err)
	}
	for _, r := range doc.Results {
		if r.RuleID == "readme" {
			t.Error("passing rules must be omitted unless requested")
		}
	}
	if len(doc.Results) != 4 {
		t.Errorf("results = %d, want 3 failures plus 1 skipped", len(doc.Results))
	}
}

func TestJSONIncludesPassedWhenRequested(t *testing.T) {
	var doc struct {
		Results []struct {
			RuleID string `json:"rule_id"`
		} `json:"results"`
	}
	opts := Options{Format: flags.FormatJSON, ShowPassed: true}
	if err := json.Unmarshal([]byte(render(t, opts)), &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Results) != 5 {
		t.Errorf("results = %d, want all 5 rules including the passed one", len(doc.Results))
	}
}

func TestNDJSONIsOneObjectPerLine(t *testing.T) {
	out := render(t, Options{Format: flags.FormatNDJSON})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected a summary line plus rule lines, got %d", len(lines))
	}
	for i, line := range lines {
		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Fatalf("line %d is not valid JSON: %v\n%s", i+1, err, line)
		}
		if _, ok := obj["type"]; !ok {
			t.Errorf("line %d has no type field", i+1)
		}
	}
}

func TestSARIFIsValidShape(t *testing.T) {
	out := render(t, Options{Format: flags.FormatSARIF})
	var doc struct {
		Version string `json:"version"`
		Runs    []struct {
			Tool struct {
				Driver struct {
					Name  string `json:"name"`
					Rules []struct {
						ID string `json:"id"`
					} `json:"rules"`
				} `json:"driver"`
			} `json:"tool"`
			Results []struct {
				RuleID string `json:"ruleId"`
				Level  string `json:"level"`
			} `json:"results"`
		} `json:"runs"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("SARIF is not valid JSON: %v\n%s", err, out)
	}
	if doc.Version != "2.1.0" {
		t.Errorf("version = %q, want 2.1.0", doc.Version)
	}
	if len(doc.Runs) != 1 {
		t.Fatalf("runs = %d, want 1", len(doc.Runs))
	}
	if doc.Runs[0].Tool.Driver.Name != "psx" {
		t.Errorf("driver = %q, want psx", doc.Runs[0].Tool.Driver.Name)
	}
	if len(doc.Runs[0].Results) != 3 {
		t.Errorf("results = %d, want 3", len(doc.Runs[0].Results))
	}
}

func TestJUnitCountsMatchSummary(t *testing.T) {
	out := render(t, Options{Format: flags.FormatJUnit})
	if !strings.Contains(out, `failures="3"`) {
		t.Errorf("junit report does not report 3 failures:\n%s", out)
	}
	if strings.Count(out, "<failure") != 3 {
		t.Errorf("expected 3 failure elements, got %d", strings.Count(out, "<failure"))
	}
}

func TestGitHubAnnotationsEscapeNewlines(t *testing.T) {
	res := result()
	res.Results[1].Message = "line one\nline two"

	opts := Options{Format: flags.FormatGitHub, NoColor: true}
	r, _ := New(opts)
	var buf bytes.Buffer
	if err := r.Render(&buf, res); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "line one\nline two") {
		t.Error("workflow command message contains a raw newline")
	}
	if !strings.Contains(buf.String(), "%0A") {
		t.Error("newline was not escaped as %0A")
	}
}

func TestLevelFilterHidesOtherSeverities(t *testing.T) {
	out := render(t, Options{Format: flags.FormatCompact, Level: "error"})
	if strings.Contains(out, "license") {
		t.Error("a warning should be hidden at --level error")
	}
	if !strings.Contains(out, "tests_folder") {
		t.Error("an error should be visible at --level error")
	}
}

// A filtered-out finding must not still be counted in the summary, otherwise
// `--level warning` prints an empty report above a line reading "1 info".
func TestSummaryCountsOnlyVisibleFindings(t *testing.T) {
	for _, format := range []string{flags.FormatTable, flags.FormatCompact} {
		out := render(t, Options{Format: format, Level: "error"})

		if strings.Contains(out, "1 info") || strings.Contains(out, "1 warnings") {
			t.Errorf("%s: summary counted findings the level filter hid:\n%s", format, out)
		}
		if !strings.Contains(out, "1 errors") {
			t.Errorf("%s: summary should still count the visible error:\n%s", format, out)
		}
	}
}

// Compact mode printed every finding but only accounted for errors and warnings,
// so "12 failed (0 errors, 1 warnings)" hid the other 11.
func TestCompactSummaryAccountsForEveryFailure(t *testing.T) {
	out := render(t, Options{Format: flags.FormatCompact})
	if !strings.Contains(out, "1 info") {
		t.Errorf("compact summary should count info failures:\n%s", out)
	}
}

func TestVerboseWidensTheLevelFilter(t *testing.T) {
	out := render(t, Options{Format: flags.FormatCompact, Level: "error", Verbose: true})
	if !strings.Contains(out, "license") {
		t.Error("verbose output should include warnings")
	}
}

func TestTableShowsSkippedOnlyWhenRequested(t *testing.T) {
	if strings.Contains(render(t, Options{Format: flags.FormatTable}), "Not applicable") {
		t.Error("skipped rules must be hidden by default")
	}
	out := render(t, Options{Format: flags.FormatTable, ShowSkipped: true})
	if !strings.Contains(out, "Skipped") {
		t.Error("skipped section should appear when requested")
	}
}

func TestQuietSuppressesTheSummary(t *testing.T) {
	out := render(t, Options{Format: flags.FormatTable, Quiet: true})
	if strings.Contains(out, "Status:") {
		t.Error("quiet mode should suppress the status line")
	}
	if !strings.Contains(out, "tests_folder") {
		t.Error("quiet mode must still report failures")
	}
}

func TestMarkdownEscapesPipes(t *testing.T) {
	res := result()
	res.Results[1].Message = "use a | b"

	opts := Options{Format: flags.FormatMarkdown, NoColor: true}
	r, _ := New(opts)
	var buf bytes.Buffer
	if err := r.Render(&buf, res); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `a \| b`) {
		t.Errorf("pipe was not escaped:\n%s", buf.String())
	}
}
