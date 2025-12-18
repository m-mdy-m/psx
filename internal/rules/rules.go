package rules

import (
	"github.com/m-mdy-m/psx/internal/config"
)

var all = []Rule{
	// GENERAL RULES
	fileRule("readme", "general", config.SeverityError, map[string][]string{
		"*": {"README.md", "README", "readme", "readme.md"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create README.md?",
		}),
	fileRule("license", "general", config.SeverityWarning, map[string][]string{
		"*": {"LICENSE", "LICENSE.md"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create LICENSE.md?",
		}),
	fileRule(".gitignore", "general", config.SeverityWarning, map[string][]string{
		"*": {".gitignore"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create .gitignore?",
		}),
	fileRule("changelog", "general", config.SeverityWarning, map[string][]string{
		"*": {"CHANGELOG.md", "CHANGELOG"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create CHANGELOG.md?",
		}),
	// Documentation RULES
	folderRule("adr", "documentation", config.SeverityInfo, map[string][]string{
		"*": {"adr/*"},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create ADR?",
	}),
	fileRule("contributing", "documentation", config.SeverityInfo, map[string][]string{
		"*": {"CONTRIBUTING.md", "CONTRIBUTING"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create CONTRIBUTING.md?",
		}),
	multiRule("api_docs", "documentation", config.SeverityInfo,
		map[string][]string{
			"*": {"docs/api/", "API.md", "api/"},
		}, FixSpec{
			Type:   FixTypeMulti,
			Prompt: "Create API docs?",
		}),

	fileRule("security", "documentation", config.SeverityInfo, map[string][]string{
		"*": {"SECURITY.md", "SECURITY"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create SECURITY.md?",
		}),
	fileRule("code_of_conduct", "documentation", config.SeverityInfo, map[string][]string{
		"*": {"CODE_OF_CONDUCT.md", "CODE_OF_CONDUCT"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create CODE_OF_CONDUCT.md?",
		}),
	fileRule("pull_request_template", "documentation", config.SeverityInfo, map[string][]string{
		"*": {
			".github/PULL_REQUEST_TEMPLATE.md",
			".github/pull_request_template.md",
		},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create Pull Request template?",
		}),

	folderRule("issue_templates", "documentation", config.SeverityInfo, map[string][]string{
		"*": {
			".github/ISSUE_TEMPLATE",
		},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create issue templates?",
	}),

	// CI CD RULES
	multiRule("ci_config", "cicd", config.SeverityInfo,
		map[string][]string{
			"*": {
				".github/workflows/", ".gitlab-ci.yml", ".circleci/config.yml",
			}}, FixSpec{
			Type:   FixTypeMulti,
			Prompt: "Create CI/CD configuration?",
		}),

	folderRule("github_actions", "cicd", config.SeverityInfo, map[string][]string{
		"*": {".github/workflows"},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create GitHub Actions workflows?",
	}),

	// Quality Rules
	fileRule("pre_commit", "quality", config.SeverityInfo, map[string][]string{
		"*": {".pre-commit-config.yaml"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create .pre-commit-config.yaml?",
		}),
	fileRule("editorconfig", "quality", config.SeverityInfo, map[string][]string{
		"*": {".editorconfig"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create .editorconfig?",
		}),
	fileRule("code_owners", "quality", config.SeverityInfo, map[string][]string{
		"*": {".github/CODEOWNERS"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create .github/CODEOWNERS?",
		}),

	// Structure
	folderRule("src_folder", "structure", config.SeverityWarning, map[string][]string{
		"go":     {"cmd/", "pkg/", "internal/"},
		"nodejs": {"src/", "lib/"},
		"*":      {"src/"},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create src/ folder?",
	}),
	multiRule("tests_folder", "structure", config.SeverityWarning,
		map[string][]string{
			"nodejs": {"test/", "tests/", "__tests__/"},
			"go":     {"*_test.go"},
			"*":      {"tests/"},
		}, FixSpec{
			Type:   FixTypeMulti,
			Prompt: "Create tests/ folder?",
		}),
	folderRule("docs_folder", "structure", config.SeverityWarning, map[string][]string{
		"*": {"docs/"},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create docs/ folder?",
	}),
	folderRule("scripts_folder", "structure", config.SeverityWarning, map[string][]string{
		"*": {"scripts/"},
	}, FixSpec{
		Type:   FixTypeFolder,
		Prompt: "Create scripts/ folder?",
	}),

	/// DevOps

	fileRule("dockerfile", "devops", config.SeverityWarning, map[string][]string{
		"*": {"Dockerfile"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create Dockerfile?",
		}),
	fileRule("dockerignore", "devops", config.SeverityWarning, map[string][]string{
		"*": {".dockerignore"},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create .dockerignore?",
		}),
	fileRule("docker_compose", "devops", config.SeverityWarning, map[string][]string{
		"*": {
			"docker-compose.yml",
			"docker-compose.yaml",
		},
	},
		FixSpec{
			Type:   FixTypeFile,
			Prompt: "Create docker-compose?",
		}),
}

func AllRules() []Rule {
	return all
}

func GetRulesByCategory(category string) []Rule {
	var filtered []Rule
	for _, rule := range all {
		if rule.Category == category {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

func GetRuleByID(id string) (Rule, bool) {
	for _, rule := range all {
		if rule.ID == id {
			return rule, true
		}
	}
	return Rule{}, false
}
