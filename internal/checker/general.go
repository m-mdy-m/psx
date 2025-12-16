package checker

import (
	"github.com/m-mdy-m/psx/internal/config"
)

func CheckReadmeRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckMinSize(100)(ctx, metadata)
}

func CheckLicenseRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckMinSize(100)(ctx, metadata)
}

func CheckGitignoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckMinSize(10)(ctx, metadata)
}

func CheckChangelogRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}
