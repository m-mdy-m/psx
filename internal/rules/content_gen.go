package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/resources"
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
	// General files (from templates.yml)
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
	case "api_docs":
		return resources.GetAPIDocs(cg.projectInfo, cg.projectType)
	// Documentation (from docs-templates.yml)
	case "security":
		return resources.GetSecurity(cg.projectInfo)
	case "code_of_conduct":
		return resources.GetCodeOfConduct(cg.projectInfo)
	case "pull_request_template":
		return resources.GetPullRequestTemplate(cg.projectInfo)
	case "codeowners":
		return resources.GetCodeowners(cg.projectInfo)

	// Quality tools (from quality-tools.yml)
	case "editorconfig":
		return resources.GetEditorconfig(cg.projectType)
	case "pre_commit":
		return resources.GetPreCommit(cg.projectType)

	// DevOps (from devops.yml)
	case "dockerfile":
		return resources.GetDockerfile(cg.projectInfo, cg.projectType)
	case "dockerignore":
		return resources.GetDockerignore(cg.projectType)
	case "docker_compose":
		return resources.GetDockerCompose(cg.projectInfo, cg.projectType)
	default:
		return ""
	}
}

func (cg *ContentGenerator) generateByPattern(filename string) string {
	lowerFilename := strings.ToLower(filename)

	// Check for specific filenames
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
		return resources.GetDockerCompose(cg.projectInfo, cg.projectType)
	case lowerFilename == "codeowners":
		return resources.GetCodeowners(cg.projectInfo)
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
	}

	return nil, nil
}
