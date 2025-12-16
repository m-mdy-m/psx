package fixer

import (
	"github.com/m-mdy-m/psx/internal/resources"
)

func FixReadme(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "readme",
		Path:         "README.md",
		Description:  "Create README.md",
		PromptText:   "Create README.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetReadme(ctx.ProjectInfo, ctx.ProjectType), nil
		},
	})
}

func FixLicense(ctx *FixContext) (*FixResult, error) {
	// Special handling for interactive license selection
	licenseType := ctx.ProjectInfo.License
	if licenseType == "" {
		licenseType = "MIT"
	}

	if ctx.Interactive && !ctx.DryRun {
		licenses := resources.ListLicenses()
		licenses = append(licenses, "Skip")
		choice, _ := ctx.ProjectInfo.ToVars()["license"]
		if choice == "" {
			licenseType = "MIT"
		}
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "license",
		Path:         "LICENSE",
		Description:  "Create LICENSE (" + licenseType + ")",
		PromptText:   "Create LICENSE?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetLicense(licenseType, ctx.ProjectInfo.Author), nil
		},
	})
}

func FixGitignore(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "gitignore",
		Path:         ".gitignore",
		Description:  "Create .gitignore",
		PromptText:   "Create .gitignore?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetGitignore(ctx.ProjectType), nil
		},
	})
}

func FixChangelog(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "changelog",
		Path:         "CHANGELOG.md",
		Description:  "Create CHANGELOG.md",
		PromptText:   "Create CHANGELOG.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetChangelog(ctx.ProjectInfo), nil
		},
	})
}
