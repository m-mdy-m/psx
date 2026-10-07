package resources

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Options carry the resolved answers a fix needs to render its templates.
type Options struct {
	ProjectType string
	Answers     map[string]string
	Force       bool
}

func (o Options) Answer(id string) string {
	if o.Answers == nil {
		return ""
	}
	return strings.TrimSpace(o.Answers[id])
}

func (o Options) AnswerOr(id, fallback string) string {
	if v := o.Answer(id); v != "" {
		return v
	}
	return fallback
}

// Generator renders rule fix specs into concrete file bodies.
//
// It replaces the per-rule switch statements that previously lived in the
// rules package: rule metadata names a template, this type resolves it.
type Generator struct {
	info    *ProjectInfo
	project string
}

// NewGenerator binds a generator to a project and its type.
func NewGenerator(info *ProjectInfo, projectType string) *Generator {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	return &Generator{info: info, project: projectType}
}

func (g *Generator) Vars() map[string]string { return g.info.ToVars() }

func (g *Generator) Project() *ProjectInfo { return g.info }

func (g *Generator) ProjectType() string { return g.project }

// resolveTemplate maps a declared template name to the concrete variant an
// answer selects, so "ci_config" plus Answer("ci_platform") == "gitlab"
// resolves to the GitLab body rather than the GitHub default.
func (g *Generator) resolveTemplate(name string, opts Options) (string, Options, error) {
	switch name {
	case TmplGitHubActions:
		return TmplGitHubActions, opts, nil

	case TmplGitLabCI:
		return TmplGitLabCI, opts, nil

	case TmplDockerCompose:
		// with_database selects the compose variant.
		if opts.AnswerOr("with_database", "no") == "yes" {
			return "docker_compose.with_db", opts, nil
		}
		return "docker_compose.basic", opts, nil

	case TmplGitignore:
		// Language-specific ignore bodies are merged with the common block.
		return TmplGitignore, opts, nil
	}
	return name, opts, nil
}

func (g *Generator) licenseBody() string {
	return GetLicense(g.info.License, g.info.Author)
}

// Render produces the body for a template name, honouring the supplied answers.
func (g *Generator) Render(name string, opts Options) (string, error) {
	resolved, opts, err := g.resolveTemplate(name, opts)
	if err != nil {
		return "", err
	}

	pt := normalizeKey(opts.ProjectType)
	if pt == "" {
		pt = normalizeKey(g.project)
	}

	// The gitignore fix composes the common block with the language block.
	if resolved == TmplGitignore {
		body, err := Template(TmplGitignore, pt, g.Vars())
		if err != nil {
			return "", err
		}
		return body, nil
	}

	// The license body comes from the project's own license choice.
	if resolved == TmplLicense {
		return g.licenseBody(), nil
	}

	body, err := Template(resolved, pt, g.Vars())
	if err != nil {
		// Fall back to a dotted path such as "docker_compose.with_db".
		if s, ok := lookupDotted(resolved, pt); ok {
			return replaceVars(s, g.Vars()), nil
		}
		if strings.TrimSpace(pt) != "" {
			return "", NotApplicable(name, pt)
		}
		return "", err
	}
	if strings.TrimSpace(body) == "" {
		return "", fmt.Errorf("template %q resolved to an empty body for %q", name, pt)
	}
	return body, nil
}

// lookupDotted resolves "group.variant" against the embedded config maps.
func lookupDotted(dotted, projectType string) (string, bool) {
	group, variant, ok := strings.Cut(dotted, ".")
	if !ok {
		return "", false
	}
	var body string
	switch group {
	case "docker_compose":
		body = devops.DockerCompose[variant]
	case "github_actions":
		body = devops.CICD.GitHubActions[variant]
	case "gitlab_ci":
		body = devops.CICD.GitLabCI[variant]
	case "dependabot":
		body = devops.Dependabot[variant]
	case "renovate":
		body = devops.Renovate[variant]
	default:
		return "", false
	}
	if body == "" {
		return "", false
	}
	if pt := normalizeKey(projectType); pt != "" && pt != "generic" {
		// Prefer the project-type tier when the variant map is language keyed.
		if typed, ok := lookupDottedVariant(group, variant, pt); ok {
			return typed, true
		}
	}
	return body, true
}

