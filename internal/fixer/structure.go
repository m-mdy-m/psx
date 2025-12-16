package fixer

import (
	"fmt"
	"os"
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

func FixScriptsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "scripts_folder",
		Changes: []Change{},
	}

	scriptsPath := filepath.Join(ctx.ProjectPath, "scripts")

	exists, _ := shared.FileExists(scriptsPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        scriptsPath,
			Description: "Create scripts/ folder with utility scripts",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create scripts folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	// Create scripts folder
	if err := shared.CreateDir(scriptsPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create common scripts
	// Install script
	installPath := filepath.Join(scriptsPath, "install.sh")
	installContent := resources.GetInstallScript(ctx.ProjectInfo, ctx.ProjectType)
	if err := shared.CreateFile(installPath, installContent); err == nil {
		os.Chmod(installPath, 0755)
	}

	// Setup script
	setupPath := filepath.Join(scriptsPath, "setup.sh")
	setupContent := resources.GetSetupScript(ctx.ProjectType)
	if err := shared.CreateFile(setupPath, setupContent); err == nil {
		os.Chmod(setupPath, 0755)
	}

	// Test script
	testPath := filepath.Join(scriptsPath, "test.sh")
	testContent := resources.GetTestScript(ctx.ProjectInfo, ctx.ProjectType)
	if err := shared.CreateFile(testPath, testContent); err == nil {
		os.Chmod(testPath, 0755)
	}

	// Build script
	buildPath := filepath.Join(scriptsPath, "build.sh")
	buildContent := resources.GetBuildScript(ctx.ProjectInfo, ctx.ProjectType)
	if err := shared.CreateFile(buildPath, buildContent); err == nil {
		os.Chmod(buildPath, 0755)
	}

	// Clean script
	cleanPath := filepath.Join(scriptsPath, "clean.sh")
	cleanContent := resources.GetCleanScript(ctx.ProjectInfo, ctx.ProjectType)
	if err := shared.CreateFile(cleanPath, cleanContent); err == nil {
		os.Chmod(cleanPath, 0755)
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        scriptsPath,
		Description: "Created scripts/ folder with install, setup, test, build, and clean scripts",
	})

	return result, nil
}

func FixEnvExample(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "env_example",
		Changes: []Change{},
	}

	envExamplePath := filepath.Join(ctx.ProjectPath, ".env.example")

	exists, _ := shared.FileExists(envExamplePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetEnvExample(ctx.ProjectType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        envExamplePath,
			Description: "Create .env.example",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .env.example?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(envExamplePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        envExamplePath,
		Description: "Created .env.example",
	})

	return result, nil
}

func FixGitHubActions(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "github_actions",
		Changes: []Change{},
	}

	workflowsPath := filepath.Join(ctx.ProjectPath, ".github", "workflows")

	exists, info := shared.FileExists(workflowsPath)
	if exists && info != nil && info.IsDir() {
		isEmpty, _ := shared.IsDirEmpty(workflowsPath)
		if !isEmpty {
			result.Skipped = true
			return result, nil
		}
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        workflowsPath,
			Description: "Create .github/workflows/ with CI workflow",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create GitHub Actions workflow?") {
			result.Skipped = true
			return result, nil
		}
	}

	// Create workflows folder
	if err := shared.CreateDir(workflowsPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create CI workflow
	ciPath := filepath.Join(workflowsPath, "ci.yml")
	ciContent := resources.GetGitHubActionsCI(ctx.ProjectInfo, ctx.ProjectType)
	if ciContent != "" {
		shared.CreateFile(ciPath, ciContent)
	}

	// Ask if user wants Docker build workflow
	createDocker := false
	if ctx.Interactive {
		createDocker = shared.Prompt("Create Docker build workflow?")
	}

	if createDocker {
		dockerPath := filepath.Join(workflowsPath, "docker.yml")
		dockerContent := resources.GetGitHubActionsDocker(ctx.ProjectInfo)
		shared.CreateFile(dockerPath, dockerContent)
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        workflowsPath,
		Description: "Created .github/workflows/ with CI workflow",
	})

	return result, nil
}

func FixRenovate(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "renovate",
		Changes: []Change{},
	}

	renovatePath := filepath.Join(ctx.ProjectPath, "renovate.json")

	exists, _ := shared.FileExists(renovatePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetRenovateConfig()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        renovatePath,
			Description: "Create renovate.json",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Renovate configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(renovatePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        renovatePath,
		Description: "Created renovate.json",
	})

	return result, nil
}

func FixDependabot(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "dependabot",
		Changes: []Change{},
	}

	dependabotPath := filepath.Join(ctx.ProjectPath, ".github", "dependabot.yml")

	exists, _ := shared.FileExists(dependabotPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetDependabotConfig(ctx.ProjectType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        dependabotPath,
			Description: "Create .github/dependabot.yml",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Dependabot configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(dependabotPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        dependabotPath,
		Description: "Created .github/dependabot.yml",
	})

	return result, nil
}
