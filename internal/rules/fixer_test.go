package rules

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

// fixEnv is a staged project plus everything needed to run a fix.
type fixEnv struct {
	root string
	cfg  *config.Config
	ctx  *Context
	opts resources.Options
}

// newFixEnv stages a project and builds a fixer context for it.
func newFixEnv(t *testing.T, files map[string]string, projectType string, ruleIDs []string) *fixEnv {
	t.Helper()
	root := stage(t, files)

	cfg := testConfig(t, ruleIDs, nil)
	// project_type comes from the same declared value a real config would set.
	cfg.Project.Type = projectType

	info := &resources.ProjectInfo{
		Name:        "demo",
		Description: "A demo project",
		Author:      "Ada Lovelace",
		Email:       "ada@example.com",
		GitHubUser:  "adal",
		RepoName:    "demo",
		License:     "MIT",
	}

	return &fixEnv{
		root: root,
		cfg:  cfg,
		ctx:  &Context{ProjectPath: root, ProjectType: projectType, ProjectInfo: info, Config: cfg},
		opts: resources.Options{ProjectType: projectType, Answers: map[string]string{}},
	}
}

// run applies a fix for the given rules.
func (e *fixEnv) run(t *testing.T, ruleIDs []string, dryRun, force bool) []*FixResult {
	t.Helper()
	return FixAll(e.cfg, &FixContext{
		Context: e.ctx,
		DryRun:  dryRun,
		Force:   force,
	}, ruleIDs, e.opts)
}

