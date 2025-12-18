package utils

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/m-mdy-m/psx/internal/logger"
)

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

func FileExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return true, info
}
func ReadFile(readFn func(string) ([]byte, error), path string) ([]byte, error) {
	data, err := readFn(path)
	if err != nil {
		return nil, logger.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}

func IsDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
func ReadYAML(path string, out any) error {
	data, err := ReadFile(os.ReadFile, path)
	if err != nil {
		return err
	}
	return UnmarshalYAML(data, path, out)
}
func UnmarshalYAML(data []byte, source string, out any) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return logger.Errorf("failed to unmarshal yaml from %s: %w", source, err)
	}

	logger.Verbosef("Parsed YAML from %s", source)
	return nil
}
func LoadEmbedded[T any](what string, parts string, EFS embed.FS) (*T, error) {
	data, err := ReadFile(EFS.ReadFile, parts)
	if err != nil {
		return nil, err
	}

	var v T
	if err := UnmarshalYAML(data, parts, &v); err != nil {
		return nil, logger.Errorf("failed to parse %s (%s): %w", parts, what, err)
	}

	logger.Verbosef("Loaded %s from %s", what, parts)
	return &v, nil
}
