package checker

import (
	"github.com/m-mdy-m/psx/internal/config"
)

func CheckPrettierRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckPrettierIgnoreRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckESLintRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckCommitlintRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckHuskyRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFolder(ctx, metadata)
}

func CheckLintStagedRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckMakefileRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}

func CheckGitattributesRule(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return CheckFile(ctx, metadata)
}
