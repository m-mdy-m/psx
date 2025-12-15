package fixer

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixEditorconfig(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "editorconfig",
		Changes: []Change{},
	}

	editorconfigPath := filepath.Join(ctx.ProjectPath, ".editorconfig")

	exists, _ := shared.FileExists(editorconfigPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetEditorconfig()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        editorconfigPath,
			Description: "Create .editorconfig",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .editorconfig?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(editorconfigPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        editorconfigPath,
		Description: "Created .editorconfig",
	})

	return result, nil
}
func FixPreCommit(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "pre_commit",
		Changes: []Change{},
	}

	precommitPath := filepath.Join(ctx.ProjectPath, ".pre-commit-config.yaml")

	exists, _ := shared.FileExists(precommitPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Basic pre-commit config
	content := `# Pre-commit hooks
# See https://pre-commit.com for more information

repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
      - id: check-merge-conflict
`

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        precommitPath,
			Description: "Create .pre-commit-config.yaml",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create pre-commit config?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(precommitPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        precommitPath,
		Description: "Created .pre-commit-config.yaml",
	})

	return result, nil
}

func FixCodeOwners(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "code_owners",
		Changes: []Change{},
	}

	codeownersPath := filepath.Join(ctx.ProjectPath, ".github", "CODEOWNERS")

	exists, _ := shared.FileExists(codeownersPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := `# Code owners
# See https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners

# Global owners
* @owner

# Specific directories
# /docs/ @docs-team
# /src/ @dev-team
`

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        codeownersPath,
			Description: "Create CODEOWNERS",
			Content:     FormatContent(content, 8),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CODEOWNERS?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(codeownersPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        codeownersPath,
		Description: "Created CODEOWNERS",
	})

	return result, nil
}