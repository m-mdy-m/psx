package checker

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/shared"
)

func CheckADRRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,info := shared.FileExists(fullPath)
		if !exists || info == nil {
			continue
		}

		if info.IsDir() {
			result.Passed = true
			result.Message = "ADR folder found"
			return result
		}
	}

	result.Passed = false
	return result
}
func CheckContributingRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,_ := shared.FileExists(fullPath)
		if exists{
			result.Passed = true
			result.Message = "CONTRIBUTING file found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckAPIDocsRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			if info.IsDir() {
				result.Passed = true
				result.Message = "API documentation folder found"
				return result
			}
			result.Passed = true
			result.Message = "API documentation file found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckCIConfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, info := shared.FileExists(fullPath)
		if exists{
			if info.IsDir() {
				isEmpty, _ := shared.IsDirEmpty(fullPath)
				if isEmpty {
					continue
				}
			}

			result.Passed = true
			result.Message = "CI/CD configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckPreCommitRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,_ := shared.FileExists(fullPath)
		if exists{
			result.Passed = true
			result.Message = "Pre-commit hooks configured"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckEditorconfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,_ := shared.FileExists(fullPath)
		if exists{
			result.Passed = true
			result.Message = ".editorconfig found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckCodeOwnersRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message:  metadata.Message,
		FixHint:  metadata.FixHint,
		DocURL:   metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists,_ := shared.FileExists(fullPath)
		if exists{
			result.Passed = true
			result.Message = "CODEOWNERS file found"
			return result
		}
	}

	result.Passed = false
	return result
}


func CheckSecurityRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "SECURITY.md found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckCodeOfConductRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "CODE_OF_CONDUCT.md found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckPullRequestTemplateRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "Pull request template found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckIssueTemplatesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			isEmpty, _ := shared.IsDirEmpty(fullPath)
			if !isEmpty {
				result.Passed = true
				result.Message = "Issue templates found"
				return result
			}
		}
	}

	result.Passed = false
	return result
}

func CheckFundingRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "Funding information found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckSupportRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "SUPPORT.md found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckRoadmapRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "ROADMAP.md found"
			return result
		}
	}

	result.Passed = false
	return result
}