package rules

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

type Checker struct {
	ctx *Context
}
type Context struct {
	ProjectPath string
	ProjectType string
	ProjectInfo *resources.ProjectInfo
	Config      *config.Config
}
type RuleResult struct {
	RuleID   string
	Passed   bool
	Severity config.Severity
	Message  string
	FixHint  string
	DocURL   string
}
type ExecutionResult struct {
	Context *Context
	Results []RuleResult
	Summary Summary
	Status  Status
}
type Summary struct {
	Total    int
	Passed   int
	Errors   int
	Warnings int
	Info     int
}
type Status string

const (
	StatusPassed   Status = "passed"
	StatusWarnings Status = "warnings"
	StatusFailed   Status = "failed"
)

type FixResult struct {
	RuleID  string
	Fixed   bool
	Skipped bool
	Error   error
	Changes []Change
}
type Change struct {
	Type        ChangeType
	Path        string
	Description string
	Content     string
}
type ChangeType string

const (
	ChangeCreateFile   ChangeType = "create_file"
	ChangeCreateFolder ChangeType = "create_folder"
	ChangeModifyFile   ChangeType = "modify_file"
)

type FixContext struct {
	Context       *Context
	Interactive   bool
	DryRun        bool
	CreateBackups bool
}
