# Software Requirements Specification
## PSX - Project Structure Checker

**Version:** 1.0.0  
**Author:** Genix (m-mdy-m)  
**Date:** 2025-12-17  
**Status:** Release

---

## 1. Introduction

### 1.1 Purpose

PSX is a command-line tool that validates and standardizes project structures across different programming languages. It automatically detects project types, checks for essential files and folders, and can fix common structural issues.

**Target Audience:**
- Developers maintaining consistent project structures
- Teams enforcing project standards
- Open source maintainers ensuring contributor consistency
- DevOps engineers integrating structure validation into CI/CD

### 1.2 Scope

**What PSX Does:**
- Detects project type (Node.js, Go, Generic)
- Validates against 47 built-in rules
- Provides clear, actionable feedback
- Automatically creates missing files and folders
- Generates templates based on project type
- Caches project information for reuse

**What PSX Does NOT Do:**
- Code quality analysis or linting
- Dependency management
- Build system configuration
- Runtime behavior checking

### 1.3 Definitions

- **Rule** - A validation check for project structure
- **Severity** - Importance level: error, warning, info
- **Auto-fix** - Automatically creating/modifying files
- **Project Type** - Detected language (nodejs, go, generic)
- **Template** - Pre-defined file content
- **Project Cache** - Stored project information (.psx-project.yml)

---

## 2. System Overview

### 2.1 Core Components

```
┌─────────────────────────────────────────┐
│             CLI Interface                │
│  (check, fix, project commands)         │
└────────────────┬────────────────────────┘
                 │
    ┌────────────┼────────────┐
    │            │            │
    ▼            ▼            ▼
┌─────────┐  ┌──────────┐  ┌──────────┐
│Detector │  │  Rules   │  │  Fixer   │
│         │  │  Engine  │  │          │
└─────────┘  └──────────┘  └──────────┘
    │            │            │
    └────────────┼────────────┘
                 ▼
         ┌───────────────┐
         │   Resources   │
         │  & Templates  │
         └───────────────┘
```

### 2.2 Supported Platforms

- **Operating Systems:** Linux, macOS, Windows
- **Architectures:** amd64, arm64
- **Languages:** Node.js, Go, Generic projects
- **Minimum Requirements:** None (single binary)

---

## 3. Functional Requirements

### 3.1 Project Detection

**FR-1: Automatic Type Detection**
- **Input:** Project directory path
- **Process:** Scan for language-specific files
  - Node.js: package.json
  - Go: go.mod
  - Generic: fallback for unknown types
- **Output:** Detected project type and metadata
- **Performance:** < 100ms for typical projects

**FR-2: Detection Caching**
- Cache detection results in `.psx-project.yml`
- Reuse cached data to avoid repeated scanning
- Update cache when changes detected

### 3.2 Rule Validation

**FR-3: Rule Execution**
- **47 Built-in Rules** across 6 categories:
  - General (4): readme, license, gitignore, changelog
  - Structure (5): src_folder, tests_folder, docs_folder, scripts_folder, env_example
  - Documentation (10): adr, contributing, api_docs, security, etc.
  - CI/CD (4): ci_config, github_actions, renovate, dependabot
  - Quality (12): editorconfig, prettier, eslint, pre_commit, etc.
  - DevOps (6): dockerfile, docker_compose, kubernetes, etc.

- **Parallel Execution:** Rules run concurrently
- **Configurable Severity:** error, warning, info, or disabled
- **Project-Specific:** Rules adapt to detected language

**FR-4: Rule Results**
- Each rule returns: pass/fail, severity, message, fix hint
- Results grouped by severity
- Summary statistics: total, passed, errors, warnings, info

### 3.3 Configuration

**FR-5: YAML Configuration**
- **Format:** Simple YAML syntax
- **Location Search Order:**
  1. `psx.yml` in current directory
  2. `.psx.yml` in current directory  
  3. Parent directories up to git root
  4. `~/.config/psx/psx.yml`
  5. Built-in defaults

**FR-6: Rule Customization**
```yaml
version: 1
rules:
  readme: error      # Enable as error
  license: warning   # Enable as warning
  docs_folder: info  # Enable as info
  prettier: false    # Disable

ignore:
  - node_modules/
  - dist/
```

### 3.4 Auto-Fix

**FR-7: File Creation**
- Create missing files from templates:
  - README.md (language-specific)
  - LICENSE (MIT, Apache, GPL, BSD)
  - .gitignore (language-specific patterns)
  - CHANGELOG.md
  - CONTRIBUTING.md
  - CODE_OF_CONDUCT.md
  - SECURITY.md
  - And 30+ other files

**FR-8: Folder Creation**
- Create missing directories:
  - src/, tests/, docs/, scripts/
  - .github/workflows/
  - docs/adr/
  - k8s/, infra/

**FR-9: Fix Modes**
- **Interactive (default):** Ask before each change
- **Dry-run:** Preview without applying
- **Non-interactive:** Fix all automatically (CI mode)
- **Selective:** Fix specific rules only

