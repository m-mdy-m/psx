package resources

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

// TestWorkflowTemplatesAreValidYAML parses every workflow template.
//
// A GitHub Actions file that is not valid YAML fails at the worst possible
// moment: on a push, in someone else's repository. Catching it here keeps the
// templates usable.
func TestWorkflowTemplatesAreValidYAML(t *testing.T) {
	if actions == nil {
		t.Fatal("github-actions.yml was not loaded")
	}

	for _, name := range actions.Names() {
		t.Run(name, func(t *testing.T) {
			body, ok := actions.Get(name)
			if !ok {
				t.Fatalf("workflow %q is empty", name)
			}

			// GitHub's own parser treats a bare `on:` as the boolean true,
			// so the key arrives as "true"; accept either spelling.
			var doc map[string]any
			if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
				t.Fatalf("workflow %q is not valid YAML: %v", name, err)
			}

			if len(doc) == 0 {
				t.Fatalf("workflow %q parsed to an empty document", name)
			}
			if _, ok := doc["name"]; !ok {
				t.Error("workflow is missing a top-level name")
			}
			if _, ok := doc["jobs"]; !ok {
				t.Error("workflow is missing a jobs section")
			}
		})
	}
}

// TestWorkflowTemplatesUseOnlyKnownVariables guards against shipping files that
// contain raw {{placeholders}} after rendering.
func TestWorkflowTemplatesUseOnlyKnownVariables(t *testing.T) {
	vars := projectVars()

	// Vars that belong to the GitHub Actions runtime, not to psx. The
	// metadata-action semver tokens use the same brace syntax but are expanded
	// by that action, so they are tolerated too.
	githubVars := map[string]bool{
		"github": true, "runner": true, "secrets": true, "steps": true,
		"jobs": true, "matrix": true, "env": true, "inputs": true, "needs": true,
		"version": true, "major": true, "minor": true, "sha": true,
	}

	for _, name := range actions.Names() {
		body, err := Workflow(name, "generic", vars)
		if err != nil {
			t.Fatalf("Workflow(%s): %v", name, err)
		}
		assertNoRawPlaceholdersWith(t, "workflow/"+name, body, githubVars)
	}
}

// TestWorkflowsRenderPerLanguage verifies the language-specific CI templates
// exist and render for each supported language.
func TestWorkflowsRenderPerLanguage(t *testing.T) {
	vars := projectVars()
	cases := map[string]string{
		"go":     TmplCIGo,
		"nodejs": TmplCINode,
		"rust":   TmplCIRust,
		"python": TmplCIPython,
	}

	for lang, name := range cases {
		t.Run(lang, func(t *testing.T) {
			if !WorkflowExists(name) {
				t.Fatalf("missing CI template for %s (%s)", lang, name)
			}
			body, err := Workflow(name, lang, vars)
			if err != nil {
				t.Fatalf("render %s: %v", name, err)
			}
			if strings.TrimSpace(body) == "" {
				t.Fatal("rendered an empty workflow")
			}
		})
	}
}

// TestEveryWorkflowGroupResolves makes sure a group always yields a workflow,
// so a rule can point at a group without risking a generation failure.
func TestEveryWorkflowGroupResolves(t *testing.T) {
	vars := projectVars()
	for _, g := range WorkflowGroups() {
		t.Run(g.Group, func(t *testing.T) {
			body, err := WorkflowForGroup(g.Group, "generic", vars)
			if err != nil {
				t.Fatalf("group %q has no usable template: %v", g.Group, err)
			}
			if strings.TrimSpace(body) == "" {
				t.Error("group resolved to an empty workflow")
			}
		})
	}
}

// TestWorkflowForGroupPrefersLanguageMatch verifies the language-specific
// template wins over the generic one.
func TestWorkflowForGroupPrefersLanguageMatch(t *testing.T) {
	vars := projectVars()

	goBody, err := WorkflowForGroup("ci", "go", vars)
	if err != nil {
		t.Fatal(err)
	}
	nodeBody, err := WorkflowForGroup("ci", "nodejs", vars)
	if err != nil {
		t.Fatal(err)
	}
	if goBody == nodeBody {
		t.Error("the go and nodejs CI templates resolved to the same body")
	}
	if !strings.Contains(goBody, "go test") {
		t.Error("the go CI template should run go test")
	}
	if !strings.Contains(nodeBody, "pnpm") {
		t.Error("the nodejs CI template should use the package manager")
	}
}

// TestUnknownWorkflowIsAnError keeps a typo from producing an empty file.
func TestUnknownWorkflowIsAnError(t *testing.T) {
	if _, err := Workflow("no_such_workflow", "generic", nil); err == nil {
		t.Error("expected an error for an unknown workflow")
	}
}

// TestReleaseWorkflowsAreTagGated is a correctness check on the trigger: a
// release workflow that runs on a branch would publish on every push.
func TestReleaseWorkflowsAreTagGated(t *testing.T) {
	vars := projectVars()

	for _, name := range []string{
		TmplReleaseGo, TmplReleaseNode, TmplReleaseRust,
		TmplReleaseLinux, TmplReleaseWindows, TmplReleaseMacOS,
		TmplReleaseAggregate, TmplDockerHub,
	} {
		t.Run(name, func(t *testing.T) {
			body, err := Workflow(name, "generic", vars)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(body, "tags:") {
				t.Errorf("%s does not gate on a tag; it would publish on every push", name)
			}
		})
	}
}

// TestDockerHubWorkflowUsesDigestPush verifies the multi-arch publish pattern:
// each platform pushes by digest and tags are applied once at the end.
func TestDockerHubWorkflowUsesDigestPush(t *testing.T) {
	body, err := Workflow(TmplDockerHub, "generic", projectVars())
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"push-by-digest=true", "imagetools create", "type=semver"} {
		if !strings.Contains(body, want) {
			t.Errorf("docker_hub workflow is missing %q", want)
		}
	}
}
