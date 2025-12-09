package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/shared/logger"
)
//go:embed embedded/*.yml
var configFS embed.FS

var (
	globalMetadata			*RulesMetadata
	defaultConfig					*Config
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
	cfg,err := ReadYamlFile(cf)
	if err != nil{
		return nil,fmt.Errorf("failed to load config from %s: %w",cf,err)
	}

	result := Validate(cfg)

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
	return cfg,nil
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
	return "", fmt.Errorf("no config file found", pp)
}
func buildConfig(uc *Config, pp string) (*Config, error) {
	config := &Config{
		Version:		uc.Version,
		Project:		uc.Project,
		Rules:			uc.Rules,
		Ignore:			uc.Ignore,
		Fix:			uc.Fix,
	}
	if config.Project.Type == ""{
		logger.Verbose("Project type not specified, will auto-detect")
	}

	for id,meta := range globalMetadata.Rules{
		severity := getSeverity(uc.Rules,id,meta.DefaultSeverity)

		if severity !=nil{
			logger.Verbose(fmt.Sprintf("Rule %s is disabled",id))
			continue
		}
		//	patterns := getPatterns(meta.Patterns,config.Project.Type)
		rulesMetadata := &RuleMetadata{
				ID:				 meta.ID,
				Category:		 meta.Category,
				Description:	 meta.Description,
				DefaultSeverity: meta.DefaultSeverity,
				Patterns:		 meta.Patterns,
				AdditionalChecks:meta.AdditionalChecks,
				Message:		 meta.Message,
				FixHint:		 meta.FixHint,
				DocURL:			 meta.DocURL,
		}
		config.Rules[id]= rulesMetadata
	}
	logger.Verbose(fmt.Sprintf("config with %d active rules",len(config.Rules)))
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

func getSeverity(rules map[string]RulesSeverity, id string, sev Severity) interface{} {
	if userSev,exists := rules[id]; exists{
		if b,ok :=userSev.(bool);ok && !b{
			return nil
		}
		if s,ok := userSev.(string);ok{
			switch s{
			case "error":
				return "error"
			case "warning":
				return "warning"
			case "info":
				return "info"
			default:
				logger.Warning(fmt.Sprintf("Invalid Severity '%s' using 'info-wrnings-error'",s))
				return sev
			}
		}
	}
	return sev
}

func getPatterns(patterns interface{}, projectType string) []string {
	switch p := patterns.(type) {
	case []interface{}:
		result := make([]string, 0, len(p))
		for _, item := range p {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result

	case map[string]interface{}:
		if projectType != "" {
			if langPatterns, exists := p[projectType]; exists {
				if arr, ok := langPatterns.([]interface{}); ok {
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
			if arr, ok := genericPatterns.([]interface{}); ok {
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

