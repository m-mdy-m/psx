package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/utils"
)

//go:embed embedded/*.yml
var configFS embed.FS

var (
	rulesMetadata *RulesMetadata
	defaultConfig *Config
)

func init() {
	var err error

	rulesMetadata, err = utils.LoadEmbedded[RulesMetadata]("rules metadata", "embedded/rules.yml", configFS)
	if err != nil {
		logger.Fatalf("Failed to load rules metadata: %v", err)
	}
	logger.Verbose(fmt.Sprintf("Loaded %d rules from metadata", len(rulesMetadata.Rules)))
	defaultConfig, err = utils.LoadEmbedded[Config]("default config", "embedded/psx.default.yml", configFS)
	if err != nil {
		logger.Fatalf("Failed to load default config: %v", err)
	}
	logger.Verbose("Default configuration loaded")
}
func GetRulesMetadata() *RulesMetadata {
	return rulesMetadata
}

// Load resolves a configuration for projectPath.
//
// configFile may be empty, in which case the standard locations are searched and
// the embedded defaults are used when nothing is found.
func Load(configFile string, projectPath string) (*Config, error) {
	var userConfig *Config
	var err error

	if configFile == "" {
		logger.Verbose("Searching for config file...")
		configFile, err = FindConfigFile(projectPath)
		if err != nil || configFile == "" {
			logger.Verbose("No config file found, using embedded defaults")
			return buildConfig(defaultConfig, projectPath, "")
		}
		logger.Verbose(fmt.Sprintf("Found config file: %s", configFile))
	}

	userConfig, err = readConfigFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configFile, err)
	}
	logger.Verbose(fmt.Sprintf("Parsed user config: %s", configFile))

	if result := Validate(userConfig); !IsValid(result) {
		logger.Error("Configuration validation failed:")
		for _, e := range result.Errors {
			logger.Error(fmt.Sprintf("  [%s] %s", e.Field, e.Message))
		}
		return nil, fmt.Errorf("config validation failed: %d errors", len(result.Errors))
	} else {
		for _, w := range result.Warnings {
			logger.Warning(w)
		}
	}

	projectType := resources.NormalizeProjectType(userConfig.Project.Type)
	cfg, err := buildConfig(userConfig, projectPath, projectType)
	if err != nil {
		return nil, err
	}
	cfg.ConfigFile = configFile
	return cfg, nil
}

func FindConfigFile(projectPath string) (string, error) {
	candidates := []string{
		"psx.yml",
		".psx.yml",
		"psx.yaml",
		".psx.yaml",
	}

	logger.Verbose(fmt.Sprintf("Looking for config in: %s", projectPath))

	// Check current directory
	for _, name := range candidates {
		path := filepath.Join(projectPath, name)
		if exists, info := utils.FileExists(path); exists && !info.IsDir() {
			logger.Verbose(fmt.Sprintf("Found config: %s", path))
			return path, nil
		}
	}

	// Check parent directories up to git root
	current := projectPath
	maxDepth := 10
	depth := 0

	for depth < maxDepth {
		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		// Check if git root
		gitPath := filepath.Join(current, ".git")
		if exists, info := utils.FileExists(gitPath); exists && info.IsDir() {
			logger.Verbose(fmt.Sprintf("Found git root: %s", current))
			for _, name := range candidates {
				path := filepath.Join(current, name)
				if exists, info := utils.FileExists(path); exists && !info.IsDir() {
					logger.Verbose(fmt.Sprintf("Found config in git root: %s", path))
					return path, nil
				}
			}
			break
		}

		current = parent
		depth++
	}

	// Check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(home, ".config", "psx")
		for _, name := range candidates {
			path := filepath.Join(configDir, name)
			if exists, info := utils.FileExists(path); exists && !info.IsDir() {
				logger.Verbose(fmt.Sprintf("Found config in home: %s", path))
				return path, nil
			}
		}
	}

	logger.Verbose("No config file found")
	return "", fmt.Errorf("no config file found")
}

// readConfigFile reads and parses a config file
func readConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Successfully parsed YAML from %s", path))
	return &cfg, nil
}

// buildConfig resolves user input into the active rule set.
// The normalized project type is written back onto the result so downstream
// consumers cannot accidentally read the raw, un-normalized value.
func buildConfig(userCfg *Config, projectPath string, projectType string) (*Config, error) {
	cfg := &Config{
		Version:     userCfg.Version,
		Project:     userCfg.Project,
		Rules:       userCfg.Rules,
		Ignore:      userCfg.Ignore,
		Fix:         userCfg.Fix,
		ProjectPath: projectPath,
		Custom:      userCfg.Custom,
		ActiveRules: make(map[string]*ActiveRule),
	}
	if projectType == "" {
		projectType = "generic"
	}
	cfg.Project.Type = projectType
	cfg.ProjectType = projectType
	// Keep the raw value: an empty project.type must stay distinguishable from
	// an explicit "generic", or detection can never run.
	cfg.DeclaredType = strings.TrimSpace(userCfg.Project.Type)

	enabledCount := 0
	disabledCount := 0
	if len(userCfg.Rules) == 0 {
		logger.Verbose("No rules declared, enabling all with default severity")
		for id, meta := range rulesMetadata.Rules {
			cfg.ActiveRules[id] = &ActiveRule{
				ID:       id,
				Metadata: meta,
				Severity: meta.DefaultSeverity,
			}
			enabledCount++
		}
	} else {
		logger.Verbose("Using declared rules as an explicit allow-list")
		for id, userSev := range userCfg.Rules {
			meta, exists := rulesMetadata.Rules[id]
			if !exists {
				logger.Warning(fmt.Sprintf("Unknown rule '%s' - skipping", id))
				continue
			}

			severity, err := ParseSeverity(userSev, meta.DefaultSeverity)
			if err != nil {
				logger.Warning(fmt.Sprintf("Rule %s: %v, skipping", id, err))
				continue
			}
			if severity == nil {
				logger.Verbose(fmt.Sprintf("Rule %s is disabled", id))
				disabledCount++
				continue
			}

			cfg.ActiveRules[id] = &ActiveRule{
				ID:       id,
				Metadata: meta,
				Severity: *severity,
			}
			enabledCount++
		}
	}

	logger.Verbose(fmt.Sprintf("Config built: %d enabled, %d disabled rules", enabledCount, disabledCount))
	return cfg, nil
}

// ResolvePatterns returns the patterns that apply to projectType.
//
// Lookup order is "<type>", then "*", then "generic". The generic tier matters:
// without it a rule that only declares a generic pattern would resolve to
// nothing for an unrecognised project type.
//
// A rule that declares only language-specific tiers resolves to nothing for an
// unrelated type, which the engine reports as skipped. Falling back to an
// arbitrary language's patterns would apply Go-only expectations to a Python
// project and produce misleading failures.
func ResolvePatterns(patterns any, projectType string) []string {
	table, ok := patterns.(map[string]any)
	if !ok {
		return toStringList(patterns)
	}

	for _, key := range []string{projectType, "*", "generic"} {
		if key == "" {
			continue
		}
		if list := toStringList(table[key]); len(list) > 0 {
			return list
		}
	}
	return nil
}

// toStringList coerces a YAML sequence into a string slice.
func toStringList(v any) []string {
	items, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		if s, ok := it.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}
