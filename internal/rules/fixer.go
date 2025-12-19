package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/utils"
)

type Fixer struct {
	ctx *Context
}
func NewFixer(ctx *Context) *Fixer {
	return &Fixer{ctx: ctx}
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

	return results, nil
}

func (f *Fixer) fix(ruleID string, rule *config.ActiveRule, fixCtx *FixContext) (*FixResult, error) {
	logger.Verbose(fmt.Sprintf("Fixing: %s", ruleID))

	patterns := config.GetPatterns(rule.Metadata.Patterns, f.ctx.ProjectType)
	if len(patterns) == 0 {
		return &FixResult{
			RuleID:  ruleID,
			Skipped: true,
		}, nil
	}

	primaryPattern := patterns[0]
	fullPath := filepath.Join(f.ctx.ProjectPath, primaryPattern)

	exists, info := utils.FileExists(fullPath)
	if exists {
		if info.IsDir() {
			if isEmpty, _ := utils.IsDirEmpty(fullPath); !isEmpty {
				return &FixResult{RuleID: ruleID, Skipped: true}, nil
			}
		} else if info.Size() > 0 {
			return &FixResult{RuleID: ruleID, Skipped: true}, nil
		}
	}

	// Ask user if interactive
	if fixCtx.Interactive && !fixCtx.DryRun {
		prompt := fmt.Sprintf("Create %s?", primaryPattern)
		if !utils.Prompt(prompt) {
			return &FixResult{RuleID: ruleID, Skipped: true}, nil
		}
	}

	// Determine fix type
	isFolder := strings.HasSuffix(primaryPattern, "/") ||
		!strings.Contains(filepath.Base(primaryPattern), ".")

	var changes []Change
	var err error

	if fixCtx.DryRun {
		// Dry run - just preview
		changeType := ChangeCreateFile
		if isFolder {
			changeType = ChangeCreateFolder
		}

		content := ""
		if !isFolder {
			content = f.generateContent(ruleID, primaryPattern)
			content = formatContent(content, 10)
		}

		changes = append(changes, Change{
			Type:        changeType,
			Path:        fullPath,
			Description: fmt.Sprintf("Create %s", primaryPattern),
			Content:     content,
		})
	} else {
		if isFolder {
			err = utils.CreateDir(fullPath)
			if err == nil {
				changes = append(changes, Change{
					Type:        ChangeCreateFolder,
					Path:        fullPath,
					Description: fmt.Sprintf("Created %s", primaryPattern),
				})
			}
		} else {
			content := f.generateContent(ruleID, primaryPattern)
			err = utils.CreateFile(fullPath, content)
			if err == nil {
				changes = append(changes, Change{
					Type:        ChangeCreateFile,
					Path:        fullPath,
					Description: fmt.Sprintf("Created %s", primaryPattern),
				})
			}
		}
	}

	if err != nil {
		return &FixResult{
			RuleID: ruleID,
			Error:  err,
		}, err
	}

	return &FixResult{
		RuleID:  ruleID,
		Fixed:   true,
		Changes: changes,
	}, nil
}

func (f *Fixer) generateContent(ruleID, pattern string) string {
	info := f.ctx.ProjectInfo
	projectType := f.ctx.ProjectType

	switch ruleID {
	case "readme":
		return resources.GetReadme(info, projectType)
	case "license":
		return resources.GetLicense(info.License, info.Author)
	case "gitignore":
		return resources.GetGitignore(projectType)
	case "changelog":
		return resources.GetChangelog(info)
	case "contributing":
		return resources.GetContributing()
	case "security":
		return resources.GetSecurity(info)
	case "code_of_conduct":
		return resources.GetCodeOfConduct(info)
	case "editorconfig":
		return resources.GetEditorconfig(projectType)
	case "dockerfile":
		return resources.GetDockerfile(info, projectType)
	case "dockerignore":
		return resources.GetDockerignore(projectType)
	default:
		return fmt.Sprintf("# %s\n\nTODO: Add content\n", filepath.Base(pattern))
	}
}

func formatContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	preview := strings.Join(lines[:maxLines], "\n")
	remaining := len(lines) - maxLines
	return fmt.Sprintf("%s\n... (%d more lines)", preview, remaining)
}
