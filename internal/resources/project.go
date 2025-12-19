package resources

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

const projectCacheFile = ".psx-project.yml"

// GetProjectInfo loads or creates project info
func GetProjectInfo(projectPath string, interactive bool) *ProjectInfo {
	// Try to load from cache
	if info, err := loadProjectInfo(projectPath); err == nil && info != nil {
		logger.Verbose("Using cached project info")
		return info
	}

	logger.Verbose("Creating new project info")

	info := &ProjectInfo{
		Name:    filepath.Base(projectPath),
		License: "MIT",
	}

	// Try to get info from git
	info.loadFromGit()

	// Ask user if interactive
	if interactive {
		info.promptUser()
	} else {
		info.setDefaults()
	}

	info.buildDerived()

	// Save to cache
	if err := saveProjectInfo(projectPath, info); err != nil {
		logger.Warning(fmt.Sprintf("Failed to save project info: %v", err))
	}

	return info
}

// loadProjectInfo loads project info from cache
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

// saveProjectInfo saves project info to cache
func saveProjectInfo(projectPath string, info *ProjectInfo) error {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := yaml.Marshal(info)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// loadFromGit tries to load info from git config
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

// promptUser asks user for missing information
func (p *ProjectInfo) promptUser() {
	fmt.Println("Project Information:")
	fmt.Println()

	p.Name = utils.PromptInput("Project name", p.Name)
	p.Description = utils.PromptInput("Description", p.Description)
	p.Author = utils.PromptInput("Author", p.Author)
	p.Email = utils.PromptInput("Email", p.Email)
	p.GitHubUser = utils.PromptInput("GitHub username", p.GitHubUser)
	p.RepoName = utils.PromptInput("Repository name", p.RepoName)
	p.License = utils.PromptInput("License", p.License)
}

func (p *ProjectInfo) setDefaults() {
	if p.Author == "" {
		p.Author = getDefaultAuthor()
	}
	if p.Email == "" {
		p.Email = fmt.Sprintf("%s@example.com", strings.ToLower(p.Author))
	}
	if p.GitHubUser == "" {
		p.GitHubUser = "yourusername"
	}
	if p.RepoName == "" {
		p.RepoName = p.Name
	}
	if p.Description == "" {
		p.Description = fmt.Sprintf("A %s project", p.Name)
	}
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

	p.SupportEmail = p.Email
	p.SecurityEmail = p.Email
	if p.Email == "" {
		p.SupportEmail = "support@example.com"
		p.SecurityEmail = "security@example.com"
	}
}

func (p *ProjectInfo) ToVars() map[string]string {
	if p == nil {
		p = &ProjectInfo{
			Name:   "project",
			Author: "Your Name",
			Email:  "you@example.com",
		}
		p.buildDerived()
	}

	vars := getCurrentVars()
	vars["project_name"] = p.Name
	vars["project_desc"] = p.Description
	vars["author"] = p.Author
	vars["fullname"] = p.Author
	vars["email"] = p.Email
	vars["github_username"] = p.GitHubUser
	vars["repo_name"] = p.RepoName
	vars["repo_url"] = p.RepoURL
	vars["license"] = p.License
	vars["domain"] = p.Domain
	vars["docker_image"] = p.DockerImage
	vars["support_email"] = p.SupportEmail
	vars["security_email"] = p.SecurityEmail
	vars["conduct_email"] = p.SupportEmail

	return vars
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
