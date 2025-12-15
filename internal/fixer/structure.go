package fixer

import (
	"fmt"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixSrcFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "src_folder",
		Changes: []Change{},
	}

	folderName := resources.GetSrcFolderName(ctx.ProjectType)
	folderPath := filepath.Join(ctx.ProjectPath, folderName)

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: fmt.Sprintf("Create %s/ folder", folderName),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt(fmt.Sprintf("Create %s/ folder?", folderName)) {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create .gitkeep
	gitkeepPath := filepath.Join(folderPath, ".gitkeep")
	shared.CreateFile(gitkeepPath, "")

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: fmt.Sprintf("Created %s/ folder", folderName),
	})

	return result, nil
}

func FixTestsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "tests_folder",
		Changes: []Change{},
	}

	folderName := resources.GetTestsFolderName(ctx.ProjectType)
	folderPath := filepath.Join(ctx.ProjectPath, folderName)

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: fmt.Sprintf("Create %s/ folder", folderName),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt(fmt.Sprintf("Create %s/ folder?", folderName)) {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create example test file
	exampleFile := resources.GetTestFileName(ctx.ProjectType)
	if exampleFile != "" {
		testPath := filepath.Join(folderPath, exampleFile)
		content := resources.GetTestExample(ctx.ProjectType)
		shared.CreateFile(testPath, content)
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: fmt.Sprintf("Created %s/ folder with example test", folderName),
	})

	return result, nil
}

func FixDocsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "docs_folder",
		Changes: []Change{},
	}

	folderPath := filepath.Join(ctx.ProjectPath, "docs")

	exists, _ := shared.FileExists(folderPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: "Create docs/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create docs/ folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create docs README
	readmePath := filepath.Join(folderPath, "README.md")
	content := "# Documentation\n\nProject documentation goes here.\n"
	shared.CreateFile(readmePath, content)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: "Created docs/ folder",
	})

	return result, nil
}