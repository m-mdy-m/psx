package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/utils"
)

type ContentGenerator struct {
	projectInfo *resources.ProjectInfo
	projectType string
}

func NewContentGenerator(projectInfo *resources.ProjectInfo, projectType string) *ContentGenerator {
	return &ContentGenerator{
		projectInfo: projectInfo,
		projectType: projectType,
	}
}

func (cg *ContentGenerator) Generate(ruleID, pattern string) (string, error) {
	filename := filepath.Base(pattern)

	if content := cg.generateByRuleID(ruleID); content != "" {
		return content, nil
	}

	if content := cg.generateByPattern(filename); content != "" {
		return content, nil
	}

	return fmt.Sprintf("# %s\n\nTODO: Add content\n", filename), nil
}

func (cg *ContentGenerator) generateByRuleID(ruleID string) string {
	switch ruleID {
	// General files
	case "readme":
		return resources.GetReadme(cg.projectInfo, cg.projectType)
	case "license":
		return resources.GetLicense(cg.projectInfo.License, cg.projectInfo.Author)
	case "gitignore":
		return resources.GetGitignore(cg.projectType)
	case "changelog":
		return resources.GetChangelog(cg.projectInfo)
	case "contributing":
		return resources.GetContributing()

	// Documentation
	case "api_docs":
		return resources.GetAPIDocs(cg.projectInfo, cg.projectType)
	case "security":
		return resources.GetSecurity(cg.projectInfo)
	case "code_of_conduct":
		return resources.GetCodeOfConduct(cg.projectInfo)
	case "pull_request_template":
		return resources.GetPullRequestTemplate(cg.projectInfo)
	case "codeowners":
		return resources.GetCodeowners(cg.projectInfo)

	// Quality tools
	case "editorconfig":
		return resources.GetEditorconfig(cg.projectType)
	case "pre_commit":
		return resources.GetPreCommit(cg.projectType)

	// DevOps
	case "dockerfile":
		return resources.GetDockerfile(cg.projectInfo, cg.projectType)
	case "dockerignore":
		return resources.GetDockerignore(cg.projectType)
	case "docker_compose":
		return resources.GetDockerComposeWithPrompt(cg.projectInfo, cg.projectType)

	// CI/CD
	case "ci_config":
		return resources.GetCIConfig(cg.projectInfo, cg.projectType)

	default:
		return ""
	}
}

func (cg *ContentGenerator) generateByPattern(filename string) string {
	lowerFilename := strings.ToLower(filename)

	switch {
	case strings.Contains(lowerFilename, "readme"):
		return resources.GetReadme(cg.projectInfo, cg.projectType)
	case strings.Contains(lowerFilename, "license"):
		return resources.GetLicense(cg.projectInfo.License, cg.projectInfo.Author)
	case lowerFilename == ".gitignore":
		return resources.GetGitignore(cg.projectType)
	case strings.Contains(lowerFilename, "changelog"):
		return resources.GetChangelog(cg.projectInfo)
	case strings.Contains(lowerFilename, "contributing"):
		return resources.GetContributing()
	case strings.Contains(lowerFilename, "security"):
		return resources.GetSecurity(cg.projectInfo)
	case strings.Contains(lowerFilename, "code_of_conduct"):
		return resources.GetCodeOfConduct(cg.projectInfo)
	case lowerFilename == ".editorconfig":
		return resources.GetEditorconfig(cg.projectType)
	case strings.Contains(lowerFilename, "dockerfile"):
		return resources.GetDockerfile(cg.projectInfo, cg.projectType)
	case lowerFilename == ".dockerignore":
		return resources.GetDockerignore(cg.projectType)
	case strings.Contains(lowerFilename, "docker-compose"):
		return resources.GetDockerComposeWithPrompt(cg.projectInfo, cg.projectType)
	case lowerFilename == "codeowners":
		return resources.GetCodeowners(cg.projectInfo)
	case strings.Contains(lowerFilename, "api"):
		return resources.GetAPIDocs(cg.projectInfo, cg.projectType)
	}

	return ""
}

func (cg *ContentGenerator) GenerateMultiple(ruleID string) (map[string]string, error) {
	result := make(map[string]string)

	switch ruleID {
	case "issue_templates":
		result[".github/ISSUE_TEMPLATE/bug_report.yml"] = resources.GetIssueBugReport()
		result[".github/ISSUE_TEMPLATE/feature_request.yml"] = resources.GetIssueFeatureRequest()
		result[".github/ISSUE_TEMPLATE/question.yml"] = resources.GetIssueQuestion()
		result[".github/ISSUE_TEMPLATE/config.yml"] = resources.GetIssueTemplatesConfig(cg.projectInfo)
		return result, nil

	case "adr":
		result["docs/adr/0001-record-architecture-decisions.md"] = resources.GetADRFirst(cg.projectInfo)
		result["docs/adr/template.md"] = resources.GetADRTemplate(cg.projectInfo)
		return result, nil

	case "scripts_folder":
		scripts := resources.GetScripts(cg.projectInfo, cg.projectType)
		for name, content := range scripts {
			result[fmt.Sprintf("scripts/%s", name)] = content
		}
		return result, nil

	case "api_docs":
		apiDocPath := cg.getAPIDocPath()
		content := resources.GetAPIDocs(cg.projectInfo, cg.projectType)
		result[apiDocPath] = content
		return result, nil
	case "ci_config":
		return cg.generateCIConfig()
	}

	return nil, nil
}

func (cg *ContentGenerator) getAPIDocPath() string {
	switch cg.projectType {
	case "nodejs", "go":
		return "docs/api/README.md"
	default:
		return "docs/API.md"
	}
}

func (cg *ContentGenerator) generateCIConfig() (map[string]string, error) {
	result := make(map[string]string)

	platform, err := utils.PromptChoice(
		"Which CI/CD platform do you want to use?",
		[]string{"GitHub Actions", "GitLab CI", "Both", "Skip"},
	)
	if err != nil || platform == "Skip" {
		return nil, nil
	}

	switch platform {
	case "GitHub Actions":
		workflow := resources.GetGitHubActionsWorkflow(cg.projectInfo, cg.projectType)
		result[".github/workflows/ci.yml"] = workflow
	case "GitLab CI":
		config := resources.GetGitLabCIConfig(cg.projectInfo, cg.projectType)
		result[".gitlab-ci.yml"] = config
	case "Both":
		workflow := resources.GetGitHubActionsWorkflow(cg.projectInfo, cg.projectType)
		result[".github/workflows/ci.yml"] = workflow
		config := resources.GetGitLabCIConfig(cg.projectInfo, cg.projectType)
		result[".gitlab-ci.yml"] = config
	}

	return result, nil
}
