package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

//go:embed embedded/*.yml
var configFS embed.FS

var (
	globalMetadata *RulesMetadata
	defaultConfig  *Config
)

func GetRulesMetadata() *RulesMetadata {
	return globalMetadata
}

func init() {
	var err error
	globalMetadata, err = LoadRulesMetadata()
	if err != nil {
		logger.Errorf("Failed to load rules metadata: %v", err)
		os.Exit(1)
	}
	logger.Verbose(fmt.Sprintf("Loaded %d rules from metadata", len(globalMetadata.Rules)))

	defaultConfig, err = LoadDefaultConfig()
	if err != nil {
		logger.Errorf("Failed to load default config: %v", err)
		os.Exit(1)
	}
	logger.Verbose("Default configuration loaded")
}

func loadEmbedded[T any](what string, parts string) (*T, error) {
	data, err := configFS.ReadFile(parts)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s (%s): %w", parts, what, err)
	}

	var v T
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to parse %s (%s): %w", parts, what, err)
	}

	logger.Verbose(fmt.Sprintf("Loaded %s from %s", what, parts))
	return &v, nil
}

func LoadRulesMetadata() (*RulesMetadata, error) {
	md, err := loadEmbedded[RulesMetadata]("rules metadata", "embedded/rules.yml")
	if err != nil {
		return nil, err
	}
	logger.Verbose("Loaded rules metadata: " + fmt.Sprint(len(md.Rules)) + " rules")
	return md, nil
}

func LoadDefaultConfig() (*Config, error) {
	cfg, err := loadEmbedded[Config]("default config", "embedded/psx.default.yml")
	if err != nil {
		return nil, err
	}
	logger.Verbose("Loaded default configuration")
	return cfg, nil
}

func Load(cf string, pp string) (*Config, error) {
	var (
		err        error
		userConfig *Config
	)

	// If no config file specified, search for it
	if cf == "" {
		logger.Verbose("Searching for config file...")
		found, err := FindConfigFile(pp)
		if err != nil {
			logger.Verbose(fmt.Sprintf("No config file found: %v", err))
			logger.Info("Using default configuration")
			return buildConfig(defaultConfig, pp)
		}
		cf = found
		logger.Verbose(fmt.Sprintf("Found config file: %s", cf))
	}

	// Try to read user config
	logger.Verbose(fmt.Sprintf("Loading config from: %s", cf))
	userConfig, err = ReadYamlFile(cf)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", cf, err)
	}

	// Validate config
	result := Validate(userConfig)

	if HasWarnings(result) {
		logger.Warning("Configuration warnings:")
		for _, warning := range result.Warnings {
			logger.Warning("  " + warning)
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
	logger.Verbose(fmt.Sprintf("Config path: %s", cf))

	return buildConfig(userConfig, pp)
}

func FindConfigFile(pp string) (string, error) {
	candidates := []string{
		"psx.yml",
		".psx.yml",
		"psx.yaml",
		".psx.yaml",
	}

	logger.Verbose(fmt.Sprintf("Searching for config in: %s", pp))
	for _, name := range candidates {
		path := filepath.Join(pp, name)
		logger.Verbose(fmt.Sprintf("  Checking: %s", path))
		if _, err := os.Stat(path); err == nil {
			logger.Verbose(fmt.Sprintf("  Found: %s", path))
			return path, nil
		}
	}

	logger.Verbose("Searching upwards for .git directory...")
	current := pp
	for {
		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		gitPath := filepath.Join(current, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			logger.Verbose(fmt.Sprintf("Found .git at: %s", current))
			for _, name := range candidates {
				path := filepath.Join(current, name)
				logger.Verbose(fmt.Sprintf("  Checking: %s", path))
				if _, err := os.Stat(path); err == nil {
					logger.Verbose(fmt.Sprintf("  Found: %s", path))
					return path, nil
				}
			}
			break
		}
		current = parent
	}

	home, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(home, ".config", "psx")
		logger.Verbose(fmt.Sprintf("Checking home config: %s", configDir))
		for _, name := range candidates {
			path := filepath.Join(configDir, name)
			logger.Verbose(fmt.Sprintf("  Checking: %s", path))
			if _, err := os.Stat(path); err == nil {
				logger.Verbose(fmt.Sprintf("  Found: %s", path))
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("no config file found")
}

func buildConfig(uc *Config, pp string) (*Config, error) {
	config := &Config{
		Version:     uc.Version,
		Project:     uc.Project,
		Rules:       uc.Rules,
		Ignore:      uc.Ignore,
		Fix:         uc.Fix,
		Path:        pp,
		ActiveRules: make(map[string]*ActiveRule),
	}

	enabledCount := 0
	disabledCount := 0

	for id, meta := range globalMetadata.Rules {
		var userSev any
		if uc.Rules != nil {
			if val, exists := uc.Rules[id]; exists {
				userSev = val
			}
		}

		severity, err := ParseSeverity(userSev, meta.DefaultSeverity)
		if err != nil {
			logger.Warning(fmt.Sprintf("Rule %s: %v, using default (%s)", id, err, meta.DefaultSeverity))
			severity = &meta.DefaultSeverity
		}

		if severity == nil {
			logger.Verbose(fmt.Sprintf("Rule %s is disabled", id))
			disabledCount++
			continue
		}

		config.ActiveRules[id] = &ActiveRule{
			ID:       id,
			Metadata: meta,
			Severity: *severity,
		}
		enabledCount++
		logger.Verbose(fmt.Sprintf("Rule %s enabled with severity: %s", id, *severity))
	}

	logger.Verbose(fmt.Sprintf("Config built: %d enabled, %d disabled rules", enabledCount, disabledCount))
	return config, nil
}

func ReadYamlFile(path string) (*Config, error) {
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

func GetPatterns(patterns any, projectType string) []string {
	switch p := patterns.(type) {
	case []any:
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
