package rules

import (
	"github.com/m-mdy-m/psx/internal/checker"
	"github.com/m-mdy-m/psx/internal/config"
)

type ExecutionResult struct {
	Results []checker.RuleResult
	Summary checker.Summary
	Status  checker.Status
	Context *checker.Context
}
type Executor interface {
	Execute(ctx *checker.Context, metadata *config.RuleMetadata) checker.RuleResult
}

type ExecutorFunc func(ctx *checker.Context, metadata *config.RuleMetadata) checker.RuleResult
func (f ExecutorFunc) Execute(ctx *checker.Context, metadata *config.RuleMetadata) checker.RuleResult {
	return f(ctx, metadata)
}
