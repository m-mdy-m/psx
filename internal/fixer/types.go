package fixer

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

type FixResult struct {
	RuleID  string
	Fixed   bool
	Skipped bool
	Error   error
	Changes []Change
}

// Change represents a single file/folder change
type Change struct {
	Type        ChangeType
	Path        string
	Description string
	Content     string
}

// ChangeType defines the type of change
type ChangeType string

const (
	ChangeCreateFile   ChangeType = "create_file"
	ChangeCreateFolder ChangeType = "create_folder"
	ChangeModifyFile   ChangeType = "modify_file"
	ChangeDeleteFile   ChangeType = "delete_file"
)

// FixContext contains context for fix operations
type FixContext struct {
	ProjectPath   string
	ProjectType   string
	Config        *config.Config
	Interactive   bool
	DryRun        bool
	CreateBackups bool
	ProjectInfo   *resources.ProjectInfo // اضافه شده
}

// Fixer defines the interface for fix operations
type Fixer interface {
	CanFix(ruleID string) bool
	Fix(ctx *FixContext, ruleID string) (*FixResult, error)
}

// FixPlan represents a plan of fixes to apply
type FixPlan struct {
	Fixes        []*FixResult
	TotalChanges int
}

// FixSummary summarizes fix results
type FixSummary struct {
	Total   int
	Fixed   int
	Skipped int
	Failed  int
	Changes int
}