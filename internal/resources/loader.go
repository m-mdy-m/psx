package resources

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

//go:embed embedded/*.yml
var embeddedFS embed.FS

var (
	templates     *TemplatesConfig
	gitignores    *GitignoresConfig
	licenses      *LicensesConfig
	qualityTools  *QualityToolsConfig
	devops        *DevOpsConfig
	docsTemplates *DocsTemplatesConfig
	messages      *MessagesConfig
	languages     *LanguagesConfig
)

func init() {
	var err error

	templates, err = utils.LoadEmbedded[TemplatesConfig]("templates", "embedded/templates.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load templates: %v", err)
	}

	gitignores, err = utils.LoadEmbedded[GitignoresConfig]("gitignores", "embedded/gitignores.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load gitignores: %v", err)
	}

	licenses, err = utils.LoadEmbedded[LicensesConfig]("licenses", "embedded/licenses.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load licenses: %v", err)
	}

	qualityTools, err = utils.LoadEmbedded[QualityToolsConfig]("quality-tools", "embedded/quality-tools.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load quality-tools: %v", err)
	}

	devops, err = utils.LoadEmbedded[DevOpsConfig]("devops", "embedded/devops.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load devops: %v", err)
	}

	docsTemplates, err = utils.LoadEmbedded[DocsTemplatesConfig]("docs-templates", "embedded/docs-templates.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load docs-templates: %v", err)
	}

	messages, err = utils.LoadEmbedded[MessagesConfig]("messages", "embedded/messages.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load messages: %v", err)
	}

	languages, err = utils.LoadEmbedded[LanguagesConfig]("languages", "embedded/languages.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load languages: %v", err)
	}

	logger.Verbose("All resources loaded successfully")
}

// replaceVars replaces {{variable}} placeholders in templates
func replaceVars(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// getCurrentVars returns current date/year variables
func getCurrentVars() map[string]string {
	now := time.Now()
	return map[string]string{
		"year": fmt.Sprintf("%d", now.Year()),
		"date": now.Format("2006-01-02"),
	}
}

// normalizeProjectType normalizes project type names
func NormalizeProjectType(projectType string) string {
	if projectType == "" {
		return "generic"
	}

	// Check aliases
	if canonical, ok := languages.Aliases[projectType]; ok {
		return canonical
	}

	return projectType
}

// === Template Getters ===

func GetReadme(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()
	template := getTemplate(templates.Readme, projectType)
	return replaceVars(template, vars)
}

func GetChangelog(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(templates.Changelog, vars)
}

func GetContributing() string {
	return templates.Contributing
}

func GetGitignore(projectType string) string {
	common := gitignores.Common
	specific := getTemplate(map[string]string{
		"nodejs": gitignores.NodeJS,
		"go":     gitignores.Go,
	}, projectType)

	if specific != "" {
		return common + "\n\n" + specific
	}
	return common
}

func GetLicense(licenseType, author string) string {
	license, ok := (*licenses)[licenseType]
	if !ok {
		license = (*licenses)["MIT"]
	}

	vars := getCurrentVars()
	vars["fullname"] = author
	if author == "" {
		vars["fullname"] = "Your Name"
	}

	return replaceVars(license.Content, vars)
}

func GetEditorconfig(projectType string) string {
	return getTemplate(qualityTools.Editorconfig, projectType)
}

func GetDockerfile(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	var template string
	switch projectType {
	case "nodejs":
		template = devops.Docker.NodeJS.Dockerfile
	case "go":
		template = devops.Docker.Go.Dockerfile
	default:
		return ""
	}

	return replaceVars(template, vars)
}

func GetDockerignore(projectType string) string {
	switch projectType {
	case "nodejs":
		return devops.Docker.NodeJS.Dockerignore
	case "go":
		return devops.Docker.Go.Dockerignore
	default:
		return ""
	}
}

func GetSecurity(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Security, vars)
}

func GetCodeOfConduct(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.CodeOfConduct, vars)
}

func GetMessage(category, key string) string {
	switch category {
	case "check":
		if msg, ok := messages.Check[key]; ok {
			return msg
		}
	case "fix":
		if msg, ok := messages.Fix[key]; ok {
			return msg
		}
	case "errors":
		if msg, ok := messages.Errors[key]; ok {
			return msg
		}
	case "help":
		if msg, ok := messages.Help[key]; ok {
			return msg
		}
	}
	return ""
}

func FormatMessage(category, key string, args ...any) string {
	msg := GetMessage(category, key)
	if msg == "" {
		return ""
	}
	return fmt.Sprintf(msg, args...)
}

// === Helpers ===

// getTemplate gets template for project type, falls back to generic
func getTemplate(templates map[string]string, projectType string) string {
	if t, ok := templates[projectType]; ok {
		return t
	}
	if t, ok := templates["generic"]; ok {
		return t
	}
	// Return first available
	for _, t := range templates {
		return t
	}
	return ""
}
