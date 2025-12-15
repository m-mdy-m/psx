package checker

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/shared"
)

func CheckDockerfileRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "Dockerfile found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckDockerIgnoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = ".dockerignore found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckDockerComposeRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "Docker Compose configuration found"
			return result
		}
	}

	result.Passed = false
	return result
}

func CheckKubernetesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
				result.Message = "Kubernetes configurations found"
				return result
			}
		}
	}

	result.Passed = false
	return result
}

func CheckNginxConfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	patterns := config.GetPatterns(metadata.Patterns, ctx.ProjectType)

	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, info := shared.FileExists(fullPath)
		if exists {
			if info != nil && info.IsDir() {
				isEmpty, _ := shared.IsDirEmpty(fullPath)
				if !isEmpty {
					result.Passed = true
					result.Message = "Nginx configuration found"
					return result
				}
			} else {
				result.Passed = true
				result.Message = "Nginx configuration found"
				return result
			}
		}
	}

	result.Passed = false
	return result
}

func CheckInfraFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
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
			result.Message = "Infrastructure folder found"
			return result
		}
	}

	result.Passed = false
	return result
}