package rules

import (
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"github.com/m-mdy-m/psx/internal/config"
)

// executeSrcFolderRule checks for source code folder
func executeSrcFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if any source folder exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			// Check if folder is not empty
			isEmpty, err := isDirEmpty(fullPath)
			if err != nil {
				result.Passed = false
				result.Message = "Source folder exists but cannot check contents"
				return result
			}

			if isEmpty {
				result.Passed = false
				result.Message = "Source folder exists but is empty"
				result.Details = map[string]interface{}{
					"folder": pattern,
				}
				return result
			}

			// Source folder found and valid
			result.Passed = true
			result.Message = "Source folder found"
			result.Details = map[string]interface{}{
				"folder": pattern,
			}
			return result
		}
	}

	// No source folder found
	result.Passed = false
	return result
}

// executeTestsFolderRule checks for tests folder or test files
func executeTestsFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check for test folders or test files
	for _, pattern := range patterns {
		// Check if it's a folder pattern
		if strings.HasSuffix(pattern, "/") {
			// Remove trailing slash for checking
			folderName := strings.TrimSuffix(pattern, "/")
			fullPath := filepath.Join(ctx.ProjectPath, folderName)

			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				// For Go projects, test/ folder might exist but tests are inline
				if ctx.ProjectType == "go" {
					// Check if there are any *_test.go files in the project
					hasGoTests, _ := hasMatchingFiles(ctx.ProjectPath, "_test.go")
					if hasGoTests {
						result.Passed = true
						result.Message = "Test files found (Go tests are inline)"
						result.Details = map[string]interface{}{
							"pattern": "*_test.go",
						}
						return result
					}
				}

				// Check if folder has test files
				hasTests, err := hasTestFiles(fullPath, ctx.ProjectType)
				if err != nil {
					result.Passed = false
					result.Message = "Tests folder exists but cannot check contents"
					return result
				}

				if !hasTests {
					result.Passed = false
					result.Message = "Tests folder exists but contains no test files"
					result.Details = map[string]interface{}{
						"folder": pattern,
					}
					return result
				}

				// Tests found
				result.Passed = true
				result.Message = "Tests folder found with test files"
				result.Details = map[string]interface{}{
					"folder": pattern,
				}
				return result
			}
		} else if strings.Contains(pattern, "*") {
			// It's a glob pattern for test files (e.g., **/*_test.go)
			// Extract the suffix to search for
			suffix := pattern
			if strings.Contains(pattern, "*") {
				parts := strings.Split(pattern, "*")
				suffix = parts[len(parts)-1]
			}

			found, count := hasMatchingFiles(ctx.ProjectPath, suffix)
			if found {
				result.Passed = true
				result.Message = fmt.Sprintf("Test files found (%d files)", count)
				result.Details = map[string]interface{}{
					"pattern": pattern,
					"count":   count,
				}
				return result
			}
		}
	}

	// No tests found
	result.Passed = false
	return result
}

// executeDocsFolderRule checks for documentation folder
func executeDocsFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if any docs folder exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			// Docs folder found
			result.Passed = true
			result.Message = "Documentation folder found"
			result.Details = map[string]interface{}{
				"folder": pattern,
			}
			return result
		}
	}

	// No docs folder found (usually just info level)
	result.Passed = false
	return result
}

// Helper functions

func isDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func hasTestFiles(path string, projectType string) (bool, error) {
	// Test file patterns by language
	testPatterns := map[string][]string{
		"nodejs": {".test.js", ".test.ts", ".spec.js", ".spec.ts"},
		"go":     {"_test.go"},
		"python": {"test_", "_test.py"},
		"rust":   {"test", "tests.rs"},
	}

	patterns, exists := testPatterns[projectType]
	if !exists {
		patterns = []string{"test", "_test", ".test"}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		for _, pattern := range patterns {
			if strings.Contains(name, pattern) {
				return true, nil
			}
		}
	}

	return false, nil
}

func hasMatchingFiles(rootPath string, suffix string) (bool, int) {
	// Simple glob matching - check for suffix in filename

	found := false
	count := 0

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip certain directories
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if filename matches suffix
		if strings.HasSuffix(info.Name(), suffix) {
			found = true
			count++
		}

		return nil
	})

	if err != nil {
		return false, 0
	}

	return found, count
}
