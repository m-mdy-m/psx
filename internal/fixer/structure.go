// strcutre.go:
package fixer

import (
	"fmt"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/shared"
)

// FixSrcFolder creates source folder
func FixSrcFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "src_folder",
		Changes: []Change{},
	}

	folderName := getSrcFolderName(ctx.ProjectType)
	folderPath := filepath.Join(ctx.ProjectPath, folderName)

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: fmt.Sprintf("Create %s/ folder", folderName),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt(fmt.Sprintf("Create %s/ folder?", folderName)) {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create .gitkeep to track empty folder
	gitkeepPath := filepath.Join(folderPath, ".gitkeep")
	if err := shared.CreateFile(gitkeepPath, ""); err != nil {
		// Non-critical error, continue
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: fmt.Sprintf("Created %s/ folder", folderName),
	})

	return result, nil
}

// FixTestsFolder creates tests folder
func FixTestsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "tests_folder",
		Changes: []Change{},
	}

	folderName := getTestsFolderName(ctx.ProjectType)
	folderPath := filepath.Join(ctx.ProjectPath, folderName)

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: fmt.Sprintf("Create %s/ folder", folderName),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt(fmt.Sprintf("Create %s/ folder?", folderName)) {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create example test file
	exampleTest := getExampleTestFile(ctx.ProjectType)
	if exampleTest != "" {
		testPath := filepath.Join(folderPath, exampleTest)
		content := getExampleTestContent(ctx.ProjectType)
		if err := shared.CreateFile(testPath, content); err != nil {
			// Non-critical error, continue
		}
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: fmt.Sprintf("Created %s/ folder", folderName),
	})

	return result, nil
}

// FixDocsFolder creates documentation folder
func FixDocsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "docs_folder",
		Changes: []Change{},
	}

	folderPath := filepath.Join(ctx.ProjectPath, "docs")

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: "Create docs/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create docs/ folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create docs README
	readmePath := filepath.Join(folderPath, "README.md")
	content := `# Documentation Tesst Files`
	if err := shared.CreateFile(readmePath, content); err != nil {
		// Non-critical error, continue
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: "Created docs/ folder",
	})

	return result, nil
}

// Helper functions
func getSrcFolderName(projectType string) string {
	folders := map[string]string{
		"nodejs": "src",
		"go":     "cmd",
		"python": "src",
		"rust":   "src",
	}

	if folder, ok := folders[projectType]; ok {
		return folder
	}
	return "src"
}

func getTestsFolderName(projectType string) string {
	folders := map[string]string{
		"nodejs": "tests",
		"go":     "tests",
		"python": "tests",
		"rust":   "tests",
	}

	if folder, ok := folders[projectType]; ok {
		return folder
	}
	return "tests"
}

func getExampleTestFile(projectType string) string {
	files := map[string]string{
		"nodejs": "example.test.js",
		"go":     "example_test.go",
		"python": "test_example.py",
		"rust":   "example_test.rs",
	}

	if file, ok := files[projectType]; ok {
		return file
	}
	return ""
}

func getExampleTestContent(projectType string) string {
	templates := map[string]string{
		"nodejs": `// Example test file
describe('Example', () => {
  it('should pass', () => {
    expect(true).toBe(true);
  });
});
`,
		"go": `package tests

import "testing"

func TestExample(t *testing.T) {
	// Your test here
	if true != true {
		t.Error("This should not fail")
	}
}
`,
		"python": `"""Example test file"""

def test_example():
    """Example test"""
    assert True
`,
	}

	if template, ok := templates[projectType]; ok {
		return template
	}
	return ""
}

