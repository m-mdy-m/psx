package resources

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/shared"
)

// ProjectInfo holds information about the project
type ProjectInfo struct {
	Name          string
	Description   string
	Author        string
	Email         string
	GitHubUser    string
	RepoName      string
	RepoURL       string
	License       string
	Domain        string
	DockerImage   string
	SupportEmail  string
	SecurityEmail string
	CurrentDir    string
}

// GetProjectInfo collects project information from various sources and user input
func GetProjectInfo(projectPath string, interactive bool) *ProjectInfo {
	info := &ProjectInfo{
		CurrentDir: projectPath,
	}

	// Try to get info from git
	info.loadFromGit()

	// Try to get info from package files
	info.loadFromPackageFiles(projectPath)

	// Get info interactively if needed
	if interactive {
		info.PromptUser()
	} else {
		// Use defaults if not interactive
		info.SetDefaults()
	}

	// Build derived fields
	info.BuildDerivedFields()

	return info
}

func (p *ProjectInfo) loadFromGit() {
	// Try to get git user name and email
	if name := runCommand("git", "config", "user.name"); name != "" {
		p.Author = strings.TrimSpace(name)
	}

	if email := runCommand("git", "config", "user.email"); email != "" {
		p.Email = strings.TrimSpace(email)
	}

	// Try to get GitHub username from git remote
	if remote := runCommand("git", "remote", "get-url", "origin"); remote != "" {
		if user := parseGitHubUserFromRemote(remote); user != "" {
			p.GitHubUser = user
		}
		if repo := parseRepoNameFromRemote(remote); repo != "" {
			p.RepoName = repo
		}
	}
}

func (p *ProjectInfo) loadFromPackageFiles(projectPath string) {
	// Try to get project name from directory name
	if p.Name == "" {
		p.Name = filepath.Base(projectPath)
	}

	// Try to get info from package.json (Node.js)
	packageJSON := filepath.Join(projectPath, "package.json")
	if exists, _ := shared.FileExists(packageJSON); exists {
		p.loadFromPackageJSON(packageJSON)
	}

	// Try to get info from Cargo.toml (Rust)
	cargoToml := filepath.Join(projectPath, "Cargo.toml")
	if exists, _ := shared.FileExists(cargoToml); exists {
		p.loadFromCargoToml(cargoToml)
	}

	// Try to get info from go.mod (Go)
	goMod := filepath.Join(projectPath, "go.mod")
	if exists, _ := shared.FileExists(goMod); exists {
		p.loadFromGoMod(goMod)
	}

	// Try to get info from setup.py or pyproject.toml (Python)
	setupPy := filepath.Join(projectPath, "setup.py")
	pyprojectToml := filepath.Join(projectPath, "pyproject.toml")
	if exists, _ := shared.FileExists(setupPy); exists {
		p.loadFromSetupPy(setupPy)
	} else if exists, _ := shared.FileExists(pyprojectToml); exists {
		p.loadFromPyprojectToml(pyprojectToml)
	}
}

func (p *ProjectInfo) loadFromPackageJSON(path string) {
	// Simple parsing - in real implementation use json.Unmarshal
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	// Extract name
	if name := extractJSONField(content, "name"); name != "" {
		if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
			p.Name = name
		}
	}

	// Extract description
	if desc := extractJSONField(content, "description"); desc != "" {
		p.Description = desc
	}

	// Extract author
	if author := extractJSONField(content, "author"); author != "" {
		if p.Author == "" {
			p.Author = author
		}
	}

	// Extract license
	if license := extractJSONField(content, "license"); license != "" {
		p.License = license
	}

	// Extract repository
	if repo := extractJSONField(content, "repository"); repo != "" {
		if strings.Contains(repo, "github.com") {
			if user := parseGitHubUserFromURL(repo); user != "" {
				p.GitHubUser = user
			}
			if repoName := parseRepoNameFromURL(repo); repoName != "" {
				p.RepoName = repoName
			}
		}
	}
}

func (p *ProjectInfo) loadFromCargoToml(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	// Extract name from [package] section
	if name := extractTomlField(content, "name"); name != "" {
		if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
			p.Name = name
		}
	}

	// Extract description
	if desc := extractTomlField(content, "description"); desc != "" {
		p.Description = desc
	}

	// Extract license
	if license := extractTomlField(content, "license"); license != "" {
		p.License = license
	}
}

func (p *ProjectInfo) loadFromGoMod(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimPrefix(line, "module ")
			parts := strings.Split(module, "/")
			if len(parts) >= 2 {
				if parts[0] == "github.com" && len(parts) >= 3 {
					p.GitHubUser = parts[1]
					p.RepoName = parts[2]
				}
				// Use last part as project name
				if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
					p.Name = parts[len(parts)-1]
				}
			}
			break
		}
	}
}

func (p *ProjectInfo) loadFromSetupPy(path string) {
	// Simple extraction from setup.py
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	if name := extractPythonStringArg(content, "name"); name != "" {
		if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
			p.Name = name
		}
	}

	if desc := extractPythonStringArg(content, "description"); desc != "" {
		p.Description = desc
	}

	if license := extractPythonStringArg(content, "license"); license != "" {
		p.License = license
	}
}

