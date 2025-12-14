package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GoDetector struct{}

func (d *GoDetector) GetSignature() LanguageSignature {
	return LanguageSignature{
		Name:         "Golang",
		PrimaryFiles: []string{"go.mod"},
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
		},
		Files: []string{},
	}

	// Check go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return nil, fmt.Errorf("not a Go project: go.mod not found")
	}

	result.Files = append(result.Files, "go.mod")

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

	return result, nil
}

