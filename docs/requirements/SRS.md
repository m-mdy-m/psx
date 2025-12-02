# Software Requirements Specification
## PSX - Project Structure Checker

- **Version:** 1.0.0  
- **Author:** Genix (m-mdy-m)  
- **Contact:** bitsgenix@gmail.com  
- **Repository:** [m-mdy-m/psx](https://github.com/m-mdy-m/psx)
- **Date:** 2025-12-02
- **Status:** Draft  
- **Confidentiality:** Internal

---

## Table of Contents

1. [Introduction](#introduction)
2. [Overall Description](#overall-description)
3. [System Features](#system-features)
4. [External Interface Requirements](#external-interface-requirements)
5. [System Requirements](#system-requirements)
6. [Appendices](#appendices)

---

## 1. Introduction

### 1.1 Purpose

This document describes the software requirements for PSX (Project Structure X), a command-line tool designed to validate and standardize project structures across different programming languages and frameworks.

The intended audience includes:
- Developers working on PSX
- Contributors wanting to understand the system
- Users looking for detailed technical information
- Anyone evaluating PSX for their workflow

### 1.2 Document Conventions

Throughout this document:
- **MUST** indicates absolute requirements
- **SHOULD** indicates recommended but not mandatory features
- **MAY** indicates optional features
- Code examples are shown in monospace font
- Priority levels: High (P0), Medium (P1), Low (P2)

### 1.3 Project Scope

PSX aims to solve a common problem in software development: inconsistent project structures. When working across multiple projects or in teams, maintaining a standard structure becomes challenging.

**What PSX does:**
- Automatically detects project type from files
- Validates project structure against configurable rules
- Provides actionable feedback with clear severity levels
- Can automatically fix common structural issues
- Supports multiple programming languages and frameworks
- Integrates into existing development workflows

**What PSX does NOT do:**
- Code quality analysis or linting
- Dependency management
- Build system configuration
- Runtime behavior checking
- Performance profiling

### 1.4 References

- Go Programming Language: https://go.dev
- GitHub Actions: https://docs.github.com/actions
- GitLab CI: https://docs.gitlab.com/ee/ci
- SARIF Format: https://sarifweb.azurewebsites.net
- Semantic Versioning: https://semver.org

---

## 2. Overall Description

### 2.1 Product Perspective

PSX is a standalone command-line tool that operates independently but can integrate into larger development workflows. It's designed with simplicity in mind - single binary, no dependencies, works everywhere.

**System Context:**
```
┌─────────────┐
│  Developer  │
└──────┬──────┘
       │
       ↓
┌─────────────────────────────────┐
│         PSX CLI Tool            │
│  ┌───────────────────────────┐  │
│  │  Project Type Detector    │  │
│  ├───────────────────────────┤  │
│  │  Rules Engine             │  │
│  ├───────────────────────────┤  │
│  │  Auto-fixer               │  │
│  ├───────────────────────────┤  │
│  │  Report Generator         │  │
│  └───────────────────────────┘  │
└──────────┬──────────────────────┘
           │
           ↓
    ┌─────────────┐
    │   Project   │
    │   Files     │
    └─────────────┘
```

### 2.2 Product Functions

At a high level, PSX performs these main functions:

1. **Detection**: Analyzes project files to determine programming language and framework
2. **Validation**: Checks project structure against rules (built-in or custom)
3. **Reporting**: Presents findings in various formats (table, JSON, HTML, etc.)
4. **Fixing**: Automatically corrects common structural issues when possible
5. **Configuration**: Allows customization through YAML configuration files

### 2.3 User Classes and Characteristics

**Individual Developers (Primary)**
- Technical expertise: Intermediate to advanced
- Usage frequency: Daily to weekly
- Main needs: Quick validation, minimal setup, clear feedback
- Characteristics: Values speed and simplicity over advanced features

**Team Leads (Secondary)**
- Technical expertise: Advanced
- Usage frequency: Weekly for setup, automated checks run continuously
- Main needs: Enforce standards across team, CI/CD integration
- Characteristics: Needs customization and reporting features

**Open Source Maintainers (Secondary)**
- Technical expertise: Advanced
- Usage frequency: Per contribution/PR
- Main needs: Consistent structure across contributions
- Characteristics: May have specific requirements per project

**DevOps Engineers (Tertiary)**
- Technical expertise: Advanced
- Usage frequency: Setup once, runs automatically
- Main needs: CI/CD integration, multiple output formats
- Characteristics: Needs reliability and clear exit codes

### 2.4 Operating Environment

PSX operates in the following environments:

**Operating Systems:**
- Linux (Ubuntu 20.04+, Debian 10+, Fedora 35+, Arch)
- macOS (11.0+, both Intel and Apple Silicon)
- Windows (10+, Windows Server 2019+)

**Hardware Requirements:**
- Minimum: 1 CPU core, 100MB RAM, 20MB disk
- Recommended: 2+ CPU cores, 256MB RAM, 50MB disk

**Software Dependencies:**
- None for binary usage
- Go 1.23+ for building from source

**Integration Points:**
- Git repositories
- CI/CD systems (GitHub Actions, GitLab CI, Jenkins, etc.)
- IDEs (through SARIF output)
- Pre-commit hooks

### 2.5 Design and Implementation Constraints

**Technical Constraints:**
- Must be written in Go for cross-platform compilation
- Binary size should remain under 20MB
- No external runtime dependencies allowed
- Must work offline (no network calls required)

**Business Constraints:**
- Open source (MIT license)
- Free for all users
- Self-funded development

**Regulatory Constraints:**
- Must not collect user data
- Must not execute arbitrary code from config files
- Should follow platform security best practices

### 2.6 Assumptions and Dependencies

**Assumptions:**
- Users have basic command-line knowledge
- Projects use standard file naming conventions
- Git is commonly used (but not required)
- UTF-8 encoding for text files

**Dependencies:**
- Go standard library
- Third-party Go packages (see go.mod):
  - gopkg.in/yaml.v3 for config parsing
  - github.com/spf13/cobra for CLI
  - github.com/fatih/color for terminal output
  - (may change during development)

---

## 3. System Features

### 3.1 Project Type Detection

**Description:**  
PSX automatically identifies what kind of project it's analyzing by examining specific files and their content. This is the first step in the validation process.

**Priority:** P0 (Critical)

**Functional Requirements:**

**FR-3.1.1: Detect Node.js Projects**
- **Input:** Project directory path
- **Process:** Check for package.json file existence
- **Output:** Project type = "nodejs"
- **Additional detection:**
  - Distinguish npm vs yarn vs pnpm (lock files)
  - Detect if TypeScript is used (tsconfig.json)
  - Identify frontend frameworks (check dependencies)

**FR-3.1.2: Detect Go Projects**
- **Input:** Project directory path
- **Process:** Check for go.mod file
- **Output:** Project type = "go"
- **Additional detection:**
  - Parse go.mod to get module name
  - Check for cmd/ directory structure
  - Identify if it's a library or application

**FR-3.1.3: Detect Rust Projects**
- **Input:** Project directory path
- **Process:** Check for Cargo.toml file
- **Output:** Project type = "rust"
- **Additional detection:**
  - Parse Cargo.toml for project name
  - Check if workspace or single crate

**FR-3.1.4: Detect Python Projects**
- **Input:** Project directory path
- **Process:** Check for requirements.txt, setup.py, or pyproject.toml
- **Output:** Project type = "python"
- **Additional detection:**
  - Distinguish pip vs poetry vs pipenv
  - Check for virtual environment markers

**FR-3.1.5: Detect Java Projects**
- **Input:** Project directory path
- **Process:** Check for pom.xml or build.gradle
- **Output:** Project type = "java"
- **Additional detection:**
  - Distinguish Maven vs Gradle
  - Check for Spring Boot indicators

**FR-3.1.6: Detect Mixed Projects**
- **Input:** Project directory path
- **Process:** If multiple language indicators found
- **Output:** Project type = "mixed", languages = ["go", "nodejs", ...]
- **Priority:** Secondary language gets validation rules too

**Acceptance Criteria:**
- Detection accuracy > 95% for supported languages
- Detection time < 100ms for projects with <10k files
- Should not fail on permission errors (warn and continue)
- Unknown projects should be marked as "generic" not error

---

### 3.2 Rule Validation Engine

**Description:**  
The core of PSX. It runs configured rules against the project structure and collects results.

**Priority:** P0 (Critical)

**Functional Requirements:**

**FR-3.2.1: Load Configuration**
- **Input:** psx.yml file or default config
- **Process:** Parse YAML, validate schema, merge with defaults
- **Output:** Configuration object with active rules
- **Error handling:** Invalid YAML should show clear error with line number

**FR-3.2.2: Rule Execution**
- **Input:** Configuration + project path
- **Process:** 
  1. Load all applicable rules for detected project type
  2. Execute rules in priority order
  3. Collect results (pass/fail/warning)
  4. Aggregate statistics
- **Output:** List of validation results
- **Performance:** Should handle 100+ rules in <1 second

**FR-3.2.3: Severity Levels**
Rules can have three severity levels:

**Error (Critical):**
- Blocks the "passing" status
- Exit code 1 if any errors found
- Examples: Missing README, missing LICENSE, invalid package.json

**Warning (Important but not blocking):**
- Doesn't block passing status
- Exit code 0 but shows in report
- Examples: Missing tests folder, no CI config, missing docs

**Info (Suggestions):**
- Just informational
- Exit code 0
- Examples: Could use EditorConfig, consider pre-commit hooks

**FR-3.2.4: Rule Categories**

Rules are organized into categories:

1. **General Rules:**
   - readme_required
   - license_required
   - gitignore_required
   - changelog_recommended

2. **Structure Rules:**
   - src_folder_required
   - tests_folder_required
   - docs_folder_recommended
   - proper_folder_naming

3. **Documentation Rules:**
   - adr_recommended
   - api_docs_required (for libraries)
   - contributing_guide_recommended

4. **CI/CD Rules:**
   - ci_config_recommended
   - workflow_syntax_valid

5. **Quality Rules:**
   - pre_commit_hooks_recommended
   - editorconfig_recommended
   - code_owners_recommended

6. **Language-Specific Rules:**
   - **Node.js:** package_json_valid, dependencies_up_to_date, scripts_defined
   - **Go:** go_mod_valid, go_sum_present, proper_module_structure
   - **Rust:** cargo_toml_valid, proper_crate_structure
   - **Python:** requirements_present, setup_py_valid
   - **Java:** pom_xml_valid, proper_maven_structure

**FR-3.2.5: Custom Rules**
- Users can define custom rules in config
- Format:
```yaml
custom_rules:
  - name: "my_custom_rule"
    description: "Check for terraform files"
    severity: warning
    check: "file_exists:terraform/main.tf"
```
- Supported check types: file_exists, folder_exists, file_contains, folder_not_empty

**Acceptance Criteria:**
- Rules execute independently (one failure doesn't stop others)
- Rule execution is deterministic (same input = same output)
- Results include: rule ID, severity, message, location, suggestion
- Can filter results by severity level
- Performance: <5 seconds for projects with 10k files

---

### 3.3 Auto-Fix Functionality

**Description:**  
PSX can automatically fix certain structural issues instead of just reporting them.

**Priority:** P1 (High)

**Functional Requirements:**

**FR-3.3.1: Safe Fixing**
- Never modify files without user confirmation (except in CI mode)
- Provide dry-run mode to preview changes
- Create backups before modifications (optional flag)
- Atomic operations (all or nothing)

**FR-3.3.2: Fixable Issues**

**File Creation:**
- Create missing README.md (from template)
- Create missing LICENSE (user chooses type)
- Create missing .gitignore (language-specific template)
- Create missing folders (src/, tests/, docs/)

**File Modification:**
- Add missing sections to README (if using template structure)
- Update .gitignore with recommended patterns
- Fix simple JSON syntax errors in package.json

**Template System:**
- README templates for each language
- LICENSE templates (MIT, Apache-2.0, GPL-3.0, etc.)
- .gitignore templates from github/gitignore repo patterns

**FR-3.3.3: Interactive Mode**
```
PSX found 5 fixable issues. Fix them?

[✓] Create README.md
    Template: nodejs-basic
    Location: ./README.md
    
[ ] Create LICENSE
    Choose type: 
    [1] MIT
    [2] Apache-2.0
    [3] GPL-3.0
    [4] Skip
    
[✓] Create tests/ folder
    Location: ./tests
    
Apply changes? (y/n/preview)
```

**FR-3.3.4: Non-Interactive Mode**
- For CI/CD usage
- Flag: `--fix-all --no-interactive`
- Uses sensible defaults
- Logs all actions

**Acceptance Criteria:**
- No fixes applied without explicit permission (except CI mode)
- Clear before/after diff shown
- Rollback mechanism available
- Fix actions logged for audit
- Templates customizable through config

---

### 3.4 Reporting System

**Description:**  
PSX presents validation results in various formats suitable for different use cases.

**Priority:** P1 (High)

**Functional Requirements:**

**FR-3.4.1: Table Output (Default)**
Human-readable terminal output with colors and formatting:
```
PSX v1.0.0 - Project Structure Checker
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Detected: Node.js (TypeScript)
Config: psx.yml ✓

ERRORS (2)
──────────────────────────────────────
✗ README_MISSING
  Location: ./
  Message: README.md file not found
  Fix: psx fix --rule readme_required
  
✗ TESTS_FOLDER_MISSING
  Location: ./
  Message: tests/ or __tests__/ folder not found

WARNINGS (3)
──────────────────────────────────────
⚠ LICENSE_MISSING
  Suggestion: Add a LICENSE file

⚠ CI_CONFIG_MISSING
  Location: .github/workflows/
  
⚠ NO_DOCUMENTATION
  Consider adding docs/ folder

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Summary: 2 errors, 3 warnings, 0 info
Status: FAILED ✗

Estimated fix time: ~3 minutes
Run: psx fix --interactive
```

**FR-3.4.2: JSON Output**
Machine-readable format for programmatic use:
```json
{
  "version": "1.0.0",
  "project": {
    "type": "nodejs",
    "subtype": "typescript",
    "path": "/path/to/project"
  },
  "config": {
    "file": "psx.yml",
    "valid": true
  },
  "results": [
    {
      "rule_id": "README_MISSING",
      "severity": "error",
      "category": "general",
      "message": "README.md file not found",
      "location": "./",
      "fixable": true,
      "fix_command": "psx fix --rule readme_required"
    }
  ],
  "summary": {
    "total_rules": 23,
    "errors": 2,
    "warnings": 3,
    "info": 0,
    "passed": 18
  },
  "status": "failed",
  "execution_time_ms": 142
}
```

**FR-3.4.3: SARIF Output**
For IDE integration (VSCode, IntelliJ, etc.):
```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "PSX",
          "version": "1.0.0",
          "informationUri": "https://github.com/m-mdy-m/psx"
        }
      },
      "results": [
        {
          "ruleId": "README_MISSING",
          "level": "error",
          "message": {
            "text": "README.md file not found"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "./"
                }
              }
            }
          ]
        }
      ]
    }
  ]
}
```

**FR-3.4.4: HTML Report**
For sharing and archiving:
- Summary at top
- Expandable sections per category
- Color-coded severity
- Clickable links to locations
- Export timestamp
- Can be opened in any browser

**FR-3.4.5: Markdown Report**
For documentation and GitHub issues:
- GitHub-flavored markdown
- Suitable for PR comments
- Can be copy-pasted

**Acceptance Criteria:**
- All formats contain same information
- Output format selectable via flag: `--output json`
- Can redirect to file: `psx check --output json > report.json`
- Table output respects terminal width
- Colors disabled when not a TTY
- Progress indicators for long operations

---

### 3.5 Configuration System

**Description:**  
PSX uses YAML configuration files to customize behavior, rules, and templates.

**Priority:** P0 (Critical)

**Functional Requirements:**

**FR-3.5.1: Config File Location**
PSX searches for configuration in this order:
1. `psx.yml` in current directory
2. `.psx.yml` in current directory
3. `psx.yml` in parent directories (up to git root)
4. `~/.config/psx/psx.yml` (user config)
5. Built-in defaults

**FR-3.5.2: Config Schema**
```yaml
version: 1

# Project type override (optional)
project_type: "nodejs"  # auto-detect if not specified

# Rules configuration
rules:
  # General
  general:
    readme_required: error
    license_required: warning
    gitignore_required: warning
    changelog_recommended: info
    
  # Structure
  structure:
    src_folder_required: warning
    tests_folder_required: error
    docs_folder_recommended: info
    
  # Documentation
  documentation:
    adr_recommended: false  # disabled
    srs_required: false
    api_docs_required: warning
    
  # CI/CD
  cicd:
    github_actions_recommended: info
    gitlab_ci_recommended: false
    
  # Quality
  quality:
    pre_commit_hooks: info
    editorconfig: info
    
  # Language-specific
  nodejs:
    package_json_valid: error
    dependencies_up_to_date: warning
    scripts_defined: info
    
  go:
    go_mod_valid: error
    go_sum_present: warning
    proper_structure: warning

# Paths to ignore
ignore:
  - node_modules/
  - vendor/
  - .git/
  - dist/
  - build/
  - "*.min.js"

# Custom rules
custom_rules:
  - name: "TERRAFORM_PRESENT"
    description: "Check for Terraform configuration"
    severity: info
    check: "folder_exists:terraform"
    
  - name: "DOCKER_PRESENT"
    description: "Check for Dockerfile"
    severity: warning
    check: "file_exists:Dockerfile"

# Templates customization
templates:
  readme:
    language: "nodejs"
    sections:
      - title
      - description
      - installation
      - usage
      - contributing
      - license
      
  license:
    default: "MIT"
    author: "Mahdi Mohamadi"
    email: "bitsgenix@gmail.com"

# Fixer behavior
fix:
  interactive: true
  create_backups: false
  default_license: "MIT"
```

**FR-3.5.3: Config Validation**
- Schema validation on load
- Clear error messages for invalid YAML
- Warning for unknown keys (typos)
- Type checking (string vs bool vs array)

**FR-3.5.4: Config Generation**
```bash
psx init                    # interactive
psx init --template nodejs  # use nodejs template
psx init --minimal          # bare minimum config
```

Interactive mode asks:
- Project type?
- Strict or relaxed rules?
- Enable custom rules?
- CI/CD platform?

**FR-3.5.5: Config Inheritance**
- User config (~/.config/psx/psx.yml) as defaults
- Project config overrides user config
- CLI flags override everything

**Acceptance Criteria:**
- Invalid config prevents execution with clear error
- Config can be validated separately: `psx config validate`
- Can show effective config: `psx config show`
- Templates can be listed: `psx config templates`

---

### 3.6 CLI Interface

**Description:**  
Command-line interface for interacting with PSX.

**Priority:** P0 (Critical)

**Functional Requirements:**

**FR-3.6.1: Commands**
```bash
# Check project
psx check
psx check --verbose
psx check --level error        # only errors
psx check --output json
psx check --config custom.yml

# Fix issues
psx fix
psx fix --interactive
psx fix --dry-run
psx fix --rule README_MISSING
psx fix --all

# Initialize config
psx init
psx init --template nodejs
psx init --force

# List rules
psx rules
psx rules --category structure
psx rules --json

# Config management
psx config validate
psx config show
psx config templates

# Version info
psx version
psx version --verbose

# Help
psx --help
psx check --help
```

**FR-3.6.2: Flags**

Global flags:
- `--verbose, -v`: More detailed output
- `--quiet, -q`: Minimal output
- `--config <file>`: Use specific config file
- `--no-color`: Disable colored output
- `--version`: Show version

Check flags:
- `--level <error|warning|info>`: Filter by severity
- `--output <table|json|sarif|html|markdown>`: Output format
- `--fail-on <error|warning>`: When to exit with code 1

Fix flags:
- `--interactive, -i`: Ask before each fix
- `--dry-run`: Show what would be fixed
- `--rule <rule_id>`: Fix specific rule
- `--all`: Fix all fixable issues
- `--no-backup`: Don't create backups

**FR-3.6.3: Exit Codes**
- `0`: Success (no errors, or only warnings/info)
- `1`: Validation failed (errors found)
- `2`: Configuration error
- `3`: File system error
- `4`: Invalid arguments

**FR-3.6.4: Shell Completion**
Generate completion scripts:
```bash
psx completion bash > /etc/bash_completion.d/psx
psx completion zsh > ~/.zsh/completion/_psx
psx completion fish > ~/.config/fish/completions/psx.fish
```

**Acceptance Criteria:**
- All commands have help text
- Flag errors show suggestions
- Progress indicators for long operations
- Keyboard interrupts handled gracefully
- Works in both interactive and non-interactive terminals

---

## 4. External Interface Requirements

### 4.1 User Interfaces

**4.1.1 Terminal Interface**
- Primary interface is command-line
- Supports 80-column terminals (minimum)
- Respects terminal color capabilities
- Uses Unicode symbols when supported (with ASCII fallback)
- Progress bars for operations >2 seconds

**4.1.2 Interactive Prompts**
When in interactive mode:
- Yes/No questions: `(y/n)` with default shown
- Multiple choice: numbered list with arrow keys
- Text input: with validation
- Cancellable with Ctrl+C

### 4.2 Software Interfaces

**4.2.1 File System**
- **Read operations:**
  - Config files (YAML)
  - Project files (for detection and validation)
  - Template files
- **Write operations:**
  - Create missing files (with permission)
  - Create folders
  - Generate reports
- **Permissions:**
  - Respects file permissions
  - Shows clear error if can't read/write

**4.2.2 Git Integration**
- Detect `.git` folder to find project root
- Read `.gitignore` for ignore patterns
- Can validate `.gitignore` content
- Does NOT require git to be installed

**4.2.3 CI/CD Integration**
- **GitHub Actions:**
  - Exit code indicates pass/fail
  - JSON/SARIF output parseable
  - Example workflow provided
  
- **GitLab CI:**
  - Same exit code behavior
  - Can output JUnit XML format
  - Example .gitlab-ci.yml provided
  
- **Jenkins:**
  - Compatible with sh step
  - Warnings plugin support

**4.2.4 IDE Integration**
- SARIF output works with:
  - VS Code (via extensions)
  - IntelliJ IDEA
  - Visual Studio
- Can be run as external tool
- File locations use absolute paths

### 4.3 Communication Interfaces

**4.3.1 Standard Streams**
- `stdout`: Normal output (reports, results)
- `stderr`: Errors and warnings
- `stdin`: Interactive input when needed

**4.3.2 Environment Variables**
- `PSX_CONFIG`: Override config file location
- `PSX_NO_COLOR`: Disable colors (if set to any value)
- `PSX_LOG_LEVEL`: Set log verbosity (debug, info, warn, error)
- `CI`: Detect CI environment (auto non-interactive mode)

**4.3.3 Network**
- PSX does NOT make network requests
- All functionality works offline
- Templates bundled in binary

---

## 5. System Requirements

### 5.1 Performance Requirements

**PER-1: Execution Speed**
- Scan projects with 1,000 files in <1 second
- Scan projects with 10,000 files in <5 seconds
- Scan projects with 100,000 files in <30 seconds
- Fix operations complete in <2 seconds

**PER-2: Resource Usage**
- Binary size: <20MB
- Memory usage: <50MB for typical projects, <200MB for huge projects
- CPU usage: Should use all available cores for large projects
- Disk I/O: Minimized through caching

**PER-3: Startup Time**
- Cold start: <100ms
- Config load: <50ms
- First scan: <200ms additional

### 5.2 Safety Requirements

**SAF-1: File Safety**
- Never delete files
- Never modify files without confirmation (except CI mode)
- Create backups when modifying (optional)
- Atomic write operations (temp file + rename)

**SAF-2: Error Handling**
- Graceful degradation on permission errors
- Clear error messages with suggestions
- No crashes on malformed input
- Handle disk full scenarios

**SAF-3: Security**
- No code execution from config files
- Validate all user inputs
- Safe path traversal (no escaping project root)
- No secrets in logs

### 5.3 Software Quality Attributes

**Reliability:**
- 99% success rate for supported scenarios
- Deterministic results (same input = same output)
- Handles edge cases gracefully

**Usability:**
- Can be learned in 5 minutes
- Clear error messages with actionable suggestions
- Sensible defaults (works without config)
- Comprehensive help system

**Maintainability:**
- Code coverage >80%
- Clear code structure
- Comprehensive documentation
- Easy to add new rules

**Portability:**
- Single binary per platform
- No platform-specific code except where necessary
- Consistent behavior across OS

**Scalability:**
- Works with projects from 10 to 100,000 files
- Rule execution parallelized
- Memory usage grows linearly with project size

---

## 6. Appendices

### Appendix A: Glossary

- **Rule**: A validation check for project structure
- **Severity**: Level of importance (error, warning, info)
- **Auto-fix**: Automatically correcting an issue
- **Project Type**: Language/framework detection result
- **Config**: YAML configuration file
- **Template**: Pre-defined content for file generation
- **SARIF**: Static Analysis Results Interchange Format
- **CI/CD**: Continuous Integration/Continuous Deployment
- **ADR**: Architecture Decision Record
- **SRS**: Software Requirements Specification

### Appendix B: Acronyms

- CLI: Command Line Interface
- YAML: YAML Ain't Markup Language
- JSON: JavaScript Object Notation
- UTF-8: Unicode Transformation Format 8-bit
- PR: Pull Request
- OS: Operating System
- IDE: Integrated Development Environment

### Appendix C: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | Dec 2024 | m-mdy-m | Initial version |

### Appendix D: Open Questions

(Questions that need resolution before implementation)

1. Should we support plugins for custom rules?
   - Pro: Extensibility
   - Con: Complexity, security concerns
   - Decision: Not in v1.0, maybe later

2. Web interface for viewing reports?
   - Pro: Better for non-technical users
   - Con: Scope creep, more maintenance
   - Decision: HTML export sufficient for v1.0

3. Automatic updates?
   - Pro: Users stay current
   - Con: Potential security issue
   - Decision: Manual updates only

4. Cloud sync for team configs?
   - Pro: Easier team coordination
   - Con: Infrastructure, privacy concerns
   - Decision: Git is sufficient for config sharing

### Appendix E: Future Enhancements

(Not in scope for v1.0 but potential future features)

- Plugin system for custom rules
- AI-powered suggestions
- Web dashboard for teams
- Integration with project management tools
- Automated pull request generation for fixes
- Template marketplace
- Multi-project scanning
- Diff mode (compare two project structures)

---

**Document Status:** Draft → **Under Review** → Approved  
**Review Required By:** Core contributors  
**Approval Required By:** Maintainer (m-mdy-m)

---

*This SRS is a living document and will be updated as requirements evolve.*