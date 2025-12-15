package fixer

import (
	"fmt"
	"path/filepath"
	"time"

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

	// Create first ADR about ADRs
	firstADR := filepath.Join(adrPath, "0001-record-architecture-decisions.md")
	content := generateFirstADR()
	if err := shared.CreateFile(firstADR, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        adrPath,
		Description: "Created docs/adr/ folder with initial ADR",
	})

	return result, nil
}

// FixContributing creates CONTRIBUTING.md
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

	content := generateContributing()

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

// FixAPIDocsFolder creates API documentation folder
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
	content := generateAPIDocsReadme(ctx.ProjectType)
	if err := shared.CreateFile(readmePath, content); err != nil {
		// Non-critical
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        apiDocsPath,
		Description: "Created docs/api/ folder",
	})

	return result, nil
}

// Template generators
func generateFirstADR() string {
	today := time.Now().Format("2006-01-02")

	return fmt.Sprintf(`# 1. Record architecture decisions

Date: %s

## Status

Accepted

## Context

We need to record the architectural decisions made on this project.

## Decision

We will use Architecture Decision Records, as described by Michael Nygard.

## Consequences

See Michael Nygard's article, linked above. For a lightweight ADR toolset, see Nat Pryce's adr-tools.
`, today)
}

func generateADRTemplate() string {
	return `# [number]. [title]

Date: [YYYY-MM-DD]

## Status

[Proposed | Accepted | Deprecated | Superseded]

## Context

What is the issue that we're seeing that is motivating this decision or change?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

What becomes easier or more difficult to do because of this change?
`
}

func generateContributing() string {
	return `# Contributing

Thank you for your interest in contributing!

## How to Contribute

### Reporting Bugs

- Check if the bug has already been reported
- Create a detailed issue with:
  - Steps to reproduce
  - Expected behavior
  - Actual behavior
  - Environment details

### Suggesting Features

- Open an issue with the enhancement label
- Describe the feature and its use case
- Explain why it would be useful

### Code Contributions

1. Fork the repository
2. Create a feature branch (git checkout -b feature/amazing-feature)
3. Make your changes
4. Run tests
5. Commit your changes (git commit -m 'Add amazing feature')
6. Push to the branch (git push origin feature/amazing-feature)
7. Open a Pull Request

## Code Style

- Follow the existing code style
- Write clear commit messages
- Add tests for new features
- Update documentation as needed

## Questions?

Open an issue or reach out to the maintainers.
`
}

func generateAPIDocsReadme(projectType string) string {
	templates := map[string]string{
		"nodejs": `# API Documentation

## Overview

Document your public API here.

## Modules

### Module Name


// Example usage
const module = require('./module');

## Functions

### functionName(param1, param2)

Description of the function.

**Parameters:**
- param1 (Type): Description
- param2 (Type): Description

**Returns:** Return type and description

**Example:**
javascript
const result = functionName('value', 123);
`,
		"go": `# API Documentation

## Overview

Document your public API here.

## Packages

### package name


import "github.com/user/project/package"

## Functions

### FunctionName(param1, param2 string) error

Description of the function.

**Parameters:**
- param1 (string): Description
- param2 (string): Description

**Returns:** error

**Example:**
go
err := FunctionName("value", "other")
if err != nil {
    log.Fatal(err)
}
`,
	}

	if template, ok := templates[projectType]; ok {
		return template
	}

	return `# API Documentation

Document your public API here.

## Usage

Add examples and descriptions of your API.
`
}
