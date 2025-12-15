package detector

import (
	"fmt"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

var (
	registry *DetectorRegistry
)

func init() {
	registry = NewDetectorRegistry()
	registry.Register(&GoDetector{})
	logger.Verbose(fmt.Sprintf("Registered %d language detectors", len(registry.GetNames())))
}

func Detect(projectPath string) (*DetectionResult, error) {
	logger.Verbose(fmt.Sprintf("Auto-detecting project type in: %s", projectPath))

	if err := validateProjectPath(projectPath); err != nil {
		return nil, err
	}

	var bestResult *DetectionResult

	for _, detector := range registry.All() {
		sig := detector.GetSignature()

		if !detector.CanDetect(projectPath) {
			logger.Verbose(fmt.Sprintf("  [%s] Quick check failed", sig.Name))
			continue
		}

		logger.Verbose(fmt.Sprintf("  [%s] Attempting detection...", sig.Name))

		result, err := detector.Detect(projectPath)
		if err != nil {
			logger.Verbose(fmt.Sprintf("  [%s] Detection failed: %v", sig.Name, err))
			continue
		}

		if bestResult == nil {
			bestResult = result
		}
	}

	if bestResult == nil {
		logger.Warning("Could not detect specific project type, using generic")
		return detectGeneric(projectPath)
	}

	return bestResult, nil
}

func DetectWithHint(projectPath string, hint string) (*DetectionResult, error) {
	logger.Verbose(fmt.Sprintf("Detecting project type with hint: %s", hint))

	if err := validateProjectPath(projectPath); err != nil {
		return nil, err
	}

	// Normalize hint using resources
	hint = resources.NormalizeLanguage(hint)

	detector, exists := registry.Get(hint)
	if !exists {
		logger.Warning(fmt.Sprintf("Unknown language hint '%s', falling back to auto-detect", hint))
		logger.Info(fmt.Sprintf("Supported languages: %v", registry.GetNames()))
		return Detect(projectPath)
	}

	sig := detector.GetSignature()
	logger.Verbose(fmt.Sprintf("Using %s detector", sig.Name))

	if !detector.CanDetect(projectPath) {
		logger.Warning(fmt.Sprintf("Project doesn't appear to be %s", sig.Name))
		logger.Info("Trying auto-detection instead...")
		return Detect(projectPath)
	}

	result, err := detector.Detect(projectPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to detect as %s: %v", sig.Name, err))
		logger.Info("Falling back to auto-detection...")
		return Detect(projectPath)
	}
	return result, nil
}

func validateProjectPath(path string) error {
	exists, info := shared.FileExists(path)
	if !exists {
		return fmt.Errorf("project path does not exist: %s", path)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", path)
	}
	return nil
}

func detectGeneric(projectPath string) (*DetectionResult, error) {
	result := &DetectionResult{
		Type: ProjectType{
			Primary: "generic",
		},
		Files: []string{},
	}

	commonFiles := []string{
		"README.md",
		"LICENSE",
		".gitignore",
		"Makefile",
		"Dockerfile",
	}

	for _, file := range commonFiles {
		exists, _ := shared.FileExists(filepath.Join(projectPath, file))
		if exists {
			result.Files = append(result.Files, file)
		}
	}

	return result, nil
}

func GetSupportedLanguages() []string {
	names := registry.GetNames()
	result := make([]string, len(names))

	for i, name := range names {
		if detector, exists := registry.Get(name); exists {
			sig := detector.GetSignature()
			result[i] = sig.Name
		} else {
			result[i] = name
		}
	}

	return result
}

func IsSupported(language string) bool {
	normalized := resources.NormalizeLanguage(language)
	_, exists := registry.Get(normalized)
	return exists
}

func GetDetector(language string) (Detector, error) {
	normalized := resources.NormalizeLanguage(language)
	detector, exists := registry.Get(normalized)
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return detector, nil
}