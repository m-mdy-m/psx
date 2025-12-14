package checker

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/shared"
)

func CheckReadmeRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,info := shared.FileExists(fullPath)
		if exists{
			// File exists, check if it's not empty
			if info.Size() < 100 {
				result.Passed = false
				result.Message = "README file exists but appears to be empty or too short"
				return result
			}

			// README found and valid
			result.Passed = true
			result.Message = "README file found"
			return result
		}
	}

	// No README found
	result.Passed = false
	return result
}

func CheckLicenseRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exist, info := shared.FileExists(fullPath)
		if exist {
			if info.Size() < 100 {
				result.Passed = false
				result.Message = "LICENSE file exists but appears to be empty"
				return result
			}

			result.Passed = true
			result.Message = "LICENSE file found"
			return result
		}
	}


	// No LICENSE found
	result.Passed = false
	return result
}
func CheckGitignoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	// Check if .gitignore exists
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exist, info := shared.FileExists(fullPath)
		if exist{
			if info.Size() < 10 {
				result.Passed = false
				result.Message = ".gitignore exists but appears to be empty"
				return result
			}

			// .gitignore found and valid
			result.Passed = true
			result.Message = ".gitignore file found"
			return result
		}
	}

	// No .gitignore found
	result.Passed = false
	return result
}

func CheckChangelogRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exist, _ := shared.FileExists(fullPath)
		if exist {
			result.Passed = true
			result.Message = "CHANGELOG file found"
			return result
		}
	}

	result.Passed = false
	return result
}
