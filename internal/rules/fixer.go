package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

type Fixer struct {
	ctx       *Context
	generator *ContentGenerator
	resolver  *PatternResolver
}

func NewFixer(ctx *Context) *Fixer {
	return &Fixer{
		ctx:       ctx,
		generator: NewContentGenerator(ctx.ProjectInfo, ctx.ProjectType),
		resolver:  NewPatternResolver(),
	}
}

func Fix(cfg *config.Config, fixCtx *FixContext, ruleID string) (*FixResult, error) {
	fixer := NewFixer(fixCtx.Context)
	activeRule, exists := cfg.ActiveRules[ruleID]
	if !exists {
		return nil, fmt.Errorf("rule not found: %s", ruleID)
	}

	return fixer.fix(ruleID, activeRule, fixCtx)
}

func FixAll(cfg *config.Config, fixCtx *FixContext, failedRules []string) ([]*FixResult, error) {
	fixer := NewFixer(fixCtx.Context)
	results := make([]*FixResult, 0, len(failedRules))

	for _, ruleID := range failedRules {
		activeRule, exists := cfg.ActiveRules[ruleID]
		if !exists {
			logger.Warning(fmt.Sprintf("Rule not found: %s", ruleID))
			continue
		}
		result, err := fixer.fix(ruleID, activeRule, fixCtx)
		if err != nil {
			logger.Warning(fmt.Sprintf("Fix failed for %s: %v", ruleID, err))
			continue
		}

		results = append(results, result)
	}

	if cfg.Custom != nil {
		handler := NewCustomHandler(fixCtx.Context, cfg.Custom)
		customResults, err := handler.ApplyCustomFiles(fixCtx)
		if err != nil {
			return results, err
		}
		results = append(results, customResults...)
	}
	return results, nil
}

func (f *Fixer) fix(ruleID string, rule *config.ActiveRule, fixCtx *FixContext) (*FixResult, error) {
	logger.Verbose(fmt.Sprintf("Fixing: %s", ruleID))

	if f.resolver.IsSpecialMultiFileRule(ruleID) || f.needsMultiFileGeneration(ruleID) {
		multiFiles, err := f.generator.GenerateMultiple(ruleID)
		if err == nil && len(multiFiles) > 0 {
			return f.fixMultipleFiles(ruleID, multiFiles, fixCtx)
		}
	}

	patterns := config.GetPatterns(rule.Metadata.Patterns, f.ctx.ProjectType)
	if len(patterns) == 0 {
		return &FixResult{
			RuleID:  ruleID,
			Skipped: true,
		}, nil
	}

	primaryPattern := patterns[0]
	return f.fixSinglePattern(ruleID, primaryPattern, fixCtx)
}

func (f *Fixer) needsMultiFileGeneration(ruleID string) bool {
	multiFileRules := []string{
		"api_docs",
		"src_folder",
		"tests_folder",
		"ci_config",
		"issue_templates",
		"adr",
		"scripts_folder",
	}

	for _, rule := range multiFileRules {
		if ruleID == rule {
			return true
		}
	}

	return false
}

func (f *Fixer) fixSinglePattern(ruleID, pattern string, fixCtx *FixContext) (*FixResult, error) {
	fullPath := filepath.Join(f.ctx.ProjectPath, pattern)

	if f.shouldSkipPattern(fullPath) {
		return &FixResult{RuleID: ruleID, Skipped: true}, nil
	}

	isFolder := f.isFolder(pattern)

	if fixCtx.Interactive && !fixCtx.DryRun {
		resourceType := "file"
		prompt := fmt.Sprintf("Create %s %s?", resourceType, pattern)
		if !utils.Prompt(prompt) {
			return &FixResult{RuleID: ruleID, Skipped: true}, nil
		}
	}

	var changes []Change
	var err error

	if fixCtx.DryRun {
		changes = f.previewChanges(ruleID, pattern, fullPath, isFolder)
	} else {
		changes, err = f.applyChanges(ruleID, pattern, fullPath, isFolder)
		if err != nil {
			return &FixResult{
				RuleID: ruleID,
				Error:  err,
			}, err
		}
	}

	return &FixResult{
		RuleID:  ruleID,
		Fixed:   true,
		Changes: changes,
	}, nil
}

