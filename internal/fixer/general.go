package fixer

import (
	"fmt"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixReadme(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "readme",
		Changes: []Change{},
	}

	readmePath := filepath.Join(ctx.ProjectPath, "README.md")

	exists, _ := shared.FileExists(readmePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Use ProjectInfo to generate README
	content := resources.GetReadme(ctx.ProjectInfo, ctx.ProjectType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        readmePath,
			Description: "Create README.md",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create README.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(readmePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        readmePath,
		Description: "Created README.md",
	})

	return result, nil
}

func FixLicense(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "license",
		Changes: []Change{},
	}

	licensePath := filepath.Join(ctx.ProjectPath, "LICENSE")

	exists, _ := shared.FileExists(licensePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	licenseType := ctx.ProjectInfo.License
	if licenseType == "" {
		licenseType = "MIT"
	}

	if ctx.Interactive {
		licenses := resources.ListLicenses()
		licenses = append(licenses, "Skip")
		choice, _ := shared.PromptChoice("Choose license type:", licenses)
		if choice == "Skip" {
			result.Skipped = true
			return result, nil
		}
		licenseType = choice
	}

	content := resources.GetLicense(licenseType, ctx.ProjectInfo.Author)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        licensePath,
			Description: fmt.Sprintf("Create LICENSE (%s)", licenseType),
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if err := shared.CreateFile(licensePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        licensePath,
		Description: fmt.Sprintf("Created LICENSE (%s)", licenseType),
	})

	return result, nil
}

func FixGitignore(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "gitignore",
		Changes: []Change{},
	}

	gitignorePath := filepath.Join(ctx.ProjectPath, ".gitignore")

	exists, _ := shared.FileExists(gitignorePath)

	content := resources.GetGitignore(ctx.ProjectType)

	if exists {
		// TODO: Append missing patterns
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        gitignorePath,
			Description: "Create .gitignore",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .gitignore?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(gitignorePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        gitignorePath,
		Description: "Created .gitignore",
	})

	return result, nil
}

func FixChangelog(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "changelog",
		Changes: []Change{},
	}

	changelogPath := filepath.Join(ctx.ProjectPath, "CHANGELOG.md")

	exists, _ := shared.FileExists(changelogPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetChangelog(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        changelogPath,
			Description: "Create CHANGELOG.md",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CHANGELOG.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(changelogPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        changelogPath,
		Description: "Created CHANGELOG.md",
	})

	return result, nil
}

func FormatContent(content string, maxLines int) string {
	lines := []string{}
	current := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	if len(lines) <= maxLines {
		return content
	}

	displayed := ""
	for i := 0; i < maxLines; i++ {
		displayed += lines[i] + "\n"
	}
	remaining := len(lines) - maxLines
	return fmt.Sprintf("%s... (%d more lines)", displayed, remaining)
}