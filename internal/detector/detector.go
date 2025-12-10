package detector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/shared/logger"
)

var (
	// Global registry of all detectors
	registry *DetectorRegistry
)

func init() {
	// Initialize registry
	registry = NewDetectorRegistry()

	registry.Register(&GoDetector{})

	logger.Verbose(fmt.Sprintf("Registered %d language detectors", len(registry.GetNames())))
}

func Detect(projectPath string) (*DetectionResult, error) {
	logger.Verbose(fmt.Sprintf("Auto-detecting project type in: %s", projectPath))

	// Validate project path
	if err := validateProjectPath(projectPath); err != nil {
		return nil, err
	}

	// Try each detector and collect results
	var bestResult *DetectionResult
	var bestConfidence float64 = 0.0
	var allResults []*DetectionResult

	for _, detector := range registry.All() {
		sig := detector.GetSignature()

		// Quick check first
		if !detector.CanDetect(projectPath) {
			logger.Verbose(fmt.Sprintf("  [%s] Quick check failed", sig.DisplayName))
			continue
		}

		logger.Verbose(fmt.Sprintf("  [%s] Attempting detection...", sig.DisplayName))

		result, err := detector.Detect(projectPath)
		if err != nil {
			logger.Verbose(fmt.Sprintf("  [%s] Detection failed: %v", sig.DisplayName, err))
			continue
		}

		logger.Verbose(fmt.Sprintf("  [%s] Detected with %.0f%% confidence",
			sig.DisplayName, result.Type.Confidence*100))

		allResults = append(allResults, result)

		// Keep track of best result
		if result.Type.Confidence > bestConfidence {
			bestResult = result
			bestConfidence = result.Type.Confidence
		}
	}

	// If no detection succeeded, try generic
	if bestResult == nil {
		logger.Warning("Could not detect specific project type, using generic")
		return detectGeneric(projectPath)
	}

	// Check for mixed projects (multiple languages with high confidence)
	if len(allResults) > 1 {
		bestResult = detectMixedProject(allResults, bestResult)
	}

	logger.Success(fmt.Sprintf("Detected: %s (%.0f%% confidence)",
		bestResult.Description, bestResult.Type.Confidence*100))

	return bestResult, nil
}

func DetectWithHint(projectPath string, hint string) (*DetectionResult, error) {
	logger.Verbose(fmt.Sprintf("Detecting project type with hint: %s", hint))

	// Validate project path
	if err := validateProjectPath(projectPath); err != nil {
		return nil, err
	}

	// Normalize hint
	hint = normalizeLanguageName(hint)

	// Try to find detector for hint
	detector, exists := registry.Get(hint)
	if !exists {
		logger.Warning(fmt.Sprintf("Unknown language hint '%s', falling back to auto-detect", hint))
		logger.Info(fmt.Sprintf("Supported languages: %v", registry.GetNames()))
		return Detect(projectPath)
	}

	sig := detector.GetSignature()
	logger.Verbose(fmt.Sprintf("Using %s detector", sig.DisplayName))

	// Quick check
	if !detector.CanDetect(projectPath) {
		logger.Warning(fmt.Sprintf("Project doesn't appear to be %s", sig.DisplayName))
		logger.Info("Trying auto-detection instead...")
		return Detect(projectPath)
	}

	// Attempt detection
	result, err := detector.Detect(projectPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to detect as %s: %v", sig.DisplayName, err))
		logger.Info("Falling back to auto-detection...")
		return Detect(projectPath)
	}

	// If confidence is low, warn but still return result
	if result.Type.Confidence < 0.7 {
		logger.Warning(fmt.Sprintf("Low confidence detection (%.0f%%) for %s",
			result.Type.Confidence*100, sig.DisplayName))
	}

	logger.Success(fmt.Sprintf("Detected: %s", result.Description))
	return result, nil
}

// validateProjectPath checks if the project path is valid
func validateProjectPath(path string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("failed to access project path: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", path)
	}

	return nil
}

