package resources

import (
	"fmt"
	"sort"
	"strings"
)

// Template names addressable from rule metadata.
const (
	TmplLicense         = "license"
	TmplReadme          = "readme"
	TmplChangelog       = "changelog"
	TmplContributing    = "contributing"
	TmplGitignore       = "gitignore"
	TmplAPIDocs         = "api_docs"
	TmplEditorconfig    = "editorconfig"
	TmplPreCommit       = "pre_commit"
	TmplDockerfile      = "dockerfile"
	TmplDockerignore    = "dockerignore"
	TmplDockerCompose   = "docker_compose"
	TmplSecurity        = "security"
	TmplCodeOfConduct   = "code_of_conduct"
	TmplPRTemplate      = "pull_request_template"
	TmplIssueBug        = "issue_bug_report"
	TmplIssueFeature    = "issue_feature_request"
	TmplIssueQuestion   = "issue_question"
	TmplIssueConfig     = "issue_templates_config"
	TmplCodeowners      = "codeowners"
	TmplADRFirst        = "adr_first"
	TmplADRTemplate     = "adr_template"
	TmplDocsIndex       = "docs_index"
	TmplScriptSetup     = "script_setup"
	TmplScriptTest      = "script_test"
	TmplScriptBuild     = "script_build"
	TmplScriptClean     = "script_clean"
	TmplGitHubActions   = "github_actions"
	TmplGitLabCI        = "gitlab_ci"
	TmplEnvExample      = "env_example"
	TmplGitattributes   = "gitattributes"
	TmplDependabot      = "dependabot"
	TmplRenovate        = "renovate"
	TmplFunding         = "funding"
	TmplSupport         = "support"
	TmplRoadmap         = "roadmap"
	TmplArchitecture    = "architecture"
	TmplRunbook         = "runbook"
	TmplOpenAPI         = "openapi"
	TmplReleaseWorkflow = "release_workflow"
	TmplNodeVersion     = "node_version"
	TmplMakefile        = "makefile"
	TmplNginx           = "nginx"
	TmplHelmChart       = "helm_chart"
	TmplKubernetes      = "kubernetes"
)

// questionTemplates are templates rendered from a question answer.
var questionTemplates = map[string]string{
	TmplDockerCompose: "docker_compose",
}

// Template returns the rendered body for name, chosen for projectType.
//
// It is the single lookup path for rule fixes, so a rule never needs a
// hand-written getter.
func Template(name, projectType string, vars map[string]string) (string, error) {
	raw, ok := lookup(name, projectType)
	if !ok {
		return "", fmt.Errorf("no template named %q for project type %q", name, projectType)
	}
	return replaceVars(raw, vars), nil
}

// scriptTemplates are rendered through Generator, which layers the language's
// command variables onto the project variables. They have no static body.
var scriptTemplates = map[string]bool{
	TmplScriptSetup: true, TmplScriptTest: true,
	TmplScriptBuild: true, TmplScriptClean: true,
}

// HasTemplate reports whether name resolves for the generic tier or any language.
//
// Script templates count as present: they resolve through the generator, not
// through a static body.
func HasTemplate(name string) bool {
	if scriptTemplates[name] {
		return projectScripts != nil
	}
	if _, ok := lookup(name, "generic"); ok {
		return true
	}
	for _, lang := range SupportedLanguages() {
		if _, ok := lookup(name, lang); ok {
			return true
		}
	}
	// Variants such as docker_compose.with_db resolve through the dotted path.
	if _, ok := lookupDotted(name, "generic"); ok {
		return true
	}
	return false
}

func IsScriptTemplate(name string) bool { return scriptTemplates[name] }

