package fixer

import (
	"fmt"

	"github.com/m-mdy-m/psx/internal/resources"
)

func FixSrcFolder(ctx *FixContext) (*FixResult, error) {
	folderName := resources.GetSrcFolderName(ctx.ProjectType)

	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "src_folder",
		Path:        folderName,
		Description: fmt.Sprintf("Create %s/ folder", folderName),
		PromptText:  fmt.Sprintf("Create %s/ folder?", folderName),
		Files: []FolderFileSpec{
			{
				Name: ".gitkeep",
				GetContent: func(ctx *FixContext) (string, error) {
					return "", nil
				},
			},
		},
	})
}

func FixTestsFolder(ctx *FixContext) (*FixResult, error) {
	folderName := resources.GetTestsFolderName(ctx.ProjectType)

	// Get example test file name and content
	testFileName := resources.GetTestFileName(ctx.ProjectType)
	files := []FolderFileSpec{}

	if testFileName != "" {
		files = append(files, FolderFileSpec{
			Name: testFileName,
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetTestExample(ctx.ProjectType), nil
			},
			FormatForDry: true,
		})
	} else {
		files = append(files, FolderFileSpec{
			Name: ".gitkeep",
			GetContent: func(ctx *FixContext) (string, error) {
				return "", nil
			},
		})
	}

	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "tests_folder",
		Path:        folderName,
		Description: fmt.Sprintf("Create %s/ folder", folderName),
		PromptText:  fmt.Sprintf("Create %s/ folder?", folderName),
		Files:       files,
	})
}

func FixDocsFolder(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "docs_folder",
		Path:        "docs",
		Description: "Create docs/ folder",
		PromptText:  "Create docs/ folder?",
		Files: []FolderFileSpec{
			{
				Name: "README.md",
				GetContent: func(ctx *FixContext) (string, error) {
					return "# Documentation\n\nProject documentation goes here.\n", nil
				},
			},
		},
	})
}

func FixScriptsFolder(ctx *FixContext) (*FixResult, error) {
	files := []FolderFileSpec{
		{
			Name: "install.sh",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetInstallScript(ctx.ProjectInfo, ctx.ProjectType), nil
			},
			PostCreate: MakeExecutable,
		},
		{
			Name: "setup.sh",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetSetupScript(ctx.ProjectType), nil
			},
			PostCreate: MakeExecutable,
		},
		{
			Name: "test.sh",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetTestScript(ctx.ProjectInfo, ctx.ProjectType), nil
			},
			PostCreate: MakeExecutable,
		},
		{
			Name: "build.sh",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetBuildScript(ctx.ProjectInfo, ctx.ProjectType), nil
			},
			PostCreate: MakeExecutable,
		},
		{
			Name: "clean.sh",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetCleanScript(ctx.ProjectInfo, ctx.ProjectType), nil
			},
			PostCreate: MakeExecutable,
		},
	}

	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "scripts_folder",
		Path:        "scripts",
		Description: "Create scripts/ folder with utility scripts",
		PromptText:  "Create scripts folder?",
		Files:       files,
	})
}

func FixEnvExample(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "env_example",
		Path:         ".env.example",
		Description:  "Create .env.example",
		PromptText:   "Create .env.example?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetEnvExample(ctx.ProjectType), nil
		},
	})
}

func FixGitHubActions(ctx *FixContext) (*FixResult, error) {
	// Get CI content
	ciContent := resources.GetGitHubActionsCI(ctx.ProjectInfo, ctx.ProjectType)
	if ciContent == "" {
		return &FixResult{RuleID: "github_actions", Skipped: true}, nil
	}

	files := []FolderFileSpec{
		{
			Name: "ci.yml",
			GetContent: func(ctx *FixContext) (string, error) {
				return ciContent, nil
			},
			FormatForDry: true,
		},
	}

	// Ask about Docker workflow if interactive
	createDocker := false
	if ctx.Interactive && !ctx.DryRun {
		// Would ask user here, but for simplicity we skip
	}

	if createDocker {
		files = append(files, FolderFileSpec{
			Name: "docker.yml",
			GetContent: func(ctx *FixContext) (string, error) {
				return resources.GetGitHubActionsDocker(ctx.ProjectInfo), nil
			},
			FormatForDry: true,
		})
	}

	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "github_actions",
		Path:        ".github/workflows",
		Description: "Create .github/workflows/ with CI workflow",
		PromptText:  "Create GitHub Actions workflow?",
		Validator:   ValidateNotEmpty,
		Files:       files,
	})
}

func FixRenovate(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "renovate",
		Path:         "renovate.json",
		Description:  "Create renovate.json",
		PromptText:   "Create Renovate configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetRenovateConfig(), nil
		},
	})
}

func FixDependabot(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "dependabot",
		Path:         ".github/dependabot.yml",
		Description:  "Create .github/dependabot.yml",
		PromptText:   "Create Dependabot configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetDependabotConfig(ctx.ProjectType), nil
		},
	})
}
