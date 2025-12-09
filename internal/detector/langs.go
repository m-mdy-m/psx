package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================
// Go Detector
// ============================================

type GoDetector struct{}

func (d *GoDetector) GetSignature() LanguageSignature {
	return LanguageSignature{
		Name:         "go",
		DisplayName:  "Go",
		PrimaryFiles: []string{"go.mod"},
		SecondaryFiles: []string{
			"go.sum",
			"go.work",
		},
		Extensions:    []string{".go"},
		MinConfidence: 0.9,
	}
}

func (d *GoDetector) CanDetect(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, "go.mod"))
	return err == nil
}

func (d *GoDetector) Detect(projectPath string) (*DetectionResult, error) {
	result := &DetectionResult{
		Type: ProjectType{
			Primary:    "go",
			Features:   make(map[string]bool),
			Confidence: 0.0,
		},
		Files: []string{},
	}

	// Check go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return nil, fmt.Errorf("not a Go project: go.mod not found")
	}

	result.Files = append(result.Files, "go.mod")
	result.Type.Confidence = 1.0

	// Check go.sum
	if _, err := os.Stat(filepath.Join(projectPath, "go.sum")); err == nil {
		result.Files = append(result.Files, "go.sum")
	}

	// Parse go.mod
	data, err := os.ReadFile(goModPath)
	if err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "go ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					result.Type.Version = parts[1]
					result.Type.Features[fmt.Sprintf("go:%s", parts[1])] = true
				}
			}
		}
	}

	// Detect structure
	d.detectStructure(projectPath, result)

	// Build description
	desc := "Go"
	if result.Type.Version != "" {
		desc += fmt.Sprintf(" %s", result.Type.Version)
	}
	switch result.Type.Structure {
	case "library":
		desc += " library"
	case "application":
		desc += " application"
	}
	result.Description = desc

	return result, nil
}

func (d *GoDetector) detectStructure(projectPath string, result *DetectionResult) {
	// Check for cmd/ (application)
	if _, err := os.Stat(filepath.Join(projectPath, "cmd")); err == nil {
		result.Type.Structure = "application"
		result.Type.Features["cmd"] = true
	} else if _, err := os.Stat(filepath.Join(projectPath, "main.go")); err == nil {
		result.Type.Structure = "application"
	} else {
		result.Type.Structure = "library"
	}

	// Check for internal/
	if _, err := os.Stat(filepath.Join(projectPath, "internal")); err == nil {
		result.Type.Features["internal"] = true
	}

	// Check for pkg/
	if _, err := os.Stat(filepath.Join(projectPath, "pkg")); err == nil {
		result.Type.Features["pkg"] = true
	}
}

