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
	ChangeDeleteFile   ChangeType = "delete_file"
)

type FixContext struct {
	ProjectPath   string
	ProjectType   string
	Config        *config.Config
	Interactive   bool
	DryRun        bool
	CreateBackups bool
	ProjectInfo   *resources.ProjectInfo
}

type Fixer interface {
	CanFix(ruleID string) bool
	Fix(ctx *FixContext, ruleID string) (*FixResult, error)
}

type FixPlan struct {
	Fixes        []*FixResult
	TotalChanges int
}

type FixSummary struct {
	Total   int
	Fixed   int
	Skipped int
	Failed  int
	Changes int
}

type ContentFunc func(ctx *FixContext) (string, error)

type ValidatorFunc func(ctx *FixContext, fullPath string) (bool, error)

type PostCreateFunc func(ctx *FixContext, fullPath string) error

type FileFixSpec struct {
	RuleID       string
	Path         string
	Description  string
	PromptText   string
	GetContent   ContentFunc
	Validator    ValidatorFunc
	PostCreate   PostCreateFunc
	FormatForDry bool
}

type FolderFileSpec struct {
	Name         string
	GetContent   ContentFunc
	PostCreate   PostCreateFunc
	FormatForDry bool
}

type FolderFixSpec struct {
	RuleID      string
	Path        string
	Description string
	PromptText  string
	Files       []FolderFileSpec
	Validator   ValidatorFunc
}

type MultiFileSpec struct {
	Path        string
	Description string
	GetContent  ContentFunc
	PostCreate  PostCreateFunc
}

type MultiFileFixSpec struct {
	RuleID     string
	PromptText string
	Files      []MultiFileSpec
	Validator  ValidatorFunc
}
