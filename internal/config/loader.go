package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

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
func Load(configFile string, projectPath string) (*Config, error) {
	var userConfig *Config
	var err error
	if configFile == "" {
		logger.Verbose("Searching for config file...")
		configFile, err = FindConfigFile(projectPath)
		if err != nil || configFile == "" {
			logger.Info("No config file found, using defaults")
			return buildConfig(defaultConfig, projectPath, "")
		}
		logger.Verbose(fmt.Sprintf("Found config file: %s", configFile))
	}
	userConfig, err = readConfigFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configFile, err)
	}

	logger.Verbose(fmt.Sprintf("Loaded user config from: %s", configFile))

	result := Validate(userConfig)
	if HasWarnings(result) {
		logger.Warning("Configuration warnings:")
		for _, warning := range result.Warnings {
			logger.Warning(warning)
		}
	}
	if !IsValid(result) {
		logger.Error("Configuration validation failed:")
		for _, err := range result.Errors {
			logger.Error(fmt.Sprintf("  [%s] %s", err.Field, err.Message))
		}
		return nil, fmt.Errorf("config validation failed: %d errors", len(result.Errors))
	}

	logger.Success("Configuration loaded and validated")
	projectType := resources.NormalizeProjectType(userConfig.Project.Type)

	return buildConfig(userConfig, projectPath, projectType)
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

// buildConfig builds a complete config with active rules
func buildConfig(userCfg *Config, projectPath string, projectType string) (*Config, error) {
	cfg := &Config{
		Version:     userCfg.Version,
		Project:     userCfg.Project,
		Rules:       userCfg.Rules,
		Ignore:      userCfg.Ignore,
		Fix:         userCfg.Fix,
		Path:        projectPath,
		Custom:      userCfg.Custom,
		ActiveRules: make(map[string]*ActiveRule),
	}
	enabledCount := 0
	disabledCount := 0
	if len(userCfg.Rules) == 0 {
		logger.Verbose("Using default config - enabling all rules")
		for id, meta := range rulesMetadata.Rules {
			cfg.ActiveRules[id] = &ActiveRule{
				ID:       id,
				Metadata: meta,
				Severity: meta.DefaultSeverity,
			}
			enabledCount++
			logger.Verbose(fmt.Sprintf("Rule %s enabled (default) with severity: %s", id, meta.DefaultSeverity))
		}
	} else {
		logger.Verbose("Using user config - enabling only specified rules")

		for id, userSev := range userCfg.Rules {
			// Get metadata
			meta, exists := rulesMetadata.Rules[id]
			if !exists {
				logger.Warning(fmt.Sprintf("Unknown rule '%s' - skipping", id))
				continue
			}

			// Parse severity
			severity, err := ParseSeverity(userSev, meta.DefaultSeverity)
			if err != nil {
				logger.Warning(fmt.Sprintf("Rule %s: %v, skipping", id, err))
				continue
			}

			// If disabled
			if severity == nil {
				logger.Verbose(fmt.Sprintf("Rule %s is disabled", id))
				disabledCount++
				continue
			}

			// Enable rule
			cfg.ActiveRules[id] = &ActiveRule{
				ID:       id,
				Metadata: meta,
				Severity: *severity,
			}
			enabledCount++
			logger.Verbose(fmt.Sprintf("Rule %s enabled with severity: %s", id, *severity))
		}
	}

	logger.Verbose(fmt.Sprintf("Config built: %d enabled, %d disabled rules", enabledCount, disabledCount))
	return cfg, nil
}

func GetPatterns(patterns any, projectType string) []string {
	switch p := patterns.(type) {
	case []any:
		// Simple list of patterns
		result := make([]string, 0, len(p))
		for _, item := range p {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result

	case map[string]any:
		if projectType != "" {
			if langPatterns, exists := p[projectType]; exists {
				if arr, ok := langPatterns.([]any); ok {
					result := make([]string, 0, len(arr))
					for _, item := range arr {
						if s, ok := item.(string); ok {
							result = append(result, s)
						}
					}
					return result
				}
			}
		}

		if genericPatterns, exists := p["*"]; exists {
			if arr, ok := genericPatterns.([]any); ok {
				result := make([]string, 0, len(arr))
				for _, item := range arr {
					if s, ok := item.(string); ok {
						result = append(result, s)
					}
				}
				return result
			}
		}
	}

	return []string{}
}