// read returns the contents of a project-relative path.
func (e *fixEnv) read(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(e.root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

func TestFixCreatesDeclaredPath(t *testing.T) {
	env := newFixEnv(t, map[string]string{}, "go", []string{"readme"})
	results := env.run(t, []string{"readme"}, false, false)

	if len(results) != 1 || !results[0].Fixed {
		t.Fatalf("expected the readme rule to be fixed: %+v", results)
	}
	if _, err := os.Stat(filepath.Join(env.root, "README.md")); err != nil {
		t.Fatalf("README.md was not created: %v", err)
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	env := newFixEnv(t, map[string]string{}, "go", []string{"readme", "license", "changelog"})
	env.run(t, []string{"readme", "license", "changelog"}, true, false)

	entries, err := os.ReadDir(env.root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("dry run created %d entries: %v", len(entries), names(entries))
	}
}

func TestFixNeverOverwritesContentWithoutForce(t *testing.T) {
	env := newFixEnv(t, map[string]string{
		"README.md": "# my own readme, do not touch",
	}, "go", []string{"readme"})

	results := env.run(t, []string{"readme"}, false, false)
	if len(results) != 1 || !results[0].Skipped {
		t.Fatalf("expected a skip: %+v", results)
	}
	if got := env.read(t, "README.md"); got != "# my own readme, do not touch" {
		t.Errorf("existing content was modified: %q", got)
	}
}

func TestForceOverwritesContent(t *testing.T) {
	env := newFixEnv(t, map[string]string{
		"README.md": "old",
	}, "go", []string{"readme"})

	results := env.run(t, []string{"readme"}, false, true)
	if len(results) != 1 || !results[0].Fixed {
		t.Fatalf("force should have overwritten: %+v", results)
	}
	if got := env.read(t, "README.md"); !strings.Contains(got, "demo") {
		t.Errorf("README was not regenerated: %q", got)
	}
}

func TestFixIsIdempotent(t *testing.T) {
	// Running fix twice must not report work the second time, otherwise
	// `psx fix` is never actually finished.
	env := newFixEnv(t, map[string]string{}, "nodejs",
		[]string{"readme", "license", "gitignore", "editorconfig", "changelog"})

	env.run(t, []string{"readme", "license", "gitignore", "editorconfig", "changelog"}, false, false)
	second := env.run(t, []string{"readme", "license", "gitignore", "editorconfig", "changelog"}, false, false)

	for _, r := range second {
		if r.Fixed {
			t.Errorf("rule %s reported a fix on a second run", r.RuleID)
		}
		if !r.Skipped {
			t.Errorf("rule %s should be skipped on a second run, got %+v", r.RuleID, r)
		}
	}
}

func TestGeneratedContentHasNoRawPlaceholders(t *testing.T) {
	// Every fixable rule is applied to an empty project and every produced file
	// is scanned for leftover {{placeholders}}.
	env := newFixEnv(t, map[string]string{}, "go", nil)

	allIDs := make([]string, 0)
	meta := config.GetRulesMetadata().Rules
	for id, m := range meta {
		if m.Fixable() {
			allIDs = append(allIDs, id)
		}
	}

	env.cfg = testConfig(t, allIDs, nil)
	env.ctx.Config = env.cfg
	// Dry run is the only mode that carries the generated body, so it is what
	// the placeholder scan inspects.
	results := env.run(t, allIDs, true, false)

	placeholderPattern := regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`)
	checked := 0

	for _, r := range results {
		if r.Error != nil {
			t.Errorf("rule %s failed: %v", r.RuleID, r.Error)
			continue
		}
		for _, c := range r.Changes {
			if c.Content == "" {
				continue
			}
			checked++
			for _, m := range placeholderPattern.FindAllStringSubmatch(c.Content, -1) {
				if resources.IsHumanPlaceholder(m[1]) {
					continue
				}
				t.Errorf("%s: generated content contains %s", c.Path, m[0])
			}
		}
	}
	if checked == 0 {
		t.Fatal("no content was inspected; the test is not exercising the generator")
	}
}

func TestEveryFixableRuleSucceeds(t *testing.T) {
	// A rule that declares a fix must be able to run it. This catches typos in
	// template names, which otherwise surface only at fix time.
	env := newFixEnv(t, map[string]string{}, "generic", nil)

	meta := config.GetRulesMetadata().Rules
	var ids []string
	for id, m := range meta {
		if m.Fixable() {
			ids = append(ids, id)
		}
	}

	env.cfg = testConfig(t, ids, nil)
	env.ctx.Config = env.cfg
	for _, r := range env.run(t, ids, false, false) {
		if r.Error != nil {
			t.Errorf("rule %s: %v", r.RuleID, r.Error)
		}
		if r.Skipped && r.Reason == "no automatic fix available" {
			t.Errorf("rule %s declares a fix but reported none", r.RuleID)
		}
	}
}

func TestFixPathIsNeverAGlob(t *testing.T) {
	// A glob in a fix path would create a literal "**" directory or a file named
	// "*_test.go", which is the class of bug this rule set previously had.
	meta := config.GetRulesMetadata().Rules
	for id, m := range meta {
		if !m.Fixable() {
			continue
		}
		if strings.ContainsAny(m.Fix.Path, "*?[]") {
			t.Errorf("rule %s: fix.path %q contains glob metacharacters", id, m.Fix.Path)
		}
		for path := range m.Fix.Files {
			if strings.ContainsAny(path, "*?[]") {
				t.Errorf("rule %s: fix file %q contains glob metacharacters", id, path)
			}
		}
	}
}

func TestEveryFixTemplateExists(t *testing.T) {
	meta := config.GetRulesMetadata().Rules
	for id, m := range meta {
		if !m.Fixable() {
			continue
		}
		for _, name := range templatesFor(m) {
			if !resources.HasTemplate(name) {
				t.Errorf("rule %s references unknown template %q", id, name)
			}
		}
	}
}

// templatesFor lists every template a rule's fix depends on.
func templatesFor(m config.RuleMetadata) []string {
	var out []string
	if m.Fix.Template != "" {
		out = append(out, m.Fix.Template)
	}
	for _, name := range m.Fix.Files {
		out = append(out, name)
	}
	return out
}

func TestScriptsAreExecutable(t *testing.T) {
	if os.PathSeparator == '\\' {
		t.Skip("unix permissions are not meaningful on windows")
	}

	env := newFixEnv(t, map[string]string{}, "go", []string{"scripts_folder"})
	for _, r := range env.run(t, []string{"scripts_folder"}, false, false) {
		for _, c := range r.Changes {
			info, err := os.Stat(c.Path)
			if err != nil {
				t.Fatalf("stat %s: %v", c.Path, err)
			}
			if info.Mode().Perm()&0o111 == 0 {
				t.Errorf("%s is not executable (mode %v)", c.Path, info.Mode().Perm())
			}
		}
	}
}

func TestCustomPathsCannotEscapeTheProject(t *testing.T) {
	// A configuration is untrusted input; custom paths must stay inside the
	// project directory.
	env := newFixEnv(t, map[string]string{}, "generic", nil)

	env.cfg.Custom = &config.CustomConfig{
		Files: []config.CustomFile{
			{Path: "../../escape.txt", Content: "nope"},
			{Path: "/absolute.txt", Content: "nope"},
			{Path: "ok.txt", Content: "fine"},
		},
	}

	results := CustomHandler{root: env.root}.Apply(env.cfg.Custom, &FixContext{Context: env.ctx})

	for _, r := range results {
		if r.Error == nil && strings.Contains(firstPath(r), "escape") {
			t.Errorf("custom fix escaped the project: %s", firstPath(r))
		}
	}
	if _, err := os.Stat(filepath.Join(env.root, "ok.txt")); err != nil {
		t.Errorf("a safe custom path should have been created: %v", err)
	}
}

func TestUnknownRuleIsReportedNotIgnored(t *testing.T) {
	env := newFixEnv(t, map[string]string{}, "go", []string{"readme"})
	results := env.run(t, []string{"does_not_exist"}, false, false)

	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}
	if results[0].Error == nil {
		t.Error("an unknown rule must produce an error result, not be skipped silently")
	}
}

// names extracts file names for error messages.
func names(entries []os.DirEntry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name())
	}
	return out
}

// firstPath returns the first path a fix result touched.
func firstPath(r *FixResult) string {
	if len(r.Changes) == 0 {
		return ""
	}
	return r.Changes[0].Path
}
