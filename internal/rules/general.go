package rules

import (
	"os"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
)

// executeReadmeRule checks for README file
func executeReadmeRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if any README file exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil {
			// File exists, check if it's not empty
			if info.Size() < 100 {
				result.Passed = false
				result.Message = "README file exists but appears to be empty or too short"
				result.Details = map[string]any{
					"file": pattern,
					"size": info.Size(),
				}
				return result
			}

			// README found and valid
			result.Passed = true
			result.Message = "README file found"
			result.Details = map[string]any{
				"file": pattern,
			}
			return result
		}
	}

	// No README found
	result.Passed = false
	return result
}

// executeLicenseRule checks for LICENSE file
func executeLicenseRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if any LICENSE file exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil {
			// File exists, check if it's not empty
			if info.Size() < 100 {
				result.Passed = false
				result.Message = "LICENSE file exists but appears to be empty"
				result.Details = map[string]any{
					"file": pattern,
					"size": info.Size(),
				}
				return result
			}

			// LICENSE found and valid
			result.Passed = true
			result.Message = "LICENSE file found"
			result.Details = map[string]any{
				"file": pattern,
			}
			return result
		}
	}

	// Check for license in package.json or other metadata
	if checkAdditionalLicenseLocations(ctx, metadata) {
		result.Passed = true
		result.Message = "License declared in project metadata"
		return result
	}

	// No LICENSE found
	result.Passed = false
	return result
}

// executeGitignoreRule checks for .gitignore file
func executeGitignoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if .gitignore exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if info, err := os.Stat(fullPath); err == nil {
			// File exists, check if it's not empty
			if info.Size() < 10 {
				result.Passed = false
				result.Message = ".gitignore exists but appears to be empty"
				result.Details = map[string]any{
					"file": pattern,
					"size": info.Size(),
				}
				return result
			}

			// .gitignore found and valid
			result.Passed = true
			result.Message = ".gitignore file found"
			result.Details = map[string]any{
				"file": pattern,
			}
			return result
		}
	}

	// No .gitignore found
	result.Passed = false
	return result
}

// executeChangelogRule checks for CHANGELOG file
func executeChangelogRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Location: ctx.ProjectPath,
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if any CHANGELOG file exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		if _, err := os.Stat(fullPath); err == nil {
			result.Passed = true
			result.Message = "CHANGELOG file found"
			result.Details = map[string]any{
				"file": pattern,
			}
			return result
		}
	}

	// No CHANGELOG found (this is usually just info level)
	result.Passed = false
	return result
}

// checkAdditionalLicenseLocations checks for license in metadata files
func checkAdditionalLicenseLocations(ctx *Context, metadata *config.RuleMetadata) bool {
	if len(metadata.AdditionalChecks) == 0 {
		return false
	}
	return false
}
