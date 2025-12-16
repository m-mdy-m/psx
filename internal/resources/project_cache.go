package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

const projectCacheFile = ".psx-project.yml"


func SaveProjectCache(projectPath string, info *ProjectInfo, detectionType string, version string, features map[string]bool, files []string) error {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	cache := &ProjectCache{
		ProjectInfo: info,
		Detection: &DetectionCache{
			ProjectType: detectionType,
			Version:     version,
			Features:    features,
			Files:       files,
		},
	}

	data, err := yaml.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal project cache: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Project cache saved to %s", cachePath))
	return nil
}

func LoadProjectCache(projectPath string) (*ProjectCache, error) {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Cache doesn't exist yet
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cache ProjectCache
	if err := yaml.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Project cache loaded from %s", cachePath))
	return &cache, nil
}

func SaveProjectInfo(projectPath string, info *ProjectInfo) error {
	cache, _ := LoadProjectCache(projectPath)
	if cache == nil {
		cache = &ProjectCache{
			ProjectInfo: info,
			Detection:   &DetectionCache{},
		}
	} else {
		cache.ProjectInfo = info
	}

	cachePath := filepath.Join(projectPath, projectCacheFile)
	data, err := yaml.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal project cache: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Project info saved to %s", cachePath))
	return nil
}

func LoadProjectInfo(projectPath string) (*ProjectInfo, error) {
	cache, err := LoadProjectCache(projectPath)
	if err != nil {
		return nil, err
	}
	if cache == nil || cache.ProjectInfo == nil {
		return nil, nil
	}
	return cache.ProjectInfo, nil
}

func GetOrCreateProjectInfo(projectPath string, interactive bool) *ProjectInfo {
	cached, err := LoadProjectInfo(projectPath)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to load project cache: %v", err))
	}

	if cached != nil {
		logger.Verbose("Using cached project info")
		cached.CurrentDir = projectPath
		return cached
	}

	logger.Verbose("No cached project info found, collecting information...")
	info := GetProjectInfo(projectPath, interactive)
	if err := SaveProjectInfo(projectPath, info); err != nil {
		logger.Warning(fmt.Sprintf("Failed to save project cache: %v", err))
	}

	return info
}
