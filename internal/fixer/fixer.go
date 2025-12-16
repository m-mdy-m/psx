package fixer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/shared"
)

func pathExists(fullPath string) (bool, os.FileInfo) {
	return shared.FileExists(fullPath)
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

func askUser(ctx *FixContext, prompt string) bool {
	if !ctx.Interactive {
		return true
	}
	return shared.Prompt(prompt)
}

func formatContent(content string, doFormat bool) string {
	if !doFormat {
		return content
	}
	return FormatContent(content, 10)
}

func FixSingleFile(ctx *FixContext, spec FileFixSpec) (*FixResult, error) {
	result := &FixResult{
		RuleID:  spec.RuleID,
		Changes: []Change{},
	}

	fullPath := filepath.Join(ctx.ProjectPath, spec.Path)

	// Custom validator
	if spec.Validator != nil {
		skip, err := spec.Validator(ctx, fullPath)
		if err != nil {
			result.Error = err
			return result, err
		}
		if skip {
			result.Skipped = true
			return result, nil
		}
	} else {
		// Default: check if exists
		if exists, _ := pathExists(fullPath); exists {
			result.Skipped = true
			return result, nil
		}
	}

	// Dry run mode
	if ctx.DryRun {
		content := ""
		if spec.GetContent != nil {
			c, err := spec.GetContent(ctx)
			if err == nil {
				content = formatContent(c, spec.FormatForDry)
			}
		}

		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        fullPath,
			Description: spec.Description,
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	// Ask user if interactive
	if !askUser(ctx, spec.PromptText) {
		result.Skipped = true
		return result, nil
	}

	// Generate content
	content := ""
	if spec.GetContent != nil {
		c, err := spec.GetContent(ctx)
		if err != nil {
			result.Error = err
			return result, err
		}
		content = c
	}

	// Create file
	if err := shared.CreateFile(fullPath, content); err != nil {
		result.Error = err
		return result, err
	}

	// Post-create hook
	if spec.PostCreate != nil {
		if err := spec.PostCreate(ctx, fullPath); err != nil {
			result.Error = err
			return result, err
		}
	}

	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        fullPath,
		Description: spec.Description,
	})
	result.Fixed = true

	return result, nil
}

func FixFolder(ctx *FixContext, spec FolderFixSpec) (*FixResult, error) {
	result := &FixResult{
		RuleID:  spec.RuleID,
		Changes: []Change{},
	}

	folderPath := filepath.Join(ctx.ProjectPath, spec.Path)

	// Custom validator
	if spec.Validator != nil {
		skip, err := spec.Validator(ctx, folderPath)
		if err != nil {
			result.Error = err
			return result, err
		}
		if skip {
			result.Skipped = true
			return result, nil
		}
	} else {
		// Default: check if exists
		if exists, _ := pathExists(folderPath); exists {
			result.Skipped = true
			return result, nil
		}
	}

	// Dry run mode
	if ctx.DryRun {
		changes := []Change{{
			Type:        ChangeCreateFolder,
			Path:        folderPath,
			Description: spec.Description,
		}}

		for _, f := range spec.Files {
			content := ""
			if f.GetContent != nil && f.FormatForDry {
				if c, err := f.GetContent(ctx); err == nil {
					content = formatContent(c, f.FormatForDry)
				}
			}
			changes = append(changes, Change{
				Type:        ChangeCreateFile,
				Path:        filepath.Join(folderPath, f.Name),
				Description: fmt.Sprintf("Create %s", f.Name),
				Content:     content,
			})
		}

		result.Changes = changes
		result.Fixed = true
		return result, nil
	}

	// Ask user if interactive
	if !askUser(ctx, spec.PromptText) {
		result.Skipped = true
		return result, nil
	}

	// Create folder
	if err := shared.CreateDir(folderPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create files inside folder
	for _, f := range spec.Files {
		filePath := filepath.Join(folderPath, f.Name)

		content := ""
		if f.GetContent != nil {
			c, err := f.GetContent(ctx)
			if err != nil {
				result.Error = err
				return result, err
			}
			content = c
		}

		if err := shared.CreateFile(filePath, content); err != nil {
			result.Error = err
			return result, err
		}

		// Post-create hook
		if f.PostCreate != nil {
			if err := f.PostCreate(ctx, filePath); err != nil {
				result.Error = err
				return result, err
			}
		}
	}

	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        folderPath,
		Description: spec.Description,
	})
	result.Fixed = true

	return result, nil
}

func FixMultipleFiles(ctx *FixContext, spec MultiFileFixSpec) (*FixResult, error) {
	result := &FixResult{
		RuleID:  spec.RuleID,
		Changes: []Change{},
	}

	// Custom validator
	if spec.Validator != nil {
		skip, err := spec.Validator(ctx, "")
		if err != nil {
			result.Error = err
			return result, err
		}
		if skip {
			result.Skipped = true
			return result, nil
		}
	}

	// Dry run mode
	if ctx.DryRun {
		for _, f := range spec.Files {
			fullPath := filepath.Join(ctx.ProjectPath, f.Path)
			content := ""
			if f.GetContent != nil {
				if c, err := f.GetContent(ctx); err == nil {
					content = FormatContent(c, 10)
				}
			}
			result.Changes = append(result.Changes, Change{
				Type:        ChangeCreateFile,
				Path:        fullPath,
				Description: f.Description,
				Content:     content,
			})
		}
		result.Fixed = true
		return result, nil
	}

	// Ask user if interactive
	if !askUser(ctx, spec.PromptText) {
		result.Skipped = true
		return result, nil
	}

	// Create all files
	for _, f := range spec.Files {
		fullPath := filepath.Join(ctx.ProjectPath, f.Path)

		// Skip if exists
		if exists, _ := pathExists(fullPath); exists {
			continue
		}

		content := ""
		if f.GetContent != nil {
			c, err := f.GetContent(ctx)
			if err != nil {
				result.Error = err
				return result, err
			}
			content = c
		}

		if err := shared.CreateFile(fullPath, content); err != nil {
			result.Error = err
			return result, err
		}

		// Post-create hook
		if f.PostCreate != nil {
			if err := f.PostCreate(ctx, fullPath); err != nil {
				result.Error = err
				return result, err
			}
		}

		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        fullPath,
			Description: f.Description,
		})
	}

	result.Fixed = true
	return result, nil
}

func ValidateNotEmpty(ctx *FixContext, fullPath string) (bool, error) {
	exists, info := pathExists(fullPath)
	if !exists {
		return false, nil
	}
	if info != nil && info.IsDir() {
		isEmpty, err := shared.IsDirEmpty(fullPath)
		if err != nil {
			return true, err
		}
		// Skip if not empty
		return !isEmpty, nil
	}
	return true, nil
}

func ValidateFileSize(minSize int64) ValidatorFunc {
	return func(ctx *FixContext, fullPath string) (bool, error) {
		exists, info := pathExists(fullPath)
		if !exists {
			return false, nil
		}
		if info != nil && !info.IsDir() {
			// Skip if file is large enough
			return info.Size() >= minSize, nil
		}
		return true, nil
	}
}

func MakeExecutable(ctx *FixContext, fullPath string) error {
	return os.Chmod(fullPath, 0755)
}
