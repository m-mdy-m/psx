package config

import (
	"fmt"
	"os"

	"path/filepath"

    "github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/shared/logger"
)

func Load(cf string, pp string)(*Config,error){
	if cf == "" {
		logger.Info("Config file is empty — searching in project path")
		found, err := FindConfigFile(pp)
		if err == nil && found != "" {
			cf = found
			logger.Verbose(fmt.Sprintf("Found config file: %s",cf))
		} else {
			logger.Info("No config file found, using default")
			return GetDefaultConfig()
		}
	}
	cfg,err := ReadYamlFile(cf)
	if err != nil{
		return nil,fmt.Errorf("failed to load config from %s: %w",cf,err)
	}

	result := Validate(cfg)

	if HasWarnings(result) {
		logger.Warning("Configuration warnings")
		for _, warnings:= range result.Warnings {
			logger.Warning(fmt.Sprintf(" - %s", warnings))
		}
	}
	if !IsValid(result){
		logger.Error("Configuration validation failed:")
		fmt.Printf("ERROR:%s\n",result.Errors)
		fmt.Printf("Warnings:%s\n",result.Warnings)
		for _,err := range result.Errors{
			logger.Error(fmt.Sprintf("  - %s",err.Message))
		}
		return nil, fmt.Errorf("config validation failed: %d errors", len(result.Errors))
	}

	logger.Verbose(fmt.Sprintf("Successfully loaded and validated config from %s", cf))
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
	return "", fmt.Errorf("no config file found in %s", pp)
}
func GetDefaultConfig() (*Config, error) {
	logger.Verbose("Loading default configuration...")

	cfg,err:= ReadYamlFile("configs", "psx.default.yml")
	if err!= nil{
		return nil,err
	}
	logger.Verbose("Default configurtion loaded successfully")
	return cfg,nil
}

func ReadYamlFile(args ...string)(*Config,error){
    pdy:= filepath.Join(args...)
	data,err := os.ReadFile(pdy)
	if err!=nil{
		 logger.Errorf("failed to read yaml file (%s): %v", pdy, err)
		 return nil,fmt.Errorf("failed to read file %s: %w",pdy,err)
	}

	var cfg Config

	if err:= yaml.Unmarshal([]byte(data),&cfg);err !=nil{
		 logger.Errorf("failed parse yaml (%s): %v",pdy,err)
		 return nil, fmt.Errorf("failed to parse YAML from %s: %w",pdy,err)
	}

	logger.Verbose(fmt.Sprintf("Successfully parsed YAML from %s",pdy))

	return &cfg,nil
}