**FR-10: Safety Measures**
- Never overwrite existing files
- Atomic operations
- Optional backup creation
- Clear change preview

### 3.5 Reporting

**FR-11: Table Output (Default)**
- Human-readable terminal format
- Color-coded by severity
- Grouped by error/warning/info
- Summary with fix suggestions

**FR-12: JSON Output**
- Machine-readable format
- Complete rule results
- Summary statistics
- Suitable for CI/CD parsing

### 3.6 Project Information

**FR-13: Information Collection**
- Gather project metadata:
  - Name, description
  - Author name and email
  - GitHub username
  - Repository URL
  - License type

**FR-14: Smart Defaults**
- Extract from git config
- Parse from package files
- Prompt interactively when needed
- Cache for reuse

**FR-15: Project Commands**
- `psx project show` - Display cached info
- `psx project edit` - Update information
- `psx project reset` - Clear cache

---

## 4. Non-Functional Requirements

### 4.1 Performance

**NFR-1: Execution Speed**
- Projects < 1,000 files: < 1 second
- Projects < 10,000 files: < 5 seconds
- Concurrent rule execution
- Efficient file system access

**NFR-2: Resource Usage**
- Binary size: < 20MB
- Memory: < 50MB for typical projects
- Single binary, no dependencies
- Fast startup: < 100ms

### 4.2 Reliability

**NFR-3: Error Handling**
- Graceful degradation on errors
- Continue checking after failures
- Clear error messages
- Proper exit codes:
  - 0: Success
  - 1: Validation failed
  - 2: Configuration error
  - 3: File system error
  - 4: Invalid arguments

**NFR-4: Deterministic**
- Same input → same output
- No random behavior
- Consistent across platforms

### 4.3 Usability

**NFR-5: Ease of Use**
- Zero configuration for basic usage
- Sensible defaults
- Clear help text
- Self-documenting commands

**NFR-6: Developer Experience**
- Single command to check
- Single command to fix
- Fast feedback
- Actionable suggestions

### 4.4 Portability

**NFR-7: Cross-Platform**
- Identical behavior on Linux/macOS/Windows
- Single binary per platform
- No platform-specific requirements

**NFR-8: No Dependencies**
- Statically linked binary
- Runs on any system
- No runtime requirements
- Offline capable

### 4.5 Security

**NFR-9: Safe Operations**
- Never deletes files
- Validates all paths
- No code execution from config
- No network requests

**NFR-10: Privacy**
- No telemetry
- No data collection
- All operations local
- Project cache stays local

---

## 5. System Constraints

### 5.1 Technical Constraints

- Written in Go 1.25+
- Single binary architecture
- YAML for configuration
- Embedded resources (templates, messages)

### 5.2 Design Constraints

- CLI-only interface (no GUI)
- Synchronous command execution
- File-based configuration
- No plugin system (v1.0)

---

## 6. Use Cases

### 6.1 Individual Developer

**Scenario:** Developer starting a new Node.js project

```bash
$ mkdir my-app && cd my-app
$ psx check

Detected: nodejs (generic fallback)

ERRORS (2)
README_MISSING
LICENSE_MISSING

$ psx fix --interactive
Create README.md? y
Create LICENSE (MIT)? y

Summary: 2 fixed
```

### 6.2 Team Lead

**Scenario:** Enforce standards across team projects

1. Create team config: `.github/psx.yml`
2. Share via git repository
3. Team members run: `psx check --config .github/psx.yml`
4. CI/CD validates on PRs

### 6.3 CI/CD Integration

**Scenario:** Automated validation in GitHub Actions

```yaml
- name: Validate Structure
  run: |
    psx check --output json > results.json
    psx check --fail-on warning
```

### 6.4 Open Source Maintainer

**Scenario:** Ensure consistent contributor submissions

1. Add `psx.yml` to repository
2. Document in CONTRIBUTING.md
3. Contributors run `psx check` before PR
4. Automated check in CI

---

## 7. Future Enhancements (Not in v1)

- Plugin system for custom rules
- Multi-project scanning
- Git pre-commit hook integration
- Language support: Python, Rust, Java

---

## 8. Appendix

### 8.1 Complete Rule List

See `internal/config/embedded/rules.yml` for full rule definitions.

### 8.2 Exit Codes

```go
const (
    ExitSuccess = 0  // No errors
    ExitFailed  = 1  // Validation failed
    ExitConfig  = 2  // Config error
    ExitFS      = 3  // File system error
    ExitArgs    = 4  // Invalid arguments
)
```

### 8.3 Configuration Schema

```yaml
version: 1                    # Required: schema version

project:
  type: ""                    # Optional: override detection

rules:                        # Rule configuration
  <rule_id>: error|warning|info|false

ignore:                       # Glob patterns to ignore
  - pattern/

fix:                          # Fix behavior
  interactive: true
  backup: false
```
