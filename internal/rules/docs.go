package rules

import (
	"os"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
)

// executeADRRule checks for Architecture Decision Records
func executeADRRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			result.Passed = true
			result.Message = "ADR folder found"
			result.Details = map[string]interface{}{
				"folder": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// executeContributingRule checks for CONTRIBUTING file
func executeContributingRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if _, err := os.Stat(fullPath); err == nil {
			result.Passed = true
			result.Message = "CONTRIBUTING file found"
			result.Details = map[string]interface{}{
				"file": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// executeAPIDocsRule checks for API documentation
func executeAPIDocsRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Only relevant for library projects
	if ctx.Detection.Type.Structure != "library" {
		result.Passed = true
		result.Message = "API docs not required for applications"
		return result
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)

		// Check if it's a folder
		if info, err := os.Stat(fullPath); err == nil {
			if info.IsDir() {
				result.Passed = true
				result.Message = "API documentation folder found"
				result.Details = map[string]interface{}{
					"folder": pattern,
				}
				return result
			} else {
				result.Passed = true
				result.Message = "API documentation file found"
				result.Details = map[string]interface{}{
					"file": pattern,
				}
				return result
			}
		}
	}

	result.Passed = false
	return result
}

// executeCIConfigRule checks for CI/CD configuration
func executeCIConfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil {
			ciType := getCIType(pattern)

			if info.IsDir() {
				// Check if folder has workflow files
				isEmpty, _ := isDirEmpty(fullPath)
				if isEmpty {
					continue
				}
			}

			result.Passed = true
			result.Message = "CI/CD configuration found"
			result.Details = map[string]interface{}{
				"type": ciType,
				"file": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// executePreCommitRule checks for pre-commit hooks
func executePreCommitRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if _, err := os.Stat(fullPath); err == nil {
			result.Passed = true
			result.Message = "Pre-commit hooks configured"
			result.Details = map[string]interface{}{
				"file": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// executeEditorconfigRule checks for .editorconfig
func executeEditorconfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if _, err := os.Stat(fullPath); err == nil {
			result.Passed = true
			result.Message = ".editorconfig found"
			result.Details = map[string]interface{}{
				"file": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// executeCodeOwnersRule checks for CODEOWNERS file
func executeCodeOwnersRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if _, err := os.Stat(fullPath); err == nil {
			result.Passed = true
			result.Message = "CODEOWNERS file found"
			result.Details = map[string]interface{}{
				"file": pattern,
			}
			return result
		}
	}

	result.Passed = false
	return result
}

// Helper functions

func getCIType(pattern string) string {
	if filepath.Base(pattern) == "workflows" || filepath.Dir(pattern) == ".github/workflows" {
		return "GitHub Actions"
	}
	if filepath.Base(pattern) == ".gitlab-ci.yml" {
		return "GitLab CI"
	}
	if filepath.Base(pattern) == "config.yml" && filepath.Dir(pattern) == ".circleci" {
		return "CircleCI"
	}
	if filepath.Base(pattern) == ".travis.yml" {
		return "Travis CI"
	}
	if filepath.Base(pattern) == "Jenkinsfile" {
		return "Jenkins"
	}
	return "Unknown CI"
}