// TemplateNames lists every addressable template, for `psx rules` and docs.
func TemplateNames() []string {
	names := map[string]struct{}{
		TmplReadme: {}, TmplChangelog: {}, TmplContributing: {}, TmplGitignore: {},
		TmplAPIDocs: {}, TmplEditorconfig: {}, TmplPreCommit: {}, TmplDockerfile: {},
		TmplDockerignore: {}, TmplDockerCompose: {}, TmplSecurity: {},
		TmplCodeOfConduct: {}, TmplPRTemplate: {}, TmplIssueBug: {}, TmplIssueFeature: {},
		TmplIssueQuestion: {}, TmplIssueConfig: {}, TmplCodeowners: {}, TmplADRFirst: {},
		TmplADRTemplate: {}, TmplGitHubActions: {}, TmplGitLabCI: {}, TmplEnvExample: {},
		TmplGitattributes: {}, TmplDependabot: {}, TmplRenovate: {}, TmplFunding: {},
		TmplSupport: {}, TmplRoadmap: {}, TmplArchitecture: {}, TmplRunbook: {},
		TmplOpenAPI: {}, TmplReleaseWorkflow: {}, TmplNodeVersion: {}, TmplMakefile: {},
		TmplNginx: {}, TmplHelmChart: {}, TmplKubernetes: {}, TmplDocsIndex: {},
		TmplScriptSetup: {}, TmplScriptTest: {}, TmplScriptBuild: {}, TmplScriptClean: {},
	}
	out := make([]string, 0, len(names))
	for n := range names {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// lookup resolves a template name to a raw body for the given project type.
//
// The license template is a special case: its body is the selected license
// text with the author substituted, so it is resolved through GetLicense rather
// than the static template maps.
func lookup(name, projectType string) (string, bool) {
	if name == "" {
		return "", false
	}
	if name == TmplLicense {
		body := GetLicense("MIT", "")
		return body, body != ""
	}

	pt := normalizeKey(projectType)
	if pt == "" {
		pt = "generic"
	}

	// The project type is tried first, then generic, then any language tier.
	for _, key := range dedupeKeys(pt) {
		if s, ok := fromLangMaps(name, key); ok && s != "" {
			return s, true
		}
	}
	if s, ok := fromScalars(name); ok {
		return s, true
	}
	if s, ok := lookupDotted(name, pt); ok {
		return s, true
	}
	// Workflow templates live in their own library but resolve through the same
	// path, so a rule can name one like any other template.
	if s, ok := actions.Get(name); ok {
		return s, true
	}
	return "", false
}

func dedupeKeys(pt string) []string {
	keys := []string{pt, "generic"}
	if pt != "generic" {
		keys = append(keys, "nodejs", "go")
	}
	out := keys[:0]
	seen := map[string]bool{}
	for _, k := range keys {
		if !seen[k] {
			seen[k] = true
			out = append(out, k)
		}
	}
	return out
}

func fromLangMaps(name, key string) (string, bool) {
	var body string
	switch name {
	case TmplReadme:
		body = templates.Readme[key]
	case TmplAPIDocs:
		body = templates.APIDocs[key]
	case TmplEditorconfig:
		body = qualityTools.Editorconfig[key]
	case TmplPreCommit:
		body = qualityTools.PreCommit[key]
	case TmplGitHubActions:
		body = devops.CICD.GitHubActions[key]
	case TmplGitLabCI:
		body = devops.CICD.GitLabCI[key]
	case TmplDependabot:
		body = devops.Dependabot[key]
	case TmplRenovate:
		body = devops.Renovate[key]
	case TmplNodeVersion:
		body = devops.NodeVersion[key]
	case TmplReleaseWorkflow:
		body = devops.ReleaseWorkflow[key]
	case TmplNginx:
		body = devops.Nginx[key]
	case TmplHelmChart:
		body = devops.HelmChart[key]
	case TmplKubernetes:
		body = devops.Kubernetes[key]
	case TmplEnvExample:
		body = gitignores.EnvExample[key]
	case TmplMakefile:
		body = qualityTools.Makefile[key]
	case TmplGitignore:
		// The language block is layered onto the common block by the generator.
		switch key {
		case "nodejs":
			return gitignores.NodeJS, true
		case "go":
			return gitignores.Go, true
		default:
			return gitignores.Common, true
		}
	case TmplDockerfile, TmplDockerignore:
		return dockerTemplate(name, key)
	case TmplDockerCompose:
		// The bare name resolves to the base variant; with_database and other
		// variants are selected by the generator from a question answer.
		if s := devops.DockerCompose[key]; s != "" {
			return s, true
		}
		if s := devops.DockerCompose["basic"]; s != "" {
			return s, true
		}
	}
	return body, body != ""
}

func fromScalars(name string) (string, bool) {
	switch name {
	case TmplChangelog:
		return templates.Changelog, true
	case TmplContributing:
		return templates.Contributing, true
	case TmplSecurity:
		return docsTemplates.Security, true
	case TmplCodeOfConduct:
		return docsTemplates.CodeOfConduct, true
	case TmplPRTemplate:
		return docsTemplates.PullRequestTemplate, true
	case TmplIssueBug:
		return docsTemplates.IssueBugReport, true
	case TmplIssueFeature:
		return docsTemplates.IssueFeatureRequest, true
	case TmplIssueQuestion:
		return docsTemplates.IssueQuestion, true
	case TmplIssueConfig:
		return docsTemplates.IssueTemplatesConfig, true
	case TmplCodeowners:
		return docsTemplates.Codeowners, true
	case TmplADRFirst:
		return docsTemplates.ADRTemplates["first"], true
	case TmplADRTemplate:
		return docsTemplates.ADRTemplates["template"], true
	case TmplDocsIndex:
		return templates.DocsIndex, true
	case TmplScriptSetup, TmplScriptTest, TmplScriptBuild, TmplScriptClean:
		// Script bodies carry command variables the language profile supplies,
		// so they are only renderable through Generator.Render.
		return "", false
	case TmplGitattributes:
		return qualityTools.Gitattributes, true
	case TmplFunding:
		return docsTemplates.Funding, true
	case TmplSupport:
		return docsTemplates.Support, true
	case TmplRoadmap:
		return docsTemplates.Roadmap, true
	case TmplArchitecture:
		return docsTemplates.Architecture, true
	case TmplRunbook:
		return docsTemplates.Runbook, true
	case TmplOpenAPI:
		return docsTemplates.OpenAPI, true
	}
	return "", false
}

func dockerTemplate(name, projectType string) (string, bool) {
	cfg := devops.Docker
	key := normalizeKey(projectType)
	var lang DockerLanguageConfig
	switch key {
	case "nodejs":
		lang = cfg.NodeJS
	case "go":
		lang = cfg.Go
	default:
		lang = cfg.Generic
	}
	if name == TmplDockerfile {
		return lang.Dockerfile, true
	}
	return lang.Dockerignore, true
}

// languageTables lists the map-backed templates, exposed so tests and
// documentation can enumerate them without duplicating the switch.
func languageTables() map[string]map[string]string {
	return map[string]map[string]string{
		TmplReadme:          templates.Readme,
		TmplAPIDocs:         templates.APIDocs,
		TmplEditorconfig:    qualityTools.Editorconfig,
		TmplPreCommit:       qualityTools.PreCommit,
		TmplMakefile:        qualityTools.Makefile,
		TmplGitHubActions:   devops.CICD.GitHubActions,
		TmplGitLabCI:        devops.CICD.GitLabCI,
		TmplDockerCompose:   devops.DockerCompose,
		TmplDependabot:      devops.Dependabot,
		TmplRenovate:        devops.Renovate,
		TmplNodeVersion:     devops.NodeVersion,
		TmplReleaseWorkflow: devops.ReleaseWorkflow,
		TmplNginx:           devops.Nginx,
		TmplHelmChart:       devops.HelmChart,
		TmplKubernetes:      devops.Kubernetes,
		TmplEnvExample:      gitignores.EnvExample,
	}
}

func normalizeKey(k string) string { return strings.ToLower(strings.TrimSpace(k)) }