func lookupDottedVariant(group, variant, projectType string) (string, bool) {
	var table map[string]string
	switch group {
	case "docker_compose":
		table = devops.DockerCompose
	case "github_actions":
		table = devops.CICD.GitHubActions
	case "gitlab_ci":
		table = devops.CICD.GitLabCI
	case "dependabot":
		table = devops.Dependabot
	case "renovate":
		table = devops.Renovate
	default:
		return "", false
	}
	body := table[variant]
	return body, body != ""
}

// Scripts renders the developer script set for a project.
//
// The script bodies are language-aware: build, test and install commands come
// from the language profile rather than being left as unresolved placeholders.
func (g *Generator) Scripts(opts Options) map[string]string {
	lang := g.language()

	out := make(map[string]string, len(devScripts))
	for _, s := range devScripts {
		body := g.scriptFor(s.key)
		if strings.TrimSpace(body) == "" {
			continue
		}
		out[s.file] = replaceVars(replaceVars(body, g.Vars()), lang)
	}
	return out
}

type scriptSpec struct {
	key  string
	file string
	mode uint32
}

// Keys match the script_setup / script_test / script_build / script_clean
// template names used by rule metadata.
var devScripts = []scriptSpec{
	{"setup", "setup.sh", 0o755},
	{"test", "test.sh", 0o755},
	{"build", "build.sh", 0o755},
	{"clean", "clean.sh", 0o755},
}

func ScriptNames() []string {
	out := make([]string, 0, len(devScripts))
	for _, s := range devScripts {
		out = append(out, s.file)
	}
	sort.Strings(out)
	return out
}

func ScriptMode(name string) uint32 {
	for _, s := range devScripts {
		if s.file == name {
			return s.mode
		}
	}
	return 0o755
}

func (g *Generator) scriptFor(key string) string {
	switch key {
	case "setup":
		return projectScripts.Setup
	case "test":
		return projectScripts.Test
	case "build":
		return projectScripts.Build
	case "clean":
		return projectScripts.Clean
	default:
		return ""
	}
}

func (g *Generator) language() map[string]string {
	key := normalizeKey(g.project)
	lang, ok := languages.Languages[key]
	if !ok {
		return nil
	}
	return map[string]string{
		"install_command":        firstOr(lang.InstallCommands, ""),
		"test_command":           firstOr(lang.TestCommands, ""),
		"build_command":          firstOr(lang.BuildCommands, ""),
		"start_command":          firstOr(lang.StartCommands, ""),
		"build_dir":              firstOr(lang.BuildPatterns, "dist"),
		"build_dirs":             joinOr(lang.BuildPatterns, "build dist"),
		"cache_dirs":             joinOr(lang.CachePatterns, ""),
		"version_update_command": versionUpdateFor(key),
	}
}

func versionUpdateFor(lang string) string {
	switch lang {
	case "nodejs":
		return "npm version $VERSION --no-git-tag-version"
	case "go":
		return "echo \"Version comes from git tags (git describe)\""
	default:
		return "echo \"No version file to update\""
	}
}

func firstOr(list []string, fallback string) string {
	if len(list) > 0 {
		return list[0]
	}
	return fallback
}

func joinOr(list []string, fallback string) string {
	if len(list) == 0 {
		return fallback
	}
	return strings.Join(list, " ")
}

func ModeFromName(name string) uint32 {
	switch filepath.Ext(name) {
	case ".sh", ".bash", ".ps1", ".bat", ".cmd":
		return 0o755
	default:
		return 0o644
	}
}

// notApplicableError marks a template that does not exist for a project type.
//
// It is distinct from a real failure: the rule is simply irrelevant here, so the
// caller should skip it rather than report an error.
type notApplicableError struct {
	name        string
	projectType string
}

func (e *notApplicableError) Error() string {
	return fmt.Sprintf("no %q template for %s projects", e.name, e.projectType)
}

// humanPlaceholders are template tokens a person fills in after generation,
// such as the ADR number and title. psx must leave them untouched.
var humanPlaceholders = map[string]bool{
	"number": true, "title": true, "status": true, "date": true,
}

// IsHumanPlaceholder reports whether a template token is meant to be filled in
// by a person rather than substituted by psx.
func IsHumanPlaceholder(name string) bool { return humanPlaceholders[name] }

// NotApplicable wraps err as a not-applicable template error.
func NotApplicable(name, projectType string) error {
	return &notApplicableError{name: name, projectType: projectType}
}

// IsNotApplicable reports whether err means "this template does not apply here".
func IsNotApplicable(err error) bool {
	var target *notApplicableError
	return errors.As(err, &target)
}

func Message(category, key string, args ...any) string {
	msg := GetMessage(category, key)
	if msg == "" {
		return ""
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}
