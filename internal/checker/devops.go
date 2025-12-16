package checker

import (
	"github.com/m-mdy-m/psx/internal/config"
)

func CheckDockerfileRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckDockerIgnoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckDockerComposeRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckKubernetesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolderNotEmpty(ctx, metadata)
}

func CheckNginxConfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: func(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
			fileInfo, ok := info.(interface{ IsDir() bool })
			if ok && fileInfo.IsDir() {
				// Check if not empty
				valid, msg, err := ValidateNotEmptyDir(ctx, fullPath, info)
				if !valid {
					return false, msg, err
				}
				return true, "Nginx configuration found", nil
			}
			return true, "Nginx configuration found", nil
		},
	})
}

func CheckInfraFolderRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolder(ctx, metadata)
}