// detectGeneric creates a generic result when no specific language is detected
func detectGeneric(projectPath string) (*DetectionResult, error) {
	result := &DetectionResult{
		Type: ProjectType{
			Primary:    "generic",
			Confidence: 0.3,
			Structure:  "unknown",
			Features:   make(map[string]bool),
		},
		Description: "Generic project (unknown language)",
		Files:       []string{},
		Suggestions: []string{
			"Consider adding a language-specific configuration file (package.json, go.mod, etc.)",
			"Add a README.md to describe your project",
		},
	}

	// Check for some common files to increase confidence slightly
	commonFiles := []string{
		"README.md",
		"LICENSE",
		".gitignore",
		"Makefile",
		"Dockerfile",
	}

	foundCount := 0
	for _, file := range commonFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			result.Files = append(result.Files, file)
			foundCount++
		}
	}

	if foundCount > 0 {
		result.Type.Confidence = 0.3 + (float64(foundCount) * 0.1)
		if result.Type.Confidence > 0.6 {
			result.Type.Confidence = 0.6 // Cap at 0.6 for generic
		}
	}

	return result, nil
}

// detectMixedProject handles projects with multiple languages
func detectMixedProject(allResults []*DetectionResult, primary *DetectionResult) *DetectionResult {
	// Count high-confidence detections
	highConfidence := []*DetectionResult{}
	for _, result := range allResults {
		if result.Type.Confidence >= 0.8 {
			highConfidence = append(highConfidence, result)
		}
	}

	// If only one high-confidence, it's not mixed
	if len(highConfidence) <= 1 {
		return primary
	}

	// Build mixed project result
	mixed := &DetectionResult{
		Type: ProjectType{
			Primary:    primary.Type.Primary,
			Secondary:  []string{},
			Confidence: primary.Type.Confidence,
			Structure:  "mixed",
			Features:   make(map[string]bool),
		},
		Files:       primary.Files,
		Suggestions: []string{},
	}

	// Add secondary languages
	for _, result := range highConfidence {
		if result.Type.Primary != primary.Type.Primary {
			mixed.Type.Secondary = append(mixed.Type.Secondary, result.Type.Primary)
			mixed.Files = append(mixed.Files, result.Files...)
		}
	}

	// Build description
	if len(mixed.Type.Secondary) > 0 {
		mixed.Description = fmt.Sprintf("%s with %v",
			primary.Description, mixed.Type.Secondary)
		mixed.Suggestions = append(mixed.Suggestions,
			"This appears to be a multi-language project (monorepo or polyglot)")
	} else {
		mixed.Description = primary.Description
	}

	return mixed
}

// normalizeLanguageName converts various language name formats to standard form
func normalizeLanguageName(name string) string {
	aliases := map[string]string{
		"node":       "nodejs",
		"javascript": "nodejs",
		"js":         "nodejs",
		"typescript": "nodejs",
		"ts":         "nodejs",
		"golang":     "go",
		"rustlang":   "rust",
		"py":         "python",
		"python3":    "python",
		"rb":         "ruby",
		"php":        "php",
		"java":       "java",
	}

	if standard, exists := aliases[name]; exists {
		return standard
	}

	return name
}

// GetSupportedLanguages returns list of all supported languages
func GetSupportedLanguages() []string {
	names := registry.GetNames()
	result := make([]string, len(names))

	for i, name := range names {
		if detector, exists := registry.Get(name); exists {
			sig := detector.GetSignature()
			result[i] = sig.DisplayName
		} else {
			result[i] = name
		}
	}

	return result
}

// IsSupported checks if a language is supported
func IsSupported(language string) bool {
	normalized := normalizeLanguageName(language)
	_, exists := registry.Get(normalized)
	return exists
}

// GetDetector returns a detector by name
func GetDetector(language string) (Detector, error) {
	normalized := normalizeLanguageName(language)
	detector, exists := registry.Get(normalized)
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return detector, nil
}
func QuickCheck(projectPath string, language string) bool {
	detector, err := GetDetector(language)
	if err != nil {
		return false
	}
	return detector.CanDetect(projectPath)
}

