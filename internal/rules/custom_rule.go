package rules

import (
	"fmt"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

type CustomHandler struct {
	ctx        *Context
	customCfg  *config.CustomConfig
	projectDir string
}

func NewCustomHandler(ctx *Context, customCfg *config.CustomConfig) *CustomHandler {
	return &CustomHandler{
		ctx:        ctx,
		customCfg:  customCfg,
		projectDir: ctx.ProjectPath,
	}
}

func (h *CustomHandler) ApplyCustomFiles(fixCtx *FixContext) ([]*FixResult, error) {
	if h.customCfg == nil {
		return nil, nil
	}

	results := []*FixResult{}
	if len(h.customCfg.Files) > 0 {
		fileResults := h.applyFiles(fixCtx)
		results = append(results, fileResults...)
	}
	if len(h.customCfg.Folders) > 0 {
		folderResults := h.applyFolders(fixCtx)
		results = append(results, folderResults...)
	}

	return results, nil
}

func (h *CustomHandler) applyFiles(fixCtx *FixContext) []*FixResult {
	results := []*FixResult{}

	for _, customFile := range h.customCfg.Files {
		fullPath := filepath.Join(h.projectDir, customFile.Path)

		// Check if file already exists
		if exists, info := utils.FileExists(fullPath); exists && info.Size() > 0 {
			logger.Verbose(fmt.Sprintf("Skipping custom file (exists): %s", customFile.Path))
			continue
		}

		// Interactive prompt
		if fixCtx.Interactive && !fixCtx.DryRun {
			if !utils.Prompt(fmt.Sprintf("Create custom file %s?", customFile.Path)) {
				continue
			}
		}

		result := &FixResult{
			RuleID: fmt.Sprintf("custom:file:%s", customFile.Path),
		}

		if fixCtx.DryRun {
			// Dry run preview
			result.Fixed = true
			result.Changes = []Change{
				{
					Type:        ChangeCreateFile,
					Path:        fullPath,
					Description: fmt.Sprintf("Create %s", customFile.Path),
					Content:     formatContent(customFile.Content, 10),
				},
			}
		} else {
			// Actually create file
			err := utils.CreateFile(fullPath, customFile.Content)
			if err != nil {
				result.Error = err
				logger.Error(fmt.Sprintf("Failed to create %s: %v", customFile.Path, err))
			} else {
				result.Fixed = true
				result.Changes = []Change{
					{
						Type:        ChangeCreateFile,
						Path:        fullPath,
						Description: fmt.Sprintf("Created %s", customFile.Path),
					},
				}
			}
		}

		results = append(results, result)
	}

	return results
}

func (h *CustomHandler) applyFolders(fixCtx *FixContext) []*FixResult {
	results := []*FixResult{}

	for _, customFolder := range h.customCfg.Folders {
		fullPath := filepath.Join(h.projectDir, customFolder.Path)

		// Interactive prompt
		if fixCtx.Interactive && !fixCtx.DryRun {
			prompt := fmt.Sprintf("Create custom folder structure %s?", customFolder.Path)
			if !utils.Prompt(prompt) {
				continue
			}
		}

		result := &FixResult{
			RuleID: fmt.Sprintf("custom:folder:%s", customFolder.Path),
		}

		if fixCtx.DryRun {
			// Dry run preview
			changes := h.previewFolderStructure(customFolder, fullPath)
			result.Fixed = true
			result.Changes = changes
		} else {
			// Actually create structure
			changes, err := h.createFolderStructure(customFolder, fullPath)
			if err != nil {
				result.Error = err
				logger.Error(fmt.Sprintf("Failed to create folder structure %s: %v", customFolder.Path, err))
			} else {
				result.Fixed = true
				result.Changes = changes
			}
		}

		results = append(results, result)
	}

	return results
}

func (h *CustomHandler) previewFolderStructure(folder config.CustomFolder, basePath string) []Change {
	changes := []Change{}

	// Root folder
	changes = append(changes, Change{
		Type:        ChangeCreateFolder,
		Path:        basePath,
		Description: fmt.Sprintf("Create %s/", folder.Path),
	})

	// Sub-structure
	if folder.Structure != nil {
		subChanges := h.previewStructure(folder.Structure, basePath, folder.Path)
		changes = append(changes, subChanges...)
	}

	return changes
}

func (h *CustomHandler) previewStructure(structure map[string]interface{}, basePath, relativePath string) []Change {
	changes := []Change{}

	for name, value := range structure {
		itemPath := filepath.Join(basePath, name)
		relItemPath := filepath.Join(relativePath, name)

		if subStructure, ok := value.(map[string]interface{}); ok {
			// It's a folder
			changes = append(changes, Change{
				Type:        ChangeCreateFolder,
				Path:        itemPath,
				Description: fmt.Sprintf("Create %s/", relItemPath),
			})

			// Recursively preview sub-structure
			if len(subStructure) > 0 {
				subChanges := h.previewStructure(subStructure, itemPath, relItemPath)
				changes = append(changes, subChanges...)
			}
		}
	}

	return changes
}

func (h *CustomHandler) createFolderStructure(folder config.CustomFolder, basePath string) ([]Change, error) {
	changes := []Change{}

	// Create root folder
	err := utils.CreateDir(basePath)
	if err != nil {
		return nil, err
	}

	changes = append(changes, Change{
		Type:        ChangeCreateFolder,
		Path:        basePath,
		Description: fmt.Sprintf("Created %s/", folder.Path),
	})

	// Create sub-structure
	if folder.Structure != nil {
		subChanges, err := h.createStructure(folder.Structure, basePath, folder.Path)
		if err != nil {
			return changes, err
		}
		changes = append(changes, subChanges...)
	}

	return changes, nil
}

func (h *CustomHandler) createStructure(structure map[string]interface{}, basePath, relativePath string) ([]Change, error) {
	changes := []Change{}

	for name, value := range structure {
		itemPath := filepath.Join(basePath, name)
		relItemPath := filepath.Join(relativePath, name)

		if subStructure, ok := value.(map[string]interface{}); ok {
			// It's a folder
			err := utils.CreateDir(itemPath)
			if err != nil {
				return changes, err
			}

			changes = append(changes, Change{
				Type:        ChangeCreateFolder,
				Path:        itemPath,
				Description: fmt.Sprintf("Created %s/", relItemPath),
			})

			// Recursively create sub-structure
			if len(subStructure) > 0 {
				subChanges, err := h.createStructure(subStructure, itemPath, relItemPath)
				if err != nil {
					return changes, err
				}
				changes = append(changes, subChanges...)
			}
		}
	}

	return changes, nil
}
