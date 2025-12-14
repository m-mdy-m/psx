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
	globalMetadata			*RulesMetadata
	defaultConfig			*Config
)
func GetRulesMetadata() *RulesMetadata{
	return globalMetadata
}
func init(){
	var err error
	globalMetadata,err = LoadRulesMetadata()
	if err!=nil{
		logger.Errorf("Failed to load rules metadata: %v",err)
		os.Exit(1)
	}
	logger.Verbose(fmt.Sprintf("Loaded %d rules from metadata", len(globalMetadata.Rules)))
	defaultConfig,err = LoadDefaultConfig()
	if err!=nil{
		logger.Errorf("Failed to load default config: %v",err)
		os.Exit(1)
	}
	logger.Verbose("Default coniguration loaded")
}
func loadEmbedded[T any](what string,parts string) (*T, error) {
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
    md, err := loadEmbedded[RulesMetadata]("rules metadata","embedded/rules.yml")
    if err != nil {
        return nil, err
    }
    logger.Verbose("Loaded rules metadata: " + fmt.Sprint(len(md.Rules)) + " rules")
    return md, nil
}

func LoadDefaultConfig() (*Config, error) {
    cfg, err := loadEmbedded[Config]("default config","embedded/psx.default.yml")
    if err != nil {
        return nil, err
    }
    logger.Verbose("Loaded default configuration")
    return cfg, nil
}

func Load(cf string, pp string)(*Config,error){
	var(
		err			error
		userConfig *Config
	)

	if cf == "" {
		logger.Info("Searching for config file...")
		found, err := FindConfigFile(pp)
		if err == nil && found != "" {
			cf = found
			logger.Verbose(fmt.Sprintf("Found config file: %s",cf))
		} else {
			logger.Info("No config file found, using default")
			return buildConfig(defaultConfig,pp)
		}
	}
	userConfig,err = ReadYamlFile(cf)
	if err != nil{
		return nil,fmt.Errorf("failed to load config from %s: %w",cf,err)
	}

	result := Validate(userConfig)

	if HasWarnings(result) {
		logger.Warning("Configuration warnings:")
		for _, warnings:= range result.Warnings {
			logger.Warning( warnings)
		}
	}
	if !IsValid(result){
		logger.Error("Configuration validation failed:")
		for _,err := range result.Errors{
			logger.Error(fmt.Sprintf("  [%s] %s",err.Field,err.Message))
		}
		return nil, fmt.Errorf("config validation failed: %d errors", len(result.Errors))
	}

	logger.Success("Configuration loaded and validated")
	return buildConfig(userConfig,pp)
}
func FindConfigFile(pp string)(string,error){
	candidates := []string {
		"psx.yml",
		".psx.yml",
		"psx.yaml",
		".psx.yaml",
	}
	for _, name := range candidates {
		path := filepath.Join(pp, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	current := pp
	for {
		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		if _,err := os.Stat(filepath.Join(current,".git")); err == nil{
			for _,name := range candidates{
				path := filepath.Join(current,name)
				if _,err :=os.Stat(path);err == nil{
					return path,nil
				}
			}
			break
		}
		current = parent
	}

	// last hope :)

	home,err := os.UserHomeDir()
	if err ==nil{
		cdir :=filepath.Join(home,".config","psx")
		for _,name := range candidates{
			path := filepath.Join(cdir,name)
			if _,err := os.Stat(path);err==nil{
				return path,nil
			}
		}
	}
	/// bro you don't install
	return "", fmt.Errorf("no config file found in %s", pp)
}
func buildConfig(uc *Config, pp string) (*Config, error) {
	config := &Config{
		Version:		uc.Version,
		Project:		uc.Project,
		Rules:		    uc.Rules,
		Ignore:			uc.Ignore,
		Fix:			uc.Fix,

		Path:           pp,
		ActiveRules:    make(map[string]*ActiveRule),
	}

	enabledCount  :=0
	disabledCount :=0

	for id,meta := range globalMetadata.Rules{
		var userSev any
		if uc.Rules != nil{
			if val,exists := uc.Rules[id];exists{
				userSev = val
			}
		}

		severity, err := ParseSeverity(userSev,meta.DefaultSeverity)
		if err != nil{
			logger.Warning(fmt.Sprintf("Rule %s: %v, using default (%s)",id,err,meta.DefaultSeverity))
			severity = &meta.DefaultSeverity
		}
		if severity == nil{
			logger.Verbose(fmt.Sprintf("Rule %s is disabled",id))
			disabledCount++
			continue
		}
		config.ActiveRules[id] = &ActiveRule{
			ID:				id,
			Metadata:		meta,
			Severity:		*severity,
		}
		enabledCount++
		logger.Verbose(fmt.Sprintf("Rule %s enabled with severity: %s",id,*severity))
	}
	logger.Verbose(fmt.Sprintf("config built: %d enabled, %d disabled rules",enabledCount,disabledCount))
	return config,nil
}

func ReadYamlFile(args string)(*Config,error){
	data,err := os.ReadFile(args)
	if err!=nil{
		return nil,fmt.Errorf("failed to read file: %w",err)
	}

	var cfg Config

	if err:= yaml.Unmarshal([]byte(data),&cfg);err !=nil{
		return nil, fmt.Errorf("failed to parse YAML: %w",err)
	}

	logger.Verbose(fmt.Sprintf("Successfully parsed YAML from %s", args))

	return &cfg,nil
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

