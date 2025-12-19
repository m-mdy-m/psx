package utils

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

func FileExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return true, info
}

func CreateFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return logger.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return logger.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

func CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return logger.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func IsDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
func LoadEmbedded[T any](name, path string, fs embed.FS) (*T, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, logger.Errorf("failed to read %s: %w", name, err)
	}

	var result T
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, logger.Errorf("failed to parse %s: %w", name, err)
	}

	return &result, nil
}
