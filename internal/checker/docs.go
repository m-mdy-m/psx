package checker

import (
	"github.com/m-mdy-m/psx/internal/config"
)

func CheckADRRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolder(ctx, metadata)
}

func CheckContributingRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckAPIDocsRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: func(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
			fileInfo, ok := info.(interface{ IsDir() bool })
			if ok && fileInfo.IsDir() {
				return true, "API documentation folder found", nil
			}
			return true, "API documentation file found", nil
		},
	})
}

func CheckCIConfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckRule(ctx, metadata, CheckSpec{
		Validator: func(ctx *Context, fullPath string, info interface{}) (bool, string, error) {
			fileInfo, ok := info.(interface{ IsDir() bool })
			if ok && fileInfo.IsDir() {
				// Check if not empty
				valid, msg, err := ValidateNotEmptyDir(ctx, fullPath, info)
				if !valid {
					return false, msg, err
				}
			}
			return true, "CI/CD configuration found", nil
		},
	})
}

func CheckPreCommitRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckEditorconfigRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckCodeOwnersRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckSecurityRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckCodeOfConductRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckPullRequestTemplateRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckIssueTemplatesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolderNotEmpty(ctx, metadata)
}

func CheckFundingRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckSupportRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckRoadmapRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}
