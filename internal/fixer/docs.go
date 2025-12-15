package fixer

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixADR(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "adr",
		Changes: []Change{},
	}

	adrPath := filepath.Join(ctx.ProjectPath, "docs", "adr")

	exists, _ := shared.FileExists(adrPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        adrPath,
			Description: "Create docs/adr/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create ADR folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(adrPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create first ADR
	firstADR := filepath.Join(adrPath, "0001-record-architecture-decisions.md")
	content := resources.GetFirstADR()
	if err := shared.CreateFile(firstADR, content); err != nil {
		result.Error = err
		return result, err
	}

	// Create template
	template := filepath.Join(adrPath, "template.md")
	templateContent := resources.GetADRTemplate()
	shared.CreateFile(template, templateContent)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        adrPath,
		Description: "Created docs/adr/ folder with initial ADR",
	})

	return result, nil
}

func FixContributing(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "contributing",
		Changes: []Change{},
	}

	contributingPath := filepath.Join(ctx.ProjectPath, "CONTRIBUTING.md")

	exists, _ := shared.FileExists(contributingPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetContributing()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        contributingPath,
			Description: "Create CONTRIBUTING.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CONTRIBUTING.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(contributingPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        contributingPath,
		Description: "Created CONTRIBUTING.md",
	})

	return result, nil
}

func FixAPIDocsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "api_docs",
		Changes: []Change{},
	}

	apiDocsPath := filepath.Join(ctx.ProjectPath, "docs", "api")

	exists, _ := shared.FileExists(apiDocsPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        apiDocsPath,
			Description: "Create docs/api/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create API docs folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(apiDocsPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create API docs README
	readmePath := filepath.Join(apiDocsPath, "README.md")
	content := resources.GetAPIDocs(ctx.ProjectType)
	shared.CreateFile(readmePath, content)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        apiDocsPath,
		Description: "Created docs/api/ folder",
	})

	return result, nil
}