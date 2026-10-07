package resources

import (
	"fmt"
	"sort"
	"strings"
)

// GitHub Actions template names, grouped by concern.
//
// Each entry addresses one workflow file in embedded/github-actions.yml, so a
// project can adopt exactly the workflows it needs rather than a single
// all-in-one pipeline.
const (
	// Continuous integration.
	TmplCIGo       = "ci_go"
	TmplCINode     = "ci_nodejs"
	TmplCIRust     = "ci_rust"
	TmplCIPython   = "ci_python"
	TmplCIMonorepo = "ci_monorepo"

	// Releases.
	TmplReleaseGo        = "release_go"
	TmplReleaseNode      = "release_nodejs"
	TmplReleaseRust      = "release_rust"
	TmplReleaseLinux     = "release_linux"
	TmplReleaseWindows   = "release_windows"
	TmplReleaseMacOS     = "release_macos"
	TmplReleaseAggregate = "release_aggregate"

	// Containers.
	TmplDockerBuild = "docker_build"
	TmplDockerHub   = "docker_hub"
	TmplDockerGHCR  = "docker_ghcr"

	// Deployment.
	TmplDeployPages     = "deploy_github_pages"
	TmplDeployDocs      = "deploy_docs"
	TmplDeployContainer = "deploy_container"

	// Security.
	TmplCodeQL           = "codeql"
	TmplDependencyReview = "security_dependency_review"
	TmplSecretScan       = "security_secret_scan"
)

// GitHubActionsConfig holds every workflow template, keyed by name.
//
// The embedded file is a flat map of name to workflow body, so the struct
// mirrors that shape rather than nesting under a key.
type GitHubActionsConfig struct {
	Workflows map[string]string `yaml:",inline"`
}

func (c *GitHubActionsConfig) Get(name string) (string, bool) {
	body, ok := c.Workflows[name]
	return body, ok && strings.TrimSpace(body) != ""
}

func (c *GitHubActionsConfig) Names() []string {
	out := make([]string, 0, len(c.Workflows))
	for name := range c.Workflows {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func WorkflowNames() []string {
	if actions == nil {
		return nil
	}
	return actions.Names()
}

func WorkflowExists(name string) bool {
	if actions == nil {
		return false
	}
	_, ok := actions.Get(name)
	return ok
}

// WorkflowGroup describes a category of workflows for discovery commands.
type WorkflowGroup struct {
	Group     string
	Purpose   string
	Templates []string
}

// WorkflowGroups lists the workflows by concern, so `psx rules` and the
// generated documentation can present them sensibly.
func WorkflowGroups() []WorkflowGroup {
	return []WorkflowGroup{
		{"ci", "Validate every push and pull request", []string{
			TmplCIGo, TmplCINode, TmplCIRust, TmplCIPython, TmplCIMonorepo,
		}},
		{"release", "Publish artifacts from a version tag", []string{
			TmplReleaseGo, TmplReleaseNode, TmplReleaseRust,
			TmplReleaseLinux, TmplReleaseWindows, TmplReleaseMacOS,
			TmplReleaseAggregate,
		}},
		{"docker", "Build images and publish them", []string{
			TmplDockerBuild, TmplDockerHub, TmplDockerGHCR,
		}},
		{"deploy", "Ship the project somewhere", []string{
			TmplDeployPages, TmplDeployDocs, TmplDeployContainer,
		}},
		{"security", "Scan for vulnerabilities and leaked secrets", []string{
			TmplCodeQL, TmplDependencyReview, TmplSecretScan,
		}},
	}
}

// Workflow returns a workflow body with project variables substituted.
func Workflow(name, projectType string, vars map[string]string) (string, error) {
	if !WorkflowExists(name) {
		return "", fmt.Errorf("no workflow template named %q (available: %s)",
			name, strings.Join(WorkflowNames(), ", "))
	}
	body, _ := actions.Get(name)
	return replaceVars(body, vars), nil
}

// workflowPaths maps a template name to the file a project should create.
//
// Most templates are standalone workflows; the release aggregate delegates to
// three per-platform workers, so those are emitted alongside it.
var workflowPaths = map[string]string{
	TmplReleaseAggregate: ".github/workflows/release.yml",
	TmplReleaseLinux:     ".github/workflows/release-linux.yml",
	TmplReleaseWindows:   ".github/workflows/release-windows.yml",
	TmplReleaseMacOS:     ".github/workflows/release-macos.yml",
}

// WorkflowPath returns the project-relative path a workflow template belongs at.
//
// Templates that are not in the table are derived from their group.
func WorkflowPath(name string) string {
	if p, ok := workflowPaths[name]; ok {
		return p
	}
	for _, g := range WorkflowGroups() {
		for _, n := range g.Templates {
			if n != name {
				continue
			}
			switch g.Group {
			case "ci":
				return ".github/workflows/ci.yml"
			case "release":
				return ".github/workflows/release.yml"
			case "docker":
				return ".github/workflows/docker.yml"
			case "deploy":
				return ".github/workflows/deploy.yml"
			default:
				return ".github/workflows/" + name + ".yml"
			}
		}
	}
	return ".github/workflows/" + name + ".yml"
}

func workflowPathFor(name string) string { return WorkflowPath(name) }

// PreviewVars returns template variables with neutral placeholder values, for
// previewing a template without a real project.
func PreviewVars() map[string]string {
	info := &ProjectInfo{
		Name:        "project",
		Description: "A project",
		Author:      "Your Name",
		Email:       "you@example.com",
		GitHubUser:  "your-username",
		RepoName:    "project",
		License:     "MIT",
	}
	return info.ToVars()
}

// WorkflowForGroup returns the workflow matching a project type within a group,
// falling back to the group's generic option.
func WorkflowForGroup(group, projectType string, vars map[string]string) (string, error) {
	for _, g := range WorkflowGroups() {
		if g.Group != group {
			continue
		}
		for _, name := range preferredOrder(g.Templates, projectType) {
			if body, err := Workflow(name, projectType, vars); err == nil {
				return body, nil
			}
		}
	}
	return "", fmt.Errorf("no workflow template for group %q", group)
}

// preferredOrder puts the project-type specific template first.
func preferredOrder(names []string, projectType string) []string {
	var langSpecific, rest []string
	for _, name := range names {
		if strings.HasSuffix(name, "_"+projectType) {
			langSpecific = append(langSpecific, name)
			continue
		}
		rest = append(rest, name)
	}
	return append(langSpecific, rest...)
}
