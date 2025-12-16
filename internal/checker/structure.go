package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func CheckSrcFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolderNotEmpty(ctx, metadata)
}

func CheckTestsFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	for _, pattern := range patterns {
		// Folder pattern
		if strings.HasSuffix(pattern, "/") {
			folderName := strings.TrimSuffix(pattern, "/")
			fullPath := filepath.Join(ctx.ProjectPath, folderName)
			exists, info := shared.FileExists(fullPath)

			if !exists || info == nil || !info.IsDir() {
				continue
			}

			// For Go, check inline tests
			if ctx.ProjectType == "go" {
				testPatterns := resources.GetTestPatterns(ctx.ProjectType)
				for _, testPattern := range testPatterns {
					hasTests, _ := hasMatchingFiles(ctx.ProjectPath, testPattern)
					if hasTests {
						result.Passed = true
						result.Message = "Test files found (Go tests are inline)"
						return result
					}
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
				return result
			}

			result.Passed = true
			result.Message = "Tests folder found with test files"
			return result
		}

		// Wildcard pattern (e.g., **/*_test.go)
		if strings.Contains(pattern, "*") {
			suffix := pattern
			parts := strings.Split(pattern, "*")
			suffix = parts[len(parts)-1]

			found, count := hasMatchingFiles(ctx.ProjectPath, suffix)
			if found {
				result.Passed = true
				result.Message = fmt.Sprintf("Test files found (%d files)", count)
				return result
			}
		}
	}

	result.Passed = false
	return result
}

func CheckDocsFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolder(ctx, metadata)
}

func CheckScriptsFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolder(ctx, metadata)
}

func CheckEnvExampleRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckGitHubActionsRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolderNotEmpty(ctx, metadata)
}

func CheckRenovateRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckDependabotRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

// Helper functions

func hasTestFiles(path string, projectType string) (bool, error) {
	patterns := resources.GetTestPatterns(projectType)

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
	found := false
	count := 0

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			// Skip common ignored directories
			if name == ".git" || name == "vendor" || name == "node_modules" ||
				name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), suffix) {
			found = true
			count++
		}

		return nil
	})

	return found, count
}
