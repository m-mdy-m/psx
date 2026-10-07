package resources

import (
	"regexp"
	"strings"
	"testing"
)

// placeholderRe finds every {{name}} token in a template body.
var placeholderRe = regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`)

// renderableVars are intentionally left in templates for humans to fill in.
var renderableVars = map[string]bool{}

func init() {
	for name := range humanPlaceholders {
		renderableVars[name] = true
	}
}

// projectVars is the full variable set a real render receives.
// A template may only reference keys from this set, otherwise a generated file
// ships with raw {{placeholder}} text.
func projectVars() map[string]string {
	info := &ProjectInfo{
		Name:        "demo",
		Description: "A demo project",
		Author:      "Ada Lovelace",
		Email:       "ada@example.com",
		GitHubUser:  "ada",
		RepoName:    "demo",
		License:     "MIT",
	}
	info.buildDerived()
	return info.ToVars()
}

// TestEveryTemplateRenders guards against the class of bug where a template
// references a variable the renderer never supplies, leaving {{raw}} text in
// generated files.
func TestEveryTemplateRenders(t *testing.T) {
	vars := projectVars()
	for _, name := range TemplateNames() {
		if scriptTemplateNames[name] {
			continue // covered by TestEveryScriptRenders
		}
		for _, pt := range []string{"generic", "nodejs", "go"} {
			body, err := Template(name, pt, vars)
			if err != nil {
				continue // template not defined for this type
			}
			assertNoRawPlaceholders(t, name+"/"+pt, body)
		}
	}
}

// TestEveryScriptRenders is the same guard for developer scripts, which carry
// extra command variables supplied by the language profile.
func TestEveryScriptRenders(t *testing.T) {
	g := NewGenerator(&ProjectInfo{
		Name: "demo", Description: "demo", Author: "Ada",
		Email: "ada@example.com", GitHubUser: "ada", RepoName: "demo", License: "MIT",
	}, "go")

	for _, lang := range []string{"go", "nodejs"} {
		gg := NewGenerator(g.Project(), lang)
		for path, body := range gg.Scripts(Options{}) {
			assertNoRawPlaceholders(t, lang+"/"+path, body)
			if strings.TrimSpace(body) == "" {
				t.Errorf("%s/%s rendered an empty script", lang, path)
			}
		}
	}
}

// TestKnownTemplatesExist pins the templates rules depend on, so a rename in the
// YAML cannot silently break every fix.
func TestKnownTemplatesExist(t *testing.T) {
	required := []string{
		TmplReadme, TmplChangelog, TmplContributing, TmplGitignore,
		TmplAPIDocs, TmplEditorconfig, TmplPreCommit, TmplDockerfile,
		TmplDockerignore, TmplDockerCompose, TmplSecurity, TmplCodeOfConduct,
		TmplPRTemplate, TmplIssueBug, TmplIssueFeature, TmplIssueQuestion,
		TmplIssueConfig, TmplCodeowners, TmplADRFirst, TmplADRTemplate,
		TmplGitHubActions, TmplGitLabCI, TmplDependabot, TmplGitattributes,
		TmplEnvExample, TmplMakefile, TmplFunding, TmplSupport, TmplRoadmap,
		TmplArchitecture, TmplRunbook, TmplOpenAPI, TmplNodeVersion,
		TmplReleaseWorkflow, TmplRenovate, TmplNginx, TmplHelmChart,
		TmplKubernetes,
	}
	for _, name := range required {
		if !HasTemplate(name) {
			t.Errorf("required template %q is not defined in embedded YAML", name)
		}
	}
}

// TestPerLanguageTemplatesExist verifies language-keyed templates resolve for
// each supported language.
func TestPerLanguageTemplatesExist(t *testing.T) {
	vars := map[string]string{}
	for _, pt := range []string{"nodejs", "go"} {
		for _, name := range []string{TmplReadme, TmplDockerfile, TmplGitHubActions, TmplGitLabCI} {
			body, err := Template(name, pt, vars)
			if err != nil {
				t.Errorf("%s/%s: %v", pt, name, err)
				continue
			}
			if strings.TrimSpace(body) == "" {
				t.Errorf("%s/%s resolved to an empty body", pt, name)
			}
		}
	}
}

// TestUnknownTemplateIsAnError keeps a typo from producing an empty file.
func TestUnknownTemplateIsAnError(t *testing.T) {
	if _, err := Template("no_such_template", "generic", nil); err == nil {
		t.Error("expected an error for an unknown template")
	}
}

// scriptTemplateNames are rendered through Generator, which layers language
// command variables on top of the project variables. They have no static body,
// so TestEveryTemplateRenders skips them in favour of TestEveryScriptRenders.
var scriptTemplateNames = map[string]bool{
	TmplScriptSetup: true, TmplScriptTest: true,
	TmplScriptBuild: true, TmplScriptClean: true,
}

// assertNoRawPlaceholders fails on any unresolved {{token}} in body.
func assertNoRawPlaceholders(t *testing.T, label, body string) {
	t.Helper()
	assertNoRawPlaceholdersWith(t, label, body, nil)
}

// assertNoRawPlaceholdersWith is assertNoRawPlaceholders with extra tolerated
// variable names, for templates that legitimately reference a foreign runtime.
func assertNoRawPlaceholdersWith(t *testing.T, label, body string, extra map[string]bool) {
	t.Helper()
	for _, m := range placeholderRe.FindAllStringSubmatch(body, -1) {
		if renderableVars[m[1]] || extra[m[1]] {
			continue
		}
		t.Errorf("%s: unresolved placeholder {{%s}}", label, m[1])
	}
}
