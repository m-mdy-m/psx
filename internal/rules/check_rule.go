package rules

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/utils"
)

func (c *Checker) CheckRule(rule Rule, spec CheckSpec) CheckResult {
	result := CheckResult{
		RuleID: rule.ID,
	}

	patterns := getPatterns(rule.Patterns, c.ctx.ProjectType)

	validator := spec.Validator
	if validator == nil {
		validator = ValidateExists
	}

	for _, pattern := range patterns {
		fullPath := filepath.Join(c.ctx.ProjectPath, pattern)
		exists, info := utils.FileExists(fullPath)

		if !exists {
			continue
		}

		valid, msg, err := validator(c.ctx, fullPath, info)
		if err != nil {
			result.Passed = false
			result.Message = msg
			return result
		}

		if !valid {
			result.Passed = false
			if msg != "" {
				result.Message = msg
			}
			return result
		}

		result.Passed = true
		if spec.SuccessMessage != "" {
			result.Message = spec.SuccessMessage
		} else {
			result.Message = "Found"
		}
		return result
	}

	result.Passed = false
	if spec.FailMessage != "" {
		result.Message = spec.FailMessage
	}
	return result
}
func getPatterns(patterns map[string][]string, projectType string) []string {
	if p, ok := patterns[projectType]; ok {
		return p
	}
	if p, ok := patterns["*"]; ok {
		return p
	}
	return nil
}
