package rules

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
)
// Status represents overall check status
type Status string
const (
	StatusPassed   Status = "passed"
	StatusWarnings Status = "warnings"
	StatusFailed   Status = "failed"
)
// Context provides information needed by rules
type Context struct {
	ProjectPath string
	ProjectType string
	Detection   *detector.DetectionResult
}
// RuleResult represents the result of executing a single rule
type RuleResult struct {
	RuleID      string
	Passed      bool
	Severity    config.Severity
	Message     string
	Location    string
	FixHint     string
	DocURL      string
	Details     map[string]any
}
// Summary contains execution statistics
type Summary struct {
	Total    int
	Passed   int
	Errors   int
	Warnings int
	Info     int
}
// ExecutionResult contains results from all rules
type ExecutionResult struct {
	Results []RuleResult
	Summary Summary
	Status  Status
	Context *Context
}
// Executor interface for rule execution
type Executor interface {
	Execute(ctx *Context, metadata *config.RuleMetadata) RuleResult
}

// ExecutorFunc is a function that implements Executor
type ExecutorFunc func(ctx *Context, metadata *config.RuleMetadata) RuleResult
// Execute implements the Executor interface
func (f ExecutorFunc) Execute(ctx *Context, metadata *config.RuleMetadata) RuleResult {
	return f(ctx, metadata)
}
