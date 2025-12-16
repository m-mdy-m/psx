package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

const projectCacheFile = ".psx-project.yml"

func SaveProjectInfo(projectPath string, info *ProjectInfo) error {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := yaml.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal project info: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Project info saved to %s", cachePath))
	return nil
}

func LoadProjectInfo(projectPath string) (*ProjectInfo, error) {
	cachePath := filepath.Join(projectPath, projectCacheFile)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Cache doesn't exist yet
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var info ProjectInfo
	if err := yaml.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Project info loaded from %s", cachePath))
	return &info, nil
}

// GetOrCreateProjectInfo gets project info from cache or creates new one
func GetOrCreateProjectInfo(projectPath string, interactive bool) *ProjectInfo {
	// Try to load from cache first
	cached, err := LoadProjectInfo(projectPath)
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to load project cache: %v", err))
	}

	if cached != nil {
		logger.Verbose("Using cached project info")
		cached.CurrentDir = projectPath // Update current dir
		return cached
	}

	// Cache doesn't exist, create new info
	logger.Verbose("No cached project info found, collecting information...")
	info := GetProjectInfo(projectPath, interactive)

	// Save to cache for next time
	if err := SaveProjectInfo(projectPath, info); err != nil {
		logger.Warning(fmt.Sprintf("Failed to save project cache: %v", err))
	}

	return info
}
