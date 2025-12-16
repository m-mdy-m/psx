package fixer

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixADR(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "adr",
		Changes: []Change{},
	}

	adrPath := filepath.Join(ctx.ProjectPath, "docs", "adr")

	exists, _ := shared.FileExists(adrPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        adrPath,
			Description: "Create docs/adr/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create ADR folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(adrPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create first ADR
	firstADR := filepath.Join(adrPath, "0001-record-architecture-decisions.md")
	content := resources.GetFirstADR(ctx.ProjectInfo)
	if err := shared.CreateFile(firstADR, content); err != nil {
		result.Error = err
		return result, err
	}

	// Create template
	template := filepath.Join(adrPath, "template.md")
	templateContent := resources.GetADRTemplate()
	shared.CreateFile(template, templateContent)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        adrPath,
		Description: "Created docs/adr/ folder with initial ADR",
	})

	return result, nil
}

func FixContributing(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "contributing",
		Changes: []Change{},
	}

	contributingPath := filepath.Join(ctx.ProjectPath, "CONTRIBUTING.md")

	exists, _ := shared.FileExists(contributingPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetContributing()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        contributingPath,
			Description: "Create CONTRIBUTING.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CONTRIBUTING.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(contributingPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        contributingPath,
		Description: "Created CONTRIBUTING.md",
	})

	return result, nil
}

func FixAPIDocsFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "api_docs",
		Changes: []Change{},
	}

	apiDocsPath := filepath.Join(ctx.ProjectPath, "docs", "api")

	exists, _ := shared.FileExists(apiDocsPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        apiDocsPath,
			Description: "Create docs/api/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create API docs folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(apiDocsPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create API docs README
	readmePath := filepath.Join(apiDocsPath, "README.md")
	content := resources.GetAPIDocs(ctx.ProjectInfo, ctx.ProjectType)
	shared.CreateFile(readmePath, content)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        apiDocsPath,
		Description: "Created docs/api/ folder",
	})

	return result, nil
}

func FixSecurity(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "security",
		Changes: []Change{},
	}

	securityPath := filepath.Join(ctx.ProjectPath, "SECURITY.md")

	exists, _ := shared.FileExists(securityPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetSecurity(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        securityPath,
			Description: "Create SECURITY.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create SECURITY.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(securityPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        securityPath,
		Description: "Created SECURITY.md",
	})

	return result, nil
}

func FixCodeOfConduct(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "code_of_conduct",
		Changes: []Change{},
	}

	cocPath := filepath.Join(ctx.ProjectPath, "CODE_OF_CONDUCT.md")

	exists, _ := shared.FileExists(cocPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetCodeOfConduct(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        cocPath,
			Description: "Create CODE_OF_CONDUCT.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CODE_OF_CONDUCT.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(cocPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        cocPath,
		Description: "Created CODE_OF_CONDUCT.md",
	})

	return result, nil
}

func FixPullRequestTemplate(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "pull_request_template",
		Changes: []Change{},
	}

	prTemplatePath := filepath.Join(ctx.ProjectPath, ".github", "PULL_REQUEST_TEMPLATE.md")

	exists, _ := shared.FileExists(prTemplatePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetPullRequestTemplate()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        prTemplatePath,
			Description: "Create .github/PULL_REQUEST_TEMPLATE.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create pull request template?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(prTemplatePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        prTemplatePath,
		Description: "Created .github/PULL_REQUEST_TEMPLATE.md",
	})

	return result, nil
}

func FixIssueTemplates(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "issue_templates",
		Changes: []Change{},
	}

	issueTemplatesPath := filepath.Join(ctx.ProjectPath, ".github", "ISSUE_TEMPLATE")

	exists, _ := shared.FileExists(issueTemplatesPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        issueTemplatesPath,
			Description: "Create .github/ISSUE_TEMPLATE/ with bug report, feature request, and question templates",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create issue templates?") {
			result.Skipped = true
			return result, nil
		}
	}

	// Create ISSUE_TEMPLATE folder
	if err := shared.CreateDir(issueTemplatesPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create bug_report.md
	bugReportPath := filepath.Join(issueTemplatesPath, "bug_report.md")
	bugReportContent := resources.GetIssueBugReport()
	shared.CreateFile(bugReportPath, bugReportContent)

	// Create feature_request.md
	featurePath := filepath.Join(issueTemplatesPath, "feature_request.md")
	featureContent := resources.GetIssueFeatureRequest()
	shared.CreateFile(featurePath, featureContent)

	// Create question.md
	questionPath := filepath.Join(issueTemplatesPath, "question.md")
	questionContent := resources.GetIssueQuestion()
	shared.CreateFile(questionPath, questionContent)

	// Create config.yml
	configPath := filepath.Join(issueTemplatesPath, "config.yml")
	configContent := resources.GetIssueTemplatesConfig(ctx.ProjectInfo)
	shared.CreateFile(configPath, configContent)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        issueTemplatesPath,
		Description: "Created .github/ISSUE_TEMPLATE/ with bug report, feature request, question, and config",
	})

	return result, nil
}

func FixFunding(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "funding",
		Changes: []Change{},
	}

	fundingPath := filepath.Join(ctx.ProjectPath, ".github", "FUNDING.yml")

	exists, _ := shared.FileExists(fundingPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetFundingYML(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        fundingPath,
			Description: "Create .github/FUNDING.yml",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create funding configuration?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(fundingPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        fundingPath,
		Description: "Created .github/FUNDING.yml",
	})

	return result, nil
}

func FixSupport(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "support",
		Changes: []Change{},
	}

	supportPath := filepath.Join(ctx.ProjectPath, "SUPPORT.md")

	exists, _ := shared.FileExists(supportPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetSupport(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        supportPath,
			Description: "Create SUPPORT.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create SUPPORT.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(supportPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        supportPath,
		Description: "Created SUPPORT.md",
	})

	return result, nil
}

func FixRoadmap(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "roadmap",
		Changes: []Change{},
	}

	roadmapPath := filepath.Join(ctx.ProjectPath, "ROADMAP.md")

	exists, _ := shared.FileExists(roadmapPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetRoadmap(ctx.ProjectInfo)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        roadmapPath,
			Description: "Create ROADMAP.md",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create ROADMAP.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(roadmapPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        roadmapPath,
		Description: "Created ROADMAP.md",
	})

	return result, nil
}
