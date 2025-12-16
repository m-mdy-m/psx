package fixer

import (
	"github.com/m-mdy-m/psx/internal/resources"
)

func FixADR(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "adr",
		Path:        "docs/adr",
		Description: "Create docs/adr/ folder with initial ADR",
		PromptText:  "Create ADR folder?",
		Files: []FolderFileSpec{
			{
				Name: "0001-record-architecture-decisions.md",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetFirstADR(ctx.ProjectInfo), nil
				},
				FormatForDry: true,
			},
			{
				Name: "template.md",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetADRTemplate(), nil
				},
				FormatForDry: true,
			},
		},
	})
}

func FixContributing(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "contributing",
		Path:         "CONTRIBUTING.md",
		Description:  "Create CONTRIBUTING.md",
		PromptText:   "Create CONTRIBUTING.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetContributing(), nil
		},
	})
}

func FixAPIDocsFolder(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "api_docs",
		Path:        "docs/api",
		Description: "Create docs/api/ folder",
		PromptText:  "Create API docs folder?",
		Files: []FolderFileSpec{
			{
				Name: "README.md",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetAPIDocs(ctx.ProjectInfo, ctx.ProjectType), nil
				},
				FormatForDry: true,
			},
		},
	})
}

func FixSecurity(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "security",
		Path:         "SECURITY.md",
		Description:  "Create SECURITY.md",
		PromptText:   "Create SECURITY.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetSecurity(ctx.ProjectInfo), nil
		},
	})
}

func FixCodeOfConduct(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "code_of_conduct",
		Path:         "CODE_OF_CONDUCT.md",
		Description:  "Create CODE_OF_CONDUCT.md",
		PromptText:   "Create CODE_OF_CONDUCT.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetCodeOfConduct(ctx.ProjectInfo), nil
		},
	})
}

func FixPullRequestTemplate(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "pull_request_template",
		Path:         ".github/PULL_REQUEST_TEMPLATE.md",
		Description:  "Create .github/PULL_REQUEST_TEMPLATE.md",
		PromptText:   "Create pull request template?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetPullRequestTemplate(), nil
		},
	})
}

func FixIssueTemplates(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "issue_templates",
		Path:        ".github/ISSUE_TEMPLATE",
		Description: "Create .github/ISSUE_TEMPLATE/ with templates",
		PromptText:  "Create issue templates?",
		Files: []FolderFileSpec{
			{
				Name: "bug_report.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetIssueBugReport(), nil
				},
				FormatForDry: true,
			},
			{
				Name: "feature_request.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetIssueFeatureRequest(), nil
				},
				FormatForDry: true,
			},
			{
				Name: "config.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetIssueTemplatesConfig(ctx.ProjectInfo), nil
				},
				FormatForDry: true,
			},
		},
	})
}

func FixFunding(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "funding",
		Path:         ".github/FUNDING.yml",
		Description:  "Create .github/FUNDING.yml",
		PromptText:   "Create funding configuration?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetFundingYML(ctx.ProjectInfo), nil
		},
	})
}

func FixSupport(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "support",
		Path:         "SUPPORT.md",
		Description:  "Create SUPPORT.md",
		PromptText:   "Create SUPPORT.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetSupport(ctx.ProjectInfo), nil
		},
	})
}

func FixRoadmap(ctx *FixContext) (*FixResult, error) {
	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "roadmap",
		Path:         "ROADMAP.md",
		Description:  "Create ROADMAP.md",
		PromptText:   "Create ROADMAP.md?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetRoadmap(ctx.ProjectInfo), nil
		},
	})
}
