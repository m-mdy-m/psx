package rules

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/tree"
)

// Outcome is the aggregate verdict of a whole run.
type Outcome string

const (
	OutcomePassed   Outcome = "passed"
	OutcomeWarnings Outcome = "warnings"
	OutcomeFailed   Outcome = "failed"
)

// Status is the verdict for a single rule.
type Status string

const (
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
	StatusSkipped Status = "skipped"
)

// Context carries everything a rule evaluation needs.
type Context struct {
	ProjectPath string
	ProjectType string
	ProjectInfo *resources.ProjectInfo
	Config      *config.Config
}

// RuleResult is the outcome of one rule against one scope.
type RuleResult struct {
	RuleID   string
	Status   Status
	Severity config.Severity
	Message  string
	FixHint  string
	DocURL   string
	Evidence string
}

// Summary aggregates rule results by status and severity.
type Summary struct {
	Total    int
	Passed   int
	Failed   int
	Skipped  int
	Errors   int
	Warnings int
	Info     int
}

// ExecutionResult is the full outcome of a run.
type ExecutionResult struct {
	Context *Context
	Results []RuleResult
	Summary Summary
	Status  Outcome
}

// Passed reports whether the rule satisfied its requirement.
func (r RuleResult) Passed() bool { return r.Status == StatusPassed }

// ChangeType enumerates the mutations a fix can perform.
type ChangeType string

const (
	ChangeCreateFile   ChangeType = "create_file"
	ChangeCreateFolder ChangeType = "create_folder"
	ChangePatchFile    ChangeType = "patch_file"
)

// Change is a single planned mutation.
type Change struct {
	Type        ChangeType
	Path        string
	Description string
	Content     string
	Mode        uint32
}

// Path returns the project-relative path of a change.
func (c Change) Rel(root string) string {
	if rel, err := filepath.Rel(root, c.Path); err == nil {
		return filepath.ToSlash(rel)
	}
	return c.Path
}

// FixResult reports what happened for one rule during a fix run.
type FixResult struct {
	RuleID  string
	Fixed   bool
	Skipped bool
	Error   error
	Reason  string
	Changes []Change
}

// FixContext configures a fix run.
type FixContext struct {
	Context       *Context
	Snapshot      *tree.Snapshot
	Interactive   bool
	DryRun        bool
	CreateBackups bool
	Force         bool
}
