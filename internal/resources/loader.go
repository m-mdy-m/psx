package resources

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

//go:embed embedded/*.yml
var resourcesFS embed.FS

var (
	languages  *LanguagesConfig
	messages   *MessagesConfig
	gitignores *GitignoresConfig
	licenses   *LicensesConfig
	templates  *TemplatesConfig
)

func init() {
	var err error

	languages, err = loadYAML[LanguagesConfig]("embedded/languages.yml")
	if err != nil {
		logger.Fatalf("Failed to load languages: %v", err)
	}

	messages, err = loadYAML[MessagesConfig]("embedded/messages.yml")
	if err != nil {
		logger.Fatalf("Failed to load messages: %v", err)
	}

	gitignores, err = loadYAML[GitignoresConfig]("embedded/gitignores.yml")
	if err != nil {
		logger.Fatalf("Failed to load gitignores: %v", err)
	}

	licenses, err = loadYAML[LicensesConfig]("embedded/licenses.yml")
	if err != nil {
		logger.Fatalf("Failed to load licenses: %v", err)
	}

	templates, err = loadYAML[TemplatesConfig]("embedded/templates.yml")
	if err != nil {
		logger.Fatalf("Failed to load templates: %v", err)
	}

	logger.Verbose("Resources loaded successfully")
}

func loadYAML[T any](filename string) (*T, error) {
	data, err := resourcesFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}

	var result T
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}

	return &result, nil
}

// Template variable replacement
func replaceVars(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// Get current year and date
func getCurrentVars() map[string]string {
	now := time.Now()
	return map[string]string{
		"year": fmt.Sprintf("%d", now.Year()),
		"date": now.Format("2006-01-02"),
	}
}

// Get user name from git config or environment
func getUserName() string {
	if name := os.Getenv("USER"); name != "" {
		return name
	}
	if name := os.Getenv("USERNAME"); name != "" {
		return name
	}
	return "Your Name"
}