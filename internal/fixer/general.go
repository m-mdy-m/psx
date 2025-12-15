package fixer

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/m-mdy-m/psx/internal/shared"
)

// FixReadme creates a README.md file
func FixReadme(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "readme",
		Changes: []Change{},
	}

	readmePath := filepath.Join(ctx.ProjectPath, "README.md")

	exists, _ := shared.FileExists(readmePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Get project name from path
	projectName := filepath.Base(ctx.ProjectPath)

	content := generateReadme(projectName, ctx.ProjectType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        readmePath,
			Description: "Create README.md",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create README.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(readmePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        readmePath,
		Description: "Created README.md",
	})

	return result, nil
}

// FixLicense creates a LICENSE file
func FixLicense(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "license",
		Changes: []Change{},
	}

	licensePath := filepath.Join(ctx.ProjectPath, "LICENSE")

	exists, _ := shared.FileExists(licensePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	licenseType := "MIT"
	if ctx.Config.Fix.Interactive || ctx.Interactive {
		options := []string{"MIT", "Apache-2.0", "GPL-3.0", "Skip"}
		choice, _ := shared.PromptChoice("Choose license type:", options)
		if choice == "Skip" {
			result.Skipped = true
			return result, nil
		}
		licenseType = choice
	}

	content := generateLicense(licenseType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        licensePath,
			Description: fmt.Sprintf("Create LICENSE (%s)", licenseType),
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if err := shared.CreateFile(licensePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        licensePath,
		Description: fmt.Sprintf("Created LICENSE (%s)", licenseType),
	})

	return result, nil
}

// FixGitignore creates or updates .gitignore
func FixGitignore(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "gitignore",
		Changes: []Change{},
	}

	gitignorePath := filepath.Join(ctx.ProjectPath, ".gitignore")

	exists, _ := shared.FileExists(gitignorePath)

	content := generateGitignore(ctx.ProjectType)

	if exists {
		// TODO: Append missing patterns
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        gitignorePath,
			Description: "Create .gitignore",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .gitignore?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(gitignorePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        gitignorePath,
		Description: "Created .gitignore",
	})

	return result, nil
}

// FixChangelog creates a CHANGELOG.md file
func FixChangelog(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "changelog",
		Changes: []Change{},
	}

	changelogPath := filepath.Join(ctx.ProjectPath, "CHANGELOG.md")

	exists, _ := shared.FileExists(changelogPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := generateChangelog()

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        changelogPath,
			Description: "Create CHANGELOG.md",
			Content:     content,
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create CHANGELOG.md?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(changelogPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        changelogPath,
		Description: "Created CHANGELOG.md",
	})

	return result, nil
}

// Template generators
func generateReadme(projectName, projectType string) string {
	return fmt.Sprintf("HELLO THIS IS TEST REDME")
}

func generateLicense(licenseType string) string {
	year := time.Now().Year()

	templates := map[string]string{
		"MIT": fmt.Sprintf(`MIT License

Copyright (c) %d

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year),
	}

	if template, ok := templates[licenseType]; ok {
		return template
	}
	return templates["MIT"]
}

func generateGitignore(projectType string) string {
	common := `# Editor directories and files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db
`

	typeSpecific := map[string]string{
		"nodejs": `
# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.pnpm-debug.log*

# Build output
dist/
build/
.next/
.nuxt/

# Environment
.env
.env.local
.env.*.local

# Coverage
coverage/
*.lcov
`,
		"go": `
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Dependencies
vendor/

# Build
build/
bin/
`,
		"python": `
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
dist/
*.egg-info/

# Virtual environments
venv/
env/
ENV/

# Testing
.pytest_cache/
.coverage
htmlcov/
`,
	}

	if specific, ok := typeSpecific[projectType]; ok {
		return common + specific
	}
	return common
}

func generateChangelog() string {
	today := time.Now().Format("2006-01-02")

	return fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - %s

### Added
- Initial release

[Unreleased]: https://github.com/username/repo/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/username/repo/releases/tag/v0.1.0
`, today)
}

// Helper to format content for display
func FormatContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	displayed := strings.Join(lines[:maxLines], "\n")
	remaining := len(lines) - maxLines
	return fmt.Sprintf("%s\n... (%d more lines)", displayed, remaining)
}
