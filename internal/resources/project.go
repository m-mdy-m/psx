package resources

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/ui"
)

const projectCacheFile = ".psx-project.yml"

// ResolveInfo returns project metadata for a directory.
//
// Metadata comes from the cache when present, then from git and the
// environment. When allowPrompt is set the user is asked to fill in the gaps and
// the result is cached; otherwise the project directory is left untouched, which
// is what read-only commands require.
func ResolveInfo(projectPath string, allowPrompt bool) (*ProjectInfo, error) {
	if info, err := loadProjectInfo(projectPath); err == nil && info != nil {
		logger.Verbose("Using cached project info")
		info.applyDefaults()
		return info, nil
	}

	info := &ProjectInfo{
		Name:    filepath.Base(projectPath),
		License: "MIT",
	}
	info.loadFromGit()
	info.applyDefaults()

	if !allowPrompt || !ui.IsInteractive() {
		logger.Verbose("Inferring project info without prompting")
		return info, nil
	}

	logger.Step("Project metadata")
	info.promptUser()
	info.applyDefaults()

	if err := saveProjectInfo(projectPath, info); err != nil {
		logger.Warning(fmt.Sprintf("Could not cache project info: %v", err))
	}
	return info, nil
}

// applyDefaults fills blanks from git data and the environment.
func (p *ProjectInfo) applyDefaults() {
	p.setDefaults()
	p.buildDerived()
}

func loadProjectInfo(projectPath string) (*ProjectInfo, error) {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var info ProjectInfo
	if err := yaml.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	info.buildDerived()
	return &info, nil
}

func saveProjectInfo(projectPath string, info *ProjectInfo) error {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := yaml.Marshal(info)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (p *ProjectInfo) loadFromGit() {
	if name := runCommand("git", "config", "user.name"); name != "" {
		p.Author = strings.TrimSpace(name)
	}

	if email := runCommand("git", "config", "user.email"); email != "" {
		p.Email = strings.TrimSpace(email)
	}

	if remote := runCommand("git", "remote", "get-url", "origin"); remote != "" {
		p.GitHubUser = parseGitHubUser(remote)
		p.RepoName = parseRepoName(remote)
	}
}

func (p *ProjectInfo) promptUser() {
	fmt.Println("Project Information:")
	fmt.Println()

	p.Name = ui.Input("Project name", p.Name)
	p.Description = ui.Input("Description", p.Description)
	p.Author = ui.Input("Author", p.Author)
	p.Email = ui.Input("Email", p.Email)
	p.GitHubUser = ui.Input("GitHub username", p.GitHubUser)
	p.RepoName = ui.Input("Repository name", p.RepoName)
	p.License = ui.Input("License (MIT/Apache-2.0/GPL-3.0/BSD-3-Clause)", p.License)
}

// setDefaults fills blanks from the environment.
//
// Guessed values are intentionally left blank rather than written as
// placeholders: persisting "yourusername" would stop the prompt on the next run
// and leak the placeholder into every generated file.
func (p *ProjectInfo) setDefaults() {
	p.Author = orDefault(p.Author, getDefaultAuthor())
	p.Email = orDefault(p.Email, slugify(p.Author)+"@example.com")
	p.RepoName = orDefault(p.RepoName, p.Name)
	p.Description = orDefault(p.Description, "A "+p.RepoName+" project")
	p.License = orDefault(p.License, "MIT")
}

// Complete reports whether every field a template may reference is known.
// Incomplete projects still generate files; unresolved fields fall back to
// sensible defaults at render time.
func (p *ProjectInfo) Complete() bool {
	return p != nil &&
		p.Name != "" && p.Author != "" && p.Email != "" &&
		p.GitHubUser != "" && p.RepoName != "" && p.License != ""
}
func (p *ProjectInfo) buildDerived() {
	if p.GitHubUser != "" && p.RepoName != "" {
		p.RepoURL = fmt.Sprintf("https://github.com/%s/%s", p.GitHubUser, p.RepoName)
		p.Domain = fmt.Sprintf("%s.github.io/%s", strings.ToLower(p.GitHubUser), strings.ToLower(p.RepoName))
		p.DockerImage = fmt.Sprintf("%s/%s", strings.ToLower(p.GitHubUser), strings.ToLower(p.RepoName))
	} else {
		p.Domain = "example.com"
		p.DockerImage = strings.ToLower(p.Name)
	}
}

// ToVars returns the template variable set for this project.
//
// This is the single source of template variables: a template may only
// reference keys defined here, which TestEveryTemplateRenders enforces.
func (p *ProjectInfo) ToVars() map[string]string {
	if p == nil {
		p = getDefaultProjectInfo()
	}

	name := orDefault(p.Name, "project")
	desc := orDefault(p.Description, "A "+name+" project")
	author := orDefault(p.Author, "Your Name")
	email := orDefault(p.Email, slugify(author)+"@example.com")
	user := orDefault(p.GitHubUser, "yourusername")
	repo := orDefault(p.RepoName, name)
	license := orDefault(p.License, "MIT")

	vars := getCurrentVars()
	vars["project_name"] = name
	vars["project_desc"] = desc
	vars["author"] = author
	vars["fullname"] = author
	vars["email"] = email
	vars["github_username"] = user
	vars["repo_name"] = repo
	vars["repo_url"] = fmt.Sprintf("https://github.com/%s/%s", user, repo)
	vars["license"] = license
	vars["domain"] = fmt.Sprintf("%s.github.io/%s", strings.ToLower(user), strings.ToLower(repo))
	vars["docker_image"] = strings.ToLower(user + "/" + repo)
	vars["module_path"] = fmt.Sprintf("github.com/%s/%s", user, repo)

	return vars
}

// orDefault returns v, or fallback when v is blank.
func orDefault(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

// slugify reduces a name to characters valid in an email local part.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ', r == '_', r == '-', r == '.':
			b.WriteRune('.')
		}
	}
	out := strings.Trim(b.String(), ".")
	if out == "" {
		return "user"
	}
	return out
}

// === Helpers ===

func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func parseGitHubUser(remote string) string {
	remote = strings.TrimSpace(remote)
	remote = strings.TrimPrefix(remote, "git@github.com:")
	remote = strings.TrimPrefix(remote, "https://github.com/")
	remote = strings.TrimPrefix(remote, "http://github.com/")
	remote = strings.TrimSuffix(remote, ".git")

	parts := strings.Split(remote, "/")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

func parseRepoName(remote string) string {
	remote = strings.TrimSpace(remote)
	remote = strings.TrimPrefix(remote, "git@github.com:")
	remote = strings.TrimPrefix(remote, "https://github.com/")
	remote = strings.TrimPrefix(remote, "http://github.com/")
	remote = strings.TrimSuffix(remote, ".git")

	parts := strings.Split(remote, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func getDefaultAuthor() string {
	if name := os.Getenv("USER"); name != "" {
		return name
	}
	if name := os.Getenv("USERNAME"); name != "" {
		return name
	}
	return "Your Name"
}