func (f *Fixer) fixMultipleFiles(ruleID string, files map[string]string, fixCtx *FixContext) (*FixResult, error) {
	var changes []Change
	var errors []string

	if fixCtx.Interactive && !fixCtx.DryRun {
		prompt := fmt.Sprintf("Create %d files for %s?", len(files), ruleID)
		if !utils.Prompt(prompt) {
			return &FixResult{RuleID: ruleID, Skipped: true}, nil
		}
	}

	for relPath, content := range files {
		fullPath := filepath.Join(f.ctx.ProjectPath, relPath)

		if f.shouldSkipPattern(fullPath) {
			continue
		}

		if fixCtx.DryRun {
			changes = append(changes, Change{
				Type:        ChangeCreateFile,
				Path:        fullPath,
				Description: fmt.Sprintf("Create %s", relPath),
				Content:     formatContent(content, 10),
			})
		} else {
			err := utils.CreateFile(fullPath, content)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", relPath, err))
				continue
			}

			changes = append(changes, Change{
				Type:        ChangeCreateFile,
				Path:        fullPath,
				Description: fmt.Sprintf("Created %s", relPath),
			})
		}
	}

	if len(errors) > 0 {
		return &FixResult{
			RuleID: ruleID,
			Error:  fmt.Errorf("failed to create some files: %s", strings.Join(errors, "; ")),
		}, nil
	}

	return &FixResult{
		RuleID:  ruleID,
		Fixed:   true,
		Changes: changes,
	}, nil
}

func (f *Fixer) shouldSkipPattern(fullPath string) bool {
	exists, info := utils.FileExists(fullPath)
	if !exists {
		return false
	}

	// If it's a directory, check if empty
	if info.IsDir() {
		if isEmpty, _ := utils.IsDirEmpty(fullPath); !isEmpty {
			return true
		}
		return false
	}

	// If it's a file, check if has content
	if info.Size() > 0 {
		return true
	}

	return false
}

func (f *Fixer) isFolder(pattern string) bool {
	return f.resolver.ResolveType(pattern) == PatternTypeFolder
}

func (f *Fixer) previewChanges(ruleID, pattern, fullPath string, isFolder bool) []Change {
	changes := []Change{}

	if isFolder {
		changes = append(changes, Change{
			Type:        ChangeCreateFolder,
			Path:        fullPath,
			Description: fmt.Sprintf("Create %s", pattern),
			Content:     "",
		})
	} else {
		content, _ := f.generator.Generate(ruleID, pattern)
		changes = append(changes, Change{
			Type:        ChangeCreateFile,
			Path:        fullPath,
			Description: fmt.Sprintf("Create %s", pattern),
			Content:     formatContent(content, 10),
		})
	}

	return changes
}

func (f *Fixer) applyChanges(ruleID, pattern, fullPath string, isFolder bool) ([]Change, error) {
	changes := []Change{}

	if isFolder {
		err := utils.CreateDir(fullPath)
		if err != nil {
			return nil, err
		}

		changes = append(changes, Change{
			Type:        ChangeCreateFolder,
			Path:        fullPath,
			Description: fmt.Sprintf("Created %s", pattern),
		})
	} else {
		content, err := f.generator.Generate(ruleID, pattern)
		if err != nil {
			return nil, err
		}

		err = utils.CreateFile(fullPath, content)
		if err != nil {
			return nil, err
		}

		changes = append(changes, Change{
			Type:        ChangeCreateFile,
			Path:        fullPath,
			Description: fmt.Sprintf("Created %s", pattern),
		})
	}

	return changes, nil
}

func formatContent(content string, maxLines int) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	preview := strings.Join(lines[:maxLines], "\n")
	remaining := len(lines) - maxLines
	return fmt.Sprintf("%s\n... (%d more lines)", preview, remaining)
}
