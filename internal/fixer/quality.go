package fixer

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixEditorconfig(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "editorconfig",
		Path:         ".editorconfig",
		Description:  "Create .editorconfig",
		PromptText:   "Create .editorconfig?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetEditorconfig(ctx.ProjectType), nil
		},
	})
}

func FixPrettier(ctx *FixContext) (*FixResult, error) {
	// Only for Node.js projects
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "prettier", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "prettier",
		Path:         ".prettierrc",
		Description:  "Create .prettierrc",
		PromptText:   "Create Prettier configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetPrettierConfig(), nil
		},
	})
}

func FixPrettierIgnore(ctx *FixContext) (*FixResult, error) {
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "prettierignore", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "prettierignore",
		Path:         ".prettierignore",
		Description:  "Create .prettierignore",
		PromptText:   "Create .prettierignore?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetPrettierIgnore(), nil
		},
	})
}

func FixESLint(ctx *FixContext) (*FixResult, error) {
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "eslint", Skipped: true}, nil
	}

	// Check if TypeScript
	tsconfigPath := filepath.Join(ctx.ProjectPath, "tsconfig.json")
	useTypeScript := false
	if exists, _ := shared.FileExists(tsconfigPath); exists {
		useTypeScript = true
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "eslint",
		Path:         ".eslintrc.json",
		Description:  "Create .eslintrc.json",
		PromptText:   "Create ESLint configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetESLintConfig(useTypeScript), nil
		},
	})
}

func FixCommitlint(ctx *FixContext) (*FixResult, error) {
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "commitlint", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "commitlint",
		Path:         "commitlint.config.js",
		Description:  "Create commitlint.config.js",
		PromptText:   "Create Commitlint configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetCommitlintConfig(), nil
		},
	})
}

func FixHusky(ctx *FixContext) (*FixResult, error) {
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "husky", Skipped: true}, nil
	}

	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "husky",
		Path:        ".husky",
		Description: "Create .husky/ folder with git hooks",
		PromptText:  "Create Husky git hooks?",
		Files: []FolderFileSpec{
			{
				Name: "pre-commit",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetHuskyPreCommit(), nil
				},
				PostCreate: MakeExecutable,
			},
			{
				Name: "commit-msg",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetHuskyCommitMsg(), nil
				},
				PostCreate: MakeExecutable,
			},
		},
	})
}

func FixLintStaged(ctx *FixContext) (*FixResult, error) {
	if ctx.ProjectType != "nodejs" {
		return &FixResult{RuleID: "lint_staged", Skipped: true}, nil
	}

	content := resources.GetLintStagedConfig(ctx.ProjectType)
	if content == "" {
		return &FixResult{RuleID: "lint_staged", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "lint_staged",
		Path:         ".lintstagedrc.json",
		Description:  "Create .lintstagedrc.json",
		PromptText:   "Create lint-staged configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return content, nil
		},
	})
}

func FixMakefile(ctx *FixContext) (*FixResult, error) {
	content := resources.GetMakefile(ctx.ProjectType)
	if content == "" {
		return &FixResult{RuleID: "makefile", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "makefile",
		Path:         "Makefile",
		Description:  "Create Makefile",
		PromptText:   "Create Makefile?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return content, nil
		},
	})
}

func FixGitattributes(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "gitattributes",
		Path:         ".gitattributes",
		Description:  "Create .gitattributes",
		PromptText:   "Create .gitattributes?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetGitattributes(ctx.ProjectType), nil
		},
	})
}

func FixPreCommit(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "pre_commit",
		Path:         ".pre-commit-config.yaml",
		Description:  "Create .pre-commit-config.yaml",
		PromptText:   "Create pre-commit config?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetPreCommitConfig(ctx.ProjectType), nil
		},
	})
}

func FixCodeOwners(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "code_owners",
		Path:         ".github/CODEOWNERS",
		Description:  "Create CODEOWNERS",
		PromptText:   "Create CODEOWNERS?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetCodeowners(ctx.ProjectInfo), nil
		},
	})
}
