package checker

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/shared"
)

func CheckPrettierRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = "Prettier configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckPrettierIgnoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = ".prettierignore found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckESLintRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = "ESLint configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckCommitlintRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = "Commitlint configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckHuskyRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, info := shared.FileExists(fullPath)
		if exists && info != nil && info.IsDir() {
			result.Passed = true
			result.Message = "Husky git hooks found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckLintStagedRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = "lint-staged configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckMakefileRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = "Makefile found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckGitattributesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, _ := shared.FileExists(fullPath)
		if exists {
			result.Passed = true
			result.Message = ".gitattributes found"
			return result
		}
	}

	result.Passed = false
	return result
}