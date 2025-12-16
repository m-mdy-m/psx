package checker

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/shared"
)

func ValidateExists(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
	return true, "", nil
}

func ValidateNotEmpty(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
	fileInfo, ok := info.(interface {
		IsDir() bool
		Size() int64
	})
	if !ok {
		return false, "Cannot validate path", nil
	}

	if fileInfo.IsDir() {
		isEmpty, err := shared.IsDirEmpty(fullPath)
		if err != nil {
			return false, "Cannot check directory contents", err
		}
		if isEmpty {
			return false, "Folder exists but is empty", nil
		}
	} else {
		if fileInfo.Size() < 10 {
			return false, "File exists but appears to be empty", nil
		}
	}

	return true, "", nil
}

func ValidateNotEmptyDir(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
	fileInfo, ok := info.(interface{ IsDir() bool })
	if !ok || !fileInfo.IsDir() {
		return false, "Path is not a directory", nil
	}

	isEmpty, err := shared.IsDirEmpty(fullPath)
	if err != nil {
		return false, "Cannot check directory contents", err
	}
	if isEmpty {
		return false, "Directory exists but is empty", nil
	}

	return true, "", nil
}

func ValidateMinSize(minSize int64) ValidatorFunc {
	return func(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
		fileInfo, ok := info.(interface{ Size() int64 })
		if !ok {
			return false, "Cannot check file size", nil
		}

		if fileInfo.Size() < minSize {
			return false, "File exists but is too small", nil
		}

		return true, "", nil
	}
}

func DefaultGetPatterns(metadata *config.RuleMetadata, projectType string) []string {
	return config.GetPatterns(metadata.Patterns, projectType)
}

func CheckRule(ctx *Context, metadata *config.RuleMetadata, spec CheckSpec) RuleResult {
	result := RuleResult{
		Message: metadata.Message,
		FixHint: metadata.FixHint,
		DocURL:  metadata.DocURL,
	}

	getPatterns := spec.GetPatterns
	if getPatterns == nil {
		getPatterns = DefaultGetPatterns
	}
	patterns := getPatterns(metadata, ctx.ProjectType)

	// Validator
	validator := spec.Validator
	if validator == nil {
		validator = ValidateExists
	}

	// Check each pattern
	for _, pattern := range patterns {
		fullPath := filepath.Join(ctx.ProjectPath, pattern)
		exists, info := shared.FileExists(fullPath)

		if !exists {
			continue
		}

		// Run custom validator
		valid, msg, err := validator(ctx, fullPath, info)
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

		// Success
		result.Passed = true
		if spec.SuccessMessage != "" {
			result.Message = spec.SuccessMessage
		} else {
			result.Message = "Found"
		}
		return result
	}

	// Not found
	result.Passed = false
	if spec.FailMessage != "" {
		result.Message = spec.FailMessage
	}
	return result
}

func CheckFile(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: ValidateExists,
	})
}

func CheckFileNotEmpty(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: ValidateNotEmpty,
	})
}

func CheckFolder(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: ValidateExists,
	})
}

func CheckFolderNotEmpty(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: ValidateNotEmptyDir,
	})
}

func CheckMinSize(minSize int64) func(*Context, *config.RuleMetadata) RuleResult {
	return func(ctx *Context, metadata *config.RuleMetadata) RuleResult {
		return CheckRule(ctx, metadata, CheckSpec{
			Validator: ValidateMinSize(minSize),
		})
	}
}
