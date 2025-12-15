package fixer

import (
	"os"
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

	content := resources.GetEditorconfig(ctx.ProjectType)

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

func FixPrettier(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "prettier",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	prettierPath := filepath.Join(ctx.ProjectPath, ".prettierrc")

	exists, _ := shared.FileExists(prettierPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetPrettierConfig()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        prettierPath,
			Description: "Create .prettierrc",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Prettier configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(prettierPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        prettierPath,
		Description: "Created .prettierrc",
	})

	return result, nil
}

func FixPrettierIgnore(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "prettierignore",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	prettierignorePath := filepath.Join(ctx.ProjectPath, ".prettierignore")

	exists, _ := shared.FileExists(prettierignorePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetPrettierIgnore()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        prettierignorePath,
			Description: "Create .prettierignore",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .prettierignore?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(prettierignorePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        prettierignorePath,
		Description: "Created .prettierignore",
	})

	return result, nil
}

func FixESLint(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "eslint",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	eslintPath := filepath.Join(ctx.ProjectPath, ".eslintrc.json")

	exists, _ := shared.FileExists(eslintPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Check if TypeScript
	useTypeScript := false
	tsconfigPath := filepath.Join(ctx.ProjectPath, "tsconfig.json")
	if exists, _ := shared.FileExists(tsconfigPath); exists {
		useTypeScript = true
	}

	content := resources.GetESLintConfig(useTypeScript)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        eslintPath,
			Description: "Create .eslintrc.json",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create ESLint configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(eslintPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        eslintPath,
		Description: "Created .eslintrc.json",
	})

	return result, nil
}

func FixCommitlint(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "commitlint",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	commitlintPath := filepath.Join(ctx.ProjectPath, "commitlint.config.js")

	exists, _ := shared.FileExists(commitlintPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetCommitlintConfig()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        commitlintPath,
			Description: "Create commitlint.config.js",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Commitlint configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(commitlintPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        commitlintPath,
		Description: "Created commitlint.config.js",
	})

	return result, nil
}

func FixHusky(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "husky",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	huskyPath := filepath.Join(ctx.ProjectPath, ".husky")

	exists, _ := shared.FileExists(huskyPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        huskyPath,
			Description: "Create .husky/ folder with git hooks",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Husky git hooks?") {
			result.Skipped = true
			return result, nil
		}
	}

	// Create .husky folder
	if err := shared.CreateDir(huskyPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create pre-commit hook
	preCommitPath := filepath.Join(huskyPath, "pre-commit")
	preCommitContent := resources.GetHuskyPreCommit()
	if err := shared.CreateFile(preCommitPath, preCommitContent); err == nil {
		// Make executable
		os.Chmod(preCommitPath, 0755)
	}

	// Create commit-msg hook
	commitMsgPath := filepath.Join(huskyPath, "commit-msg")
	commitMsgContent := resources.GetHuskyCommitMsg()
	if err := shared.CreateFile(commitMsgPath, commitMsgContent); err == nil {
		// Make executable
		os.Chmod(commitMsgPath, 0755)
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        huskyPath,
		Description: "Created .husky/ folder with pre-commit and commit-msg hooks",
	})

	return result, nil
}

func FixLintStaged(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "lint_staged",
		Changes: []Change{},
	}

	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		result.Skipped = true
		return result, nil
	}

	lintStagedPath := filepath.Join(ctx.ProjectPath, ".lintstagedrc.json")

	exists, _ := shared.FileExists(lintStagedPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetLintStagedConfig(ctx.ProjectType)
	if content == "" {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        lintStagedPath,
			Description: "Create .lintstagedrc.json",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create lint-staged configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(lintStagedPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        lintStagedPath,
		Description: "Created .lintstagedrc.json",
	})

	return result, nil
}

func FixMakefile(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "makefile",
		Changes: []Change{},
	}

	makefilePath := filepath.Join(ctx.ProjectPath, "Makefile")

	exists, _ := shared.FileExists(makefilePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetMakefile(ctx.ProjectType)
	if content == "" {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        makefilePath,
			Description: "Create Makefile",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Makefile?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(makefilePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        makefilePath,
		Description: "Created Makefile",
	})

	return result, nil
}

func FixGitattributes(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "gitattributes",
		Changes: []Change{},
	}

	gitattributesPath := filepath.Join(ctx.ProjectPath, ".gitattributes")

	exists, _ := shared.FileExists(gitattributesPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetGitattributes(ctx.ProjectType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        gitattributesPath,
			Description: "Create .gitattributes",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .gitattributes?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(gitattributesPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        gitattributesPath,
		Description: "Created .gitattributes",
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

	content := resources.GetPreCommitConfig(ctx.ProjectType)

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

	content := resources.GetCodeowners(ctx.ProjectInfo)

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