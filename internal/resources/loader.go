// Add this to internal/resources/loader.go init() function

package resources

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

//go:embed embedded/*.yml
var embeddedFS embed.FS

var (
	templates      *TemplatesConfig
	gitignores     *GitignoresConfig
	licenses       *LicensesConfig
	qualityTools   *QualityToolsConfig
	devops         *DevOpsConfig
	docsTemplates  *DocsTemplatesConfig
	messages       *MessagesConfig
	languages      *LanguagesConfig
	scripts        *ScriptsConfig
	projectScripts *ProjectScriptsConfig
	actions        *GitHubActionsConfig
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

	scripts, err = utils.LoadEmbedded[ScriptsConfig]("scripts", "embedded/scripts.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load scripts: %v", err)
	}

	actions, err = utils.LoadEmbedded[GitHubActionsConfig]("github-actions", "embedded/github-actions.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load github-actions: %v", err)
	}

	projectScripts, err = utils.LoadEmbedded[ProjectScriptsConfig]("project-scripts", "embedded/project-scripts.yml", embeddedFS)
	if err != nil {
		logger.Fatalf("Failed to load project scripts: %v", err)
	}

	logger.Verbose("All resources loaded successfully")
}

func replaceVars(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
func getCurrentVars() map[string]string {
	now := time.Now()
	return map[string]string{
		"year": fmt.Sprintf("%d", now.Year()),
		"date": now.Format("2006-01-02"),
	}
}
func NormalizeProjectType(projectType string) string {
	if projectType == "" {
		return "generic"
	}

	if canonical, ok := languages.Aliases[projectType]; ok {
		return canonical
	}

	return projectType
}

// getTemplate picks the project-type tier, falling back to generic.
//
// The iteration order over a map is random, so keys are sorted before the last
// resort: template selection must be reproducible across runs.
func getTemplate(table map[string]string, projectType string) string {
	if t, ok := table[projectType]; ok && t != "" {
		return t
	}
	if t, ok := table["generic"]; ok && t != "" {
		return t
	}
	keys := make([]string, 0, len(table))
	for k := range table {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if t := table[k]; t != "" {
			return t
		}
	}
	return ""
}

func getDefaultProjectInfo() *ProjectInfo {
	info := &ProjectInfo{
		Name:        "project",
		Description: "A new project",
		Author:      "Your Name",
		Email:       "you@example.com",
		GitHubUser:  "yourusername",
		RepoName:    "project",
		License:     "MIT",
	}
	info.buildDerived()
	return info
}