func (p *ProjectInfo) loadFromPyprojectToml(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	if name := extractTomlField(content, "name"); name != "" {
		if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
			p.Name = name
		}
	}

	if desc := extractTomlField(content, "description"); desc != "" {
		p.Description = desc
	}

	if license := extractTomlField(content, "license"); license != "" {
		p.License = license
	}
}

func (p *ProjectInfo) PromptUser() {
	fmt.Println("Let's set up your project information:")
	fmt.Println()

	// Project name
	if p.Name == "" || p.Name == filepath.Base(p.CurrentDir) {
		p.Name = shared.PromptInput("Project name", filepath.Base(p.CurrentDir))
	} else {
		p.Name = shared.PromptInput("Project name", p.Name)
	}

	// Description
	p.Description = shared.PromptInput("Project description", p.Description)

	// Author
	if p.Author == "" {
		p.Author = getUserName()
	}
	p.Author = shared.PromptInput("Author name", p.Author)

	// Email
	p.Email = shared.PromptInput("Email", p.Email)

	// GitHub username
	p.GitHubUser = shared.PromptInput("GitHub username", p.GitHubUser)

	// Repo name
	if p.RepoName == "" {
		p.RepoName = p.Name
	}
	p.RepoName = shared.PromptInput("Repository name", p.RepoName)

	// License
	if p.License == "" {
		p.License = "MIT"
	}
	p.License = shared.PromptInput("License", p.License)
}

func (p *ProjectInfo) SetDefaults() {
	// Set defaults for missing fields
	if p.Name == "" {
		p.Name = filepath.Base(p.CurrentDir)
	}

	if p.Author == "" {
		p.Author = getUserName()
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

	if p.License == "" {
		p.License = "MIT"
	}

	if p.Description == "" {
		p.Description = fmt.Sprintf("A %s project", p.Name)
	}
}

func (p *ProjectInfo) BuildDerivedFields() {
	// Build repo URL
	if p.GitHubUser != "" && p.RepoName != "" {
		p.RepoURL = fmt.Sprintf("https://github.com/%s/%s", p.GitHubUser, p.RepoName)
	}

	// Build Docker image name
	if p.GitHubUser != "" && p.RepoName != "" {
		p.DockerImage = fmt.Sprintf("%s/%s", strings.ToLower(p.GitHubUser), strings.ToLower(p.RepoName))
	} else {
		p.DockerImage = strings.ToLower(p.Name)
	}

	// Build support email
	if p.Email != "" {
		p.SupportEmail = p.Email
		p.SecurityEmail = p.Email
	} else {
		p.SupportEmail = "support@example.com"
		p.SecurityEmail = "security@example.com"
	}

	// Build domain (if not set)
	if p.Domain == "" {
		if p.GitHubUser != "" && p.RepoName != "" {
			p.Domain = fmt.Sprintf("%s.github.io/%s", strings.ToLower(p.GitHubUser), strings.ToLower(p.RepoName))
		} else {
			p.Domain = "example.com"
		}
	}
}

// ToVars converts ProjectInfo to template variables
func (p *ProjectInfo) ToVars() map[string]string {
	// Handle nil pointer
	if p == nil {
		p = &ProjectInfo{
			Name:        "my-project",
			Description: "A project",
			Author:      "Your Name",
			Email:       "email@example.com",
			GitHubUser:  "yourusername",
			RepoName:    "my-project",
			License:     "MIT",
		}
		p.BuildDerivedFields()
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

// Helper functions

func parseGitHubUserFromRemote(remote string) string {
	// git@github.com:user/repo.git -> user
	// https://github.com/user/repo.git -> user
	return parseGitHubUserFromURL(remote)
}

func parseRepoNameFromRemote(remote string) string {
	return parseRepoNameFromURL(remote)
}

func parseGitHubUserFromURL(url string) string {
	url = strings.TrimSpace(url)

	// Remove common prefixes
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "http://github.com/")

	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) >= 1 {
		return parts[0]
	}

	return ""
}

// runCommand runs a command and returns its output
func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func parseRepoNameFromURL(url string) string {
	url = strings.TrimSpace(url)

	// Remove common prefixes
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "http://github.com/")

	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

func extractJSONField(content, field string) string {
	// Simple extraction - in production use json.Unmarshal
	pattern := fmt.Sprintf(`"%s"\s*:\s*"([^"]*)"`, field)
	return extractPattern(content, pattern)
}

func extractTomlField(content, field string) string {
	pattern := fmt.Sprintf(`%s\s*=\s*"([^"]*)"`, field)
	return extractPattern(content, pattern)
}

func extractPythonStringArg(content, field string) string {
	pattern := fmt.Sprintf(`%s\s*=\s*["']([^"']*)["']`, field)
	return extractPattern(content, pattern)
}

func extractPattern(content, pattern string) string {
	// Simple regex-like extraction
	// In production, use regexp package
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, pattern) {
			// Simple extraction
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, `"'`)
				value = strings.TrimSuffix(value, ",")
				return value
			}
		}
	}
	return ""
}
