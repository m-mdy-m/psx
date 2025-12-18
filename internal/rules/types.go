package rules

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

type Fixer struct {
	ctx         *Context
	interactive bool
	dryRun      bool
	projectInfo *resources.ProjectInfo
}
type FixSummary struct {
	Total   int
	Fixed   int
	Skipped int
	Failed  int
	Changes int
}

type FixResult struct {
	RuleID  string
	Fixed   bool
	Skipped bool
	Error   error
	Changes []string
}

type CheckResult struct {
	RuleID   string
	Passed   bool
	Severity config.Severity
	Message  string
	FixHint  string
	DocURL   string
}
type Context struct {
	ProjectPath string
	ProjectType string
}
type Rule struct {
	ID       string
	Type     RuleType
	Category string
	Severity config.Severity
	Patterns map[string][]string
	FixSpec  FixSpec
}

type FixSpec struct {
	Type        FixType
	Prompt      string
	CustomCheck func(*Context) (*CheckResult, error)
	CustomFix   func(*Context) error
}

type RuleType int

const (
	RuleTypeFile RuleType = iota
	RuleTypeFolder
	RuleTypeMulti
)

type FixType int

const (
	FixTypeFile FixType = iota
	FixTypeFolder
	FixTypeMulti
)

type CheckSpec struct {
	GetPatterns    func(*config.RuleMetadata, string) []string
	Validator      ValidatorFunc
	SuccessMessage string
	FailMessage    string
}
type Checker struct {
	ctx *Context
}
type ValidatorFunc func(ctx *Context, fullPath string, info any) (bool, string, error)
