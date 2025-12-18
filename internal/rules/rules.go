package rules

import (
	"github.com/m-mdy-m/psx/internal/config"
)

var all = []Rule{
	// GENERAL RULES
	{
		ID:       "readme",
		Type:     RuleTypeFile,
		Category: "general",
		Severity: config.SeverityError,
		Patterns: map[string][]string{
			"*": {"README.md", "README", "readme", "readme.md"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "readme",
			Prompt:       "Create README.md?",
		},
	},
	{
		ID:       "license",
		Type:     RuleTypeFile,
		Category: "general",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"LICENSE", "LICENSE.md"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "license",
			Prompt:       "Create LICENSE?",
		},
	},
	{
		ID:       ".gitignore",
		Type:     RuleTypeFile,
		Category: "general",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".gitignore"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "gitignore",
			Prompt:       "Create .gitignore?",
		},
	},
	{
		ID:       "changelog",
		Type:     RuleTypeFile,
		Category: "general",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"CHANGELOG.md", "CHANGELOG"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "changelog",
			Prompt:       "Create CHANGELOG.md?",
		},
	},
	// DOCS RULES
	{
		ID:       "adr",
		Type:     RuleTypeFolder,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"adr/*"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "adr",
			Prompt:       "Create ADR?",
		},
	},
	{
		ID:       "contributing",
		Type:     RuleTypeFile,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"CONTRIBUTING.md", "CONTRIBUTING"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "contributing",
			Prompt:       "Create CONTRIBUTING.md?",
		},
	},
	{
		ID:       "api_docs",
		Type:     RuleTypeMulti,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				// possible file names
				"API.md",
				"api.md",
				"docs/API.md",
				"docs/api.md",
				"docs/api.md",
				// possible directories
				"api",
				"docs/api",
				"apidocs",
				"docs/apidocs",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeMulti,
			TemplateName: "api_docs",
			Prompt:       "Create API docs (API.md file or api/ directory)?",
		},
	},
	{
		ID:       "security",
		Type:     RuleTypeFile,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"SECURITY.md", "SECURITY"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "security",
			Prompt:       "Create SECURITY.md?",
		},
	},
	{
		ID:       "code_of_conduct",
		Type:     RuleTypeFile,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"CODE_OF_CONDUCT.md", "CODE_OF_CONDUCT"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "code_of_conduct",
			Prompt:       "Create CODE_OF_CONDUCT.md?",
		},
	},
	{
		ID:       "pull_request_template",
		Type:     RuleTypeFile,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				".github/PULL_REQUEST_TEMPLATE.md",
				".github/pull_request_template.md",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "pull_request_template",
			Prompt:       "Create Pull Request template?",
		},
	},
	{
		ID:       "issue_templates",
		Type:     RuleTypeFolder,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				".github/ISSUE_TEMPLATE",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFolder,
			TemplateName: "issue_templates",
			Prompt:       "Create issue templates?",
		},
	},
	{
		ID:       "funding",
		Type:     RuleTypeMulti,
		Category: "docs",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				".github/FUNDING.yml",
				".github/FUNDING.yaml",

				"FUNDING.md",
				"FUNDING",
				"docs/FUNDING.md",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeMulti,
			TemplateName: "funding",
			Prompt:       "Add funding information (GitHub Sponsors or FUNDING.md)?",
		},
	},

	// CI CD RULES
	{
		ID:       "ci_config",
		Type:     RuleTypeMulti,
		Category: "ci",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				// GitHub Actions
				".github/workflows",
				".github/workflows/ci.yml",
				".github/workflows/ci.yaml",

				// GitLab CI
				".gitlab-ci.yml",

				// CircleCI
				".circleci",
				".circleci/config.yml",

				// Travis CI
				".travis.yml",

				// Azure Pipelines
				"azure-pipelines.yml",

				// Bitbucket Pipelines
				"bitbucket-pipelines.yml",

				// Generic CI folders
				"ci",
				".ci",
				"cicd",
				".cicd",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeMulti,
			TemplateName: "ci_config",
			Prompt:       "Create CI/CD configuration (GitHub Actions, GitLab CI, etc.)?",
		},
	},
	{
		ID:       "github_actions",
		Type:     RuleTypeFolder,
		Category: "ci",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".github/workflows"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFolder,
			TemplateName: "github_actions",
			Prompt:       "Create GitHub Actions workflows?",
		},
	},

	// Quality Rules
	{
		ID:       "pre_commit",
		Type:     RuleTypeFile,
		Category: "quality",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".pre-commit-config.yaml"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "pre_commit",
			Prompt:       "Create .pre-commit-config.yaml?",
		},
	},
	{
		ID:       "editorconfig",
		Type:     RuleTypeFile,
		Category: "quality",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".editorconfig"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "editorconfig",
			Prompt:       "Create .editorconfig?",
		},
	},
	{
		ID:       "code_owners",
		Type:     RuleTypeFile,
		Category: "quality",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".github/CODEOWNERS"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "code_owners",
			Prompt:       "Create .github/CODEOWNERS?",
		},
	},

	// Structure

	{
		ID:       "src_folder",
		Type:     RuleTypeFolder,
		Category: "structure",
		Severity: config.SeverityError,
		Patterns: map[string][]string{
			"go": {"cmd/", "pkg/", "internal/", "src/"},
			"*":  {"src/"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFolder,
			TemplateName: "src_folder",
			Prompt:       "Create src/ folder?",
		},
	},
	{
		ID:       "tests_folder",
		Type:     RuleTypeMulti,
		Category: "structure",
		Severity: config.SeverityError,
		Patterns: map[string][]string{
			"nodejs": {"test/", "tests/", "__tests__/"},
			"go":     {"*_test.go"},
			"*":      {"tests/"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeMulti,
			TemplateName: "tests_folder",
			Prompt:       "Create tests/ folder?",
		},
	},
	{
		ID:       "docs_folder",
		Type:     RuleTypeFolder,
		Category: "structure",
		Severity: config.SeverityError,
		Patterns: map[string][]string{
			"*": {"docs/"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFolder,
			TemplateName: "docs_folder",
			Prompt:       "Create docs/ folder?",
		},
	},
	{
		ID:       "scripts_folder",
		Type:     RuleTypeFolder,
		Category: "structure",
		Severity: config.SeverityError,
		Patterns: map[string][]string{
			"*": {"scripts/"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFolder,
			TemplateName: "scripts_folder",
			Prompt:       "Create scripts/ folder?",
		},
	},

	/// DevOps
	{
		ID:       "dockerfile",
		Type:     RuleTypeFile,
		Category: "devops",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {"Dockerfile"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "dockerfile",
			Prompt:       "Create Dockerfile?",
		},
	},
	{
		ID:       "dockerignore",
		Type:     RuleTypeFile,
		Category: "devops",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {".dockerignore"},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "dockerignore",
			Prompt:       "Create dockerignore?",
		},
	},
	{
		ID:       "docker_compose",
		Type:     RuleTypeFile,
		Category: "devops",
		Severity: config.SeverityWarning,
		Patterns: map[string][]string{
			"*": {
				"docker-compose.yml",
				"docker-compose.yaml",
			},
		},
		FixSpec: FixSpec{
			Type:         FixTypeFile,
			TemplateName: "docker_compose",
			Prompt:       "Create docker_compose?",
		},
	},
}

func AllRules() []Rule {
	out := make([]Rule, len(all))
	copy(out, all)
	return out
}
