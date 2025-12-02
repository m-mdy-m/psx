# Functional Requirements
## PSX - Project Structure Checker

- **Version:** 1.0.0  
- **Author:** Genix (m-mdy-m)  
- **Contact:** bitsgenix@gmail.com  
- **Repository:** [m-mdy-m/psx](https://github.com/m-mdy-m/psx)
- **Date:** 2025-12-02
- **Status:** Draft  
- **Confidentiality:** Internal

---

## Overview

This document details all functional requirements for PSX. Each requirement is assigned a unique ID, priority, and acceptance criteria.

Requirements are categorized into:
1. Core Functionality
2. Detection System
3. Validation System
4. Auto-Fix System
5. Reporting System
6. Configuration System
7. CLI Interface

---

## Priority Levels

- **P0 (Critical)**: Must have for MVP, blocks release
- **P1 (High)**: Important, should be in v1.0
- **P2 (Medium)**: Nice to have, can wait for v1.1
- **P3 (Low)**: Future consideration

---

## 1. Core Functionality

### FR-1.1: Single Binary
Distribution
**Priority:** P0  
**Category:** Core

**Description:**  
PSX must be distributed as a single executable binary with no external dependencies.

**Requirements:**
- One binary per OS/architecture combination
- Binary size <20MB
- No runtime dependencies (no Python, Node, etc. required)
- Binary includes all templates and configurations

**Acceptance Criteria:**
- [ ] Binary runs on fresh system without installing anything else
- [ ] `file psx` shows it's a statically linked executable
- [ ] Binary size is under 20MB
- [ ] Works offline (no network required)

**Test Scenarios:**
```bash
# Fresh Ubuntu container
docker run -it ubuntu:20.04
# copy binary
./psx check  # should work without apt install anything
```

---

### FR-1.2: Cross-Platform Support
**Priority:** P0  
**Category:** Core

**Description:**  
PSX must work identically on Linux, macOS, and Windows.

**Requirements:**
- Support Linux (amd64, arm64)
- Support macOS (amd64, arm64/M1)
- Support Windows (amd64)
- Path handling works on all platforms
- Line endings handled correctly

**Acceptance Criteria:**
- [ ] Same command syntax on all platforms
- [ ] Config file format identical
- [ ] Output format consistent
- [ ] File operations respect OS conventions

**Test Scenarios:**
```bash
# Same project checked on all three OS
# Results should be identical (except path format)
Linux:   ./src/main.go
macOS:   ./src/main.go
Windows: .\src\main.go
```

---

### FR-1.3: Project Root Detection
**Priority:** P0  
**Category:** Core

**Description:**  
PSX must correctly identify the project root directory.

**Requirements:**
- Check current directory first
- Look for `.git` folder (most reliable)
- Look for language-specific markers (package.json, go.mod, etc.)
- Stop at filesystem root if no markers found
- Allow manual override with `--root` flag

**Acceptance Criteria:**
- [ ] Finds project root from subdirectory
- [ ] Works in projects without git
- [ ] Handles monorepos correctly
- [ ] Shows clear error if no project found

**Test Scenarios:**
```bash
# Run from subdirectory
cd my-project/src/utils
psx check  # should detect my-project/ as root

# Override
psx check --root /path/to/project
```

---

## 2. Detection System

### FR-2.1: Node.js Project Detection
**Priority:** P0  
**Category:** Detection

**Description:**  
Detect if project is Node.js based and identify variants.

**Detection Criteria:**
- Primary: `package.json` exists
- Variants:
  - npm: `package-lock.json` exists
  - yarn: `yarn.lock` exists
  - pnpm: `pnpm-lock.yaml` exists
  - bun: `bun.lockb` exists

**Additional Detection:**
- TypeScript: `tsconfig.json` exists
- Frontend framework:
  - React: "react" in dependencies
  - Vue: "vue" in dependencies
  - Angular: "@angular/core" in dependencies
  - Next.js: "next" in dependencies

**Output:**
```json
{
  "type": "nodejs",
  "package_manager": "pnpm",
  "typescript": true,
  "framework": "next"
}
```

**Acceptance Criteria:**
- [ ] Detects all package managers correctly
- [ ] Identifies TypeScript projects
- [ ] Recognizes major frameworks
- [ ] Handles multiple lock files (picks most recent)

---

### FR-2.2: Go Project Detection
**Priority:** P0  
**Category:** Detection

**Description:**  
Detect Go projects and identify structure type.

**Detection Criteria:**
- Primary: `go.mod` exists
- Parse module name from go.mod
- Check structure:
  - Application: `cmd/` directory exists
  - Library: No `cmd/`, has `pkg/` or just `.go` files in root
  - CLI tool: `main.go` in root or `cmd/`

**Output:**
```json
{
  "type": "go",
  "module": "github.com/m-mdy-m/psx",
  "structure": "application",
  "go_version": "1.23"
}
```

**Acceptance Criteria:**
- [ ] Parses go.mod correctly
- [ ] Identifies application vs library
- [ ] Handles malformed go.mod gracefully
- [ ] Extracts Go version requirement

---

### FR-2.3: Rust Project Detection
**Priority:** P0  
**Category:** Detection

**Description:**  
Detect Rust projects and workspace configuration.

**Detection Criteria:**
- Primary: `Cargo.toml` exists
- Parse package name and version
- Check if workspace:
  - `[workspace]` section in Cargo.toml
  - Multiple crates in subdirectories

**Output:**
```json
{
  "type": "rust",
  "name": "my-tool",
  "version": "0.1.0",
  "workspace": false,
  "edition": "2021"
}
```

**Acceptance Criteria:**
- [ ] Parses Cargo.toml correctly
- [ ] Identifies workspaces
- [ ] Handles TOML syntax errors gracefully
- [ ] Extracts Rust edition

---

### FR-2.4: Python Project Detection
**Priority:** P0  
**Category:** Detection

**Description:**  
Detect Python projects and dependency management tool.

**Detection Criteria:**
- Check for (in priority order):
  1. `pyproject.toml` (modern, preferred)
  2. `setup.py` (traditional)
  3. `requirements.txt` (simple)
  4. `Pipfile` (pipenv)

**Additional Detection:**
- Virtual environment: `.venv/`, `venv/`, `.virtualenv/`
- Django: `manage.py` exists
- Flask: "flask" in dependencies

**Output:**
```json
{
  "type": "python",
  "dependency_manager": "poetry",
  "has_venv": true,
  "framework": "django"
}
```

**Acceptance Criteria:**
- [ ] Detects all dependency file formats
- [ ] Identifies virtual environments
- [ ] Recognizes popular frameworks
- [ ] Works with Python 2 and 3 projects

---

### FR-2.5: Java Project Detection
**Priority:** P1  
**Category:** Detection

**Description:**  
Detect Java projects and build tool.

**Detection Criteria:**
- Maven: `pom.xml` exists
- Gradle: `build.gradle` or `build.gradle.kts` exists
- Check for Spring Boot: dependency on spring-boot-starter

**Output:**
```json
{
  "type": "java",
  "build_tool": "maven",
  "spring_boot": true
}
```

**Acceptance Criteria:**
- [ ] Distinguishes Maven vs Gradle
- [ ] Handles multi-module projects
- [ ] Identifies Spring Boot projects
- [ ] Parses build file versions

---

### FR-2.6: Mixed Project Detection
**Priority:** P1  
**Category:** Detection

**Description:**  
Handle projects with multiple languages (monorepos, polyglot projects).

**Detection Strategy:**
- Scan entire tree for language markers
- Identify primary language (most files, root-level configs)
- List secondary languages

**Example Scenario:**
```
project/
├── package.json          # Node.js frontend
├── go.mod                # Go backend
└── Dockerfile            # Deployment
```

**Output:**
```json
{
  "type": "mixed",
  "primary": "go",
  "secondary": ["nodejs"],
  "components": {
    "backend": "go",
    "frontend": "nodejs"
  }
}
```

**Acceptance Criteria:**
- [ ] Detects all languages present
- [ ] Identifies primary language correctly
- [ ] Applies rules for all detected languages
- [ ] Handles nested projects (monorepos)

---

## 3. Validation System

### FR-3.1: General Rules

#### FR-3.1.1: README Required
**Priority:** P0  
**Rule ID:** `README_REQUIRED`  
**Default Severity:** Error

**Check:**
- `README.md` or `README` or `readme.md` exists in project root

**Pass Criteria:**
- File exists
- File is not empty (>100 bytes)

**Fail Message:**
```
✗ README_REQUIRED
  Location: ./
  Message: No README file found in project root
  Suggestion: Create README.md with project description
  Fix: psx fix --rule README_REQUIRED
```

**Auto-fix:**
- Creates `README.md` from template
- Template includes: title, description, installation, usage, license

---

#### FR-3.1.2: LICENSE Required
**Priority:** P0  
**Rule ID:** `LICENSE_REQUIRED`  
**Default Severity:** Error

**Check:**
- `LICENSE` or `LICENSE.md` or `COPYING` exists in project root

**Pass Criteria:**
- File exists
- File contains recognized license text (MIT, Apache, GPL, etc.)

**Fail Message:**
```
✗ LICENSE_REQUIRED
  Location: ./
  Message: No LICENSE file found
  Suggestion: Add a license to specify usage terms
  Fix: psx fix --rule LICENSE_REQUIRED (will prompt for license type)
```

**Auto-fix:**
- Interactive: prompts user to choose license
- Non-interactive: uses config default (MIT)
- Fills in copyright holder from git config or config file

---

#### FR-3.1.3: .gitignore Required
**Priority:** P0  
**Rule ID:** `GITIGNORE_REQUIRED`  
**Default Severity:** Warning

**Check:**
- `.gitignore` exists in project root

**Pass Criteria:**
- File exists
- Contains language-specific patterns

**Additional Checks:**
- Warns if common patterns missing (node_modules/, .env, etc.)
- Validates pattern syntax

**Auto-fix:**
- Creates `.gitignore` with language-specific template
- Appends missing patterns if file exists

---

### FR-3.2: Structure Rules

#### FR-3.2.1: Source Folder Required
**Priority:** P0  
**Rule ID:** `SRC_FOLDER_REQUIRED`  
**Default Severity:** Warning

**Check:**
- Folder named `src/`, `app/`, `lib/`, or language-specific exists

**Language-Specific:**
- Node.js: `src/` or `lib/`
- Go: `cmd/`, `internal/`, `pkg/`
- Rust: `src/` (always)
- Python: Package folder or `src/`
- Java: `src/main/java/`

**Pass Criteria:**
- At least one source folder exists
- Folder is not empty

**Auto-fix:**
- Creates appropriate folder for detected language
- Does NOT move existing files (too risky)

---

#### FR-3.2.2: Tests Folder Required
**Priority:** P0  
**Rule ID:** `TESTS_FOLDER_REQUIRED`  
**Default Severity:** Error

**Check:**
- Test folder exists and contains test files

**Language-Specific:**
- Node.js: `test/`, `tests/`, `__tests__/`, or `*.test.js` files
- Go: `*_test.go` files
- Rust: `tests/` or inline tests
- Python: `tests/`, `test/`, or `test_*.py` files
- Java: `src/test/java/`

**Pass Criteria:**
- Test folder exists OR test files exist
- At least one test file present

**Auto-fix:**
- Creates test folder
- Optionally creates example test file

---

#### FR-3.2.3: Documentation Folder Recommended
**Priority:** P1  
**Rule ID:** `DOCS_FOLDER_RECOMMENDED`  
**Default Severity:** Info

**Check:**
- `docs/` or `documentation/` folder exists

**Additional Checks:**
- Warns if folder exists but empty
- Suggests structure (API, guides, architecture)

**Auto-fix:**
- Creates `docs/` folder
- Creates subdirectories: `api/`, `guides/`, `architecture/`
- Creates `docs/README.md` index

---

### FR-3.3: Documentation Rules

#### FR-3.3.1: ADR Recommended
**Priority:** P1  
**Rule ID:** `ADR_RECOMMENDED`  
**Default Severity:** Info

**Check:**
- `docs/adr/` or `docs/architecture/decisions/` exists
- Contains at least one ADR file

**Pass Criteria:**
- ADR folder exists
- Contains files matching pattern `NNNN-*.md`

**Auto-fix:**
- Creates `docs/adr/` folder
- Creates `0001-record-architecture-decisions.md` (first ADR about ADRs)
- Creates `template.md` for future ADRs

---

#### FR-3.3.2: API Documentation Required (for libraries)
**Priority:** P1  
**Rule ID:** `API_DOCS_REQUIRED`  
**Default Severity:** Warning

**Applicability:**
- Only for library projects (detected from project structure)

**Check:**
- API documentation exists in one of:
  - `docs/api/`
  - Generated docs (Go: pkg.go.dev link, Rust: docs.rs, etc.)
  - README API section

**Pass Criteria:**
- Documentation describes public API
- Examples included

**Auto-fix:**
- Creates `docs/api/README.md` template
- Suggests documentation generators per language

---

### FR-3.4: CI/CD Rules

#### FR-3.4.1: CI Configuration Recommended
**Priority:** P1  
**Rule ID:** `CI_CONFIG_RECOMMENDED`  
**Default Severity:** Info

**Check:**
- CI/CD configuration exists:
  - GitHub Actions: `.github/workflows/*.yml`
  - GitLab CI: `.gitlab-ci.yml`
  - Jenkins: `Jenkinsfile`
  - CircleCI: `.circleci/config.yml`
  - Travis: `.travis.yml`

**Pass Criteria:**
- At least one CI config file exists
- Config file is valid YAML/Groovy

**Auto-fix:**
- Creates basic CI config for detected platform
- Includes: build, test, lint steps
- Language-specific configuration

---

### FR-3.5: Quality Rules

#### FR-3.5.1: Pre-commit Hooks Recommended
**Priority:** P2  
**Rule ID:** `PRECOMMIT_HOOKS_RECOMMENDED`  
**Default Severity:** Info

**Check:**
- Pre-commit configuration exists:
  - `.pre-commit-config.yaml`
  - `.husky/` folder
  - `.git/hooks/pre-commit` script

**Pass Criteria:**
- Configuration file exists
- Hooks are set up

**Auto-fix:**
- Creates `.pre-commit-config.yaml` with sensible defaults
- For Node.js: sets up husky and lint-staged
- Suggests running installation command

---

#### FR-3.5.2: EditorConfig Recommended
**Priority:** P2  
**Rule ID:** `EDITORCONFIG_RECOMMENDED`  
**Default Severity:** Info

**Check:**
- `.editorconfig` file exists

**Pass Criteria:**
- File exists
- Contains basic settings (charset, indent, line endings)

**Auto-fix:**
- Creates `.editorconfig` with standard settings:
  - UTF-8 charset
  - Unix line endings
  - Indent style based on detected language
  - Trailing whitespace trim

---

### FR-3.6: Language-Specific Rules

#### FR-3.6.1: Node.js - package.json Valid
**Priority:** P0  
**Rule ID:** `NODEJS_PACKAGE_JSON_VALID`  
**Default Severity:** Error

**Check:**
- `package.json` is valid JSON
- Required fields present: name, version
- Scripts section includes: test, build (if applicable)

**Validation:**
- JSON syntax correct
- Version follows semver
- Dependencies versions valid
- No deprecated packages in dependencies

**Auto-fix:**
- Fixes JSON syntax errors (if simple)
- Adds missing required fields
- Suggests adding common scripts

---

#### FR-3.6.2: Go - go.mod Valid
**Priority:** P0  
**Rule ID:** `GO_MOD_VALID`  
**Default Severity:** Error

**Check:**
- `go.mod` is syntactically correct
- Module name is valid
- Go version specified
- `go.sum` exists and is in sync

**Validation:**
- Module name follows Go conventions
- Go version is supported
- No replace directives pointing to missing paths

**Auto-fix:**
- Runs `go mod tidy` to fix go.sum
- Updates Go version if too old

---

#### FR-3.6.3: Rust - Cargo.toml Valid
**Priority:** P0  
**Rule ID:** `RUST_CARGO_TOML_VALID`  
**Default Severity:** Error

**Check:**
- `Cargo.toml` is valid TOML
- Package section complete: name, version, edition
- Dependencies versions valid

**Validation:**
- TOML syntax correct
- Version follows semver
- Edition is valid (2015, 2018, 2021)

**Auto-fix:**
- Fixes simple TOML syntax errors
- Suggests running `cargo update`

---

## 4. Auto-Fix System

### FR-4.1: Fix Safety

#### FR-4.1.1: Confirmation Required
**Priority:** P0  
**Category:** Safety

**Requirements:**
- Never modify files without user confirmation
- Exception: CI mode with `--fix-all --no-interactive`
- Show clear diff before applying changes
- Allow selective fixing (choose which fixes to apply)

**Acceptance Criteria:**
- [ ] Interactive mode is default
- [ ] CI detection works (`CI` env var)
- [ ] Can cancel at any time
- [ ] No partial applies (atomic operations)

---

#### FR-4.1.2: Backup Creation
**Priority:** P1  
**Category:** Safety

**Requirements:**
- Optional backup before modifying files
- Flag: `--create-backups`
- Backup naming: `filename.backup.YYYYMMDD-HHMMSS`
- Cleanup old backups (keep last 5)

**Acceptance Criteria:**
- [ ] Backups created before modifications
- [ ] Backup location configurable
- [ ] Cleanup works correctly
- [ ] Can restore from backup

---

#### FR-4.1.3: Dry Run Mode
**Priority:** P0  
**Category:** Safety

**Requirements:**
- Flag: `--dry-run`
- Shows what would be done
- Shows file diffs
- No actual modifications made
- Can be combined with other flags

**Example Output:**
```
DRY RUN MODE - No changes will be made

Would create: ./README.md
Content preview:
┌────────────────────────────────────
│ # My Project
│ 
│ Description of my project
│ 
│ ## Installation
│ ...
└────────────────────────────────────

Would create: ./tests/
Would create: ./tests/example.test.js

Summary: 2 files, 1 folder would be created
Run without --dry-run to apply changes
```

**Acceptance Criteria:**
- [ ] No files modified in dry-run mode
- [ ] Output shows complete changes
- [ ] Can redirect output to file
- [ ] Exit code same as actual run would be

---

### FR-4.2: Template System

#### FR-4.2.1: Built-in Templates
**Priority:** P0  
**Category:** Templates

**Requirements:**
- Templates for common files:
  - README.md (per language)
  - LICENSE (MIT, Apache-2.0, GPL-3.0, etc.)
  - .gitignore (per language)
  - CONTRIBUTING.md
  - CODE_OF_CONDUCT.md
- Templates embedded in binary
- Variables substitution (project name, author, etc.)

**Template Variables:**
```
{{PROJECT_NAME}}     - from package.json, go.mod, etc.
{{PROJECT_DESC}}     - from package.json or prompt
{{AUTHOR_NAME}}      - from git config or prompt
{{AUTHOR_EMAIL}}     - from git config or prompt
{{LICENSE_TYPE}}     - from choice or config
{{YEAR}}             - current year
{{DATE}}             - current date
```

**Acceptance Criteria:**
- [ ] All variables substituted correctly
- [ ] Templates render properly
- [ ] Missing variables prompted for
- [ ] Can preview before creating

---

#### FR-4.2.2: Custom Templates
**Priority:** P2  
**Category:** Templates

**Requirements:**
- Users can provide custom templates
- Template location: `~/.config/psx/templates/`
- Custom templates override built-in
- Mustache or Go template syntax

**Config:**
```yaml
templates:
  readme:
    path: "~/.config/psx/templates/readme.md"
    variables:
      author: "Mahdi Mohamadi"
      email: "bitsgenix@gmail.com"
```

**Acceptance Criteria:**
- [ ] Custom templates loaded correctly
- [ ] Override mechanism works
- [ ] Template syntax validated
- [ ] Clear error if template invalid

---

### FR-4.3: Fix Operations

#### FR-4.3.1: Create Missing Files
**Priority:** P0  
**Category:** Operations

**Supported Files:**
- README.md
- LICENSE
- .gitignore
- CONTRIBUTING.md
- CODE_OF_CONDUCT.md
- CHANGELOG.md
- .editorconfig

**Process:**
1. Check if file already exists
2. Select appropriate template
3. Gather variables (from config, git, or prompt)
4. Render template
5. Show preview
6. Create file if confirmed

**Acceptance Criteria:**
- [ ] Never overwrites existing files
- [ ] Templates appropriate for detected language
- [ ] Variables filled correctly
- [ ] File permissions set correctly (644)

---

#### FR-4.3.2: Create Missing Folders
**Priority:** P0  
**Category:** Operations

**Supported Folders:**
- src/ (or language equivalent)
- tests/ (or language equivalent)
- docs/
- docs/adr/
- .github/workflows/

**Process:**
1. Check if folder already exists
2. Create folder with proper permissions
3. Optionally create `.gitkeep` to track empty folder
4. Create index/readme in folder

**Acceptance Criteria:**
- [ ] Never removes existing folders
- [ ] Creates parent directories if needed
- [ ] Permissions set correctly (755)
- [ ] Can create multiple folders in one operation

---

#### FR-4.3.3: Fix File Content
**Priority:** P1  
**Category:** Operations

**Supported Fixes:**
- Add missing sections to README
- Update .gitignore with missing patterns
- Fix simple JSON syntax errors
- Add missing fields to package.json

**Process:**
1. Parse existing file
2. Identify what's missing
3. Generate additions
4. Show diff
5. Apply if confirmed

**Acceptance Criteria:**
- [ ] Preserves existing content
- [ ] Diff shows exactly what changes
- [ ] Handles merge conflicts gracefully
- [ ] Validates result after modification

---

## 5. Reporting System

### FR-5.1: Output Formats

#### FR-5.1.1: Table Output
**Priority:** P0  
**Category:** Reporting

**Requirements:**
- Human-readable terminal output
- Colors for severity (red=error, yellow=warning, blue=info)
- Unicode box drawing characters (with ASCII fallback)
- Summary section at bottom
- Respects terminal width

**Layout:**
```
PSX v1.0.0 - Project Structure Checker
═══════════════════════════════════════════════

Detected: Node.js (TypeScript)
Config: psx.yml ✓
Rules checked: 23

ERRORS (2)
───────────────────────────────────────────────
✗ README_MISSING
  Location: ./
  Message: No README.md file found
  Fix: psx fix --rule README_MISSING

✗ LICENSE_MISSING
  Location: ./
  Message: No LICENSE file found
  
WARNINGS (3)
───────────────────────────────────────────────
⚠ TESTS_FOLDER_MISSING
  Location: ./
  Message: No tests/ folder found
  Fix: psx fix --rule TESTS_FOLDER_MISSING

[... more warnings ...]

═══════════════════════════════════════════════
Summary: 2 errors, 3 warnings, 1 info
Status: FAILED ✗
Time: 0.23s

Fix command: psx fix --interactive
```

**Acceptance Criteria:**
- [ ] Colors work in color terminals
- [ ] Graceful fallback for no-color
- [ ] Respects `NO_COLOR` env var
- [ ] Works in 80-column terminals
- [ ] Links are clickable (if terminal supports)

---

#### FR-5.1.2: JSON Output
**Priority:** P0  
**Category:** Reporting

**Requirements:**
- Machine-readable format
- Schema-validated
- All information from table output
- Can be piped to jq or other tools

**Schema:**
```json
{
  "version": "1.0.0",
  "timestamp": "2024-12-02T10:30:00Z",
  "project": {
    "path": "/path/to/project",
    "type": "nodejs",
    "subtype": "typescript"
  },
  "config": {
    "file": "psx.yml",
    "valid": true
  },
  "results": [
    {
      "rule_id": "README_MISSING",
      "category": "general",
      "severity": "error",
      "message": "No README.md file found",
      "location": "./",
      "fixable": true,
      "fix_command": "psx fix --rule README_MISSING"
    }
  ],
  "summary": {
    "total_rules": 23,
    "passed": 18,
    "errors": 2,
    "warnings": 3,
    "info": 0
  },
  "status": "failed",
  "execution_time_ms": 234
}
```

**Acceptance Criteria:**
- [ ] Valid JSON always (even on error)
- [ ] Schema version included
- [ ] Can be pretty-printed
- [ ] Includes all relevant information

---

#### FR-5.1.3: SARIF Output
**Priority:** P1  
**Category:** Reporting

**Requirements:**
- SARIF 2.1.0 format
- Compatible with major IDEs
- Includes fix suggestions
- Proper location mapping

**Use Cases:**
- VS Code SARIF Viewer extension
- GitHub Code Scanning
- IntelliJ IDEA
- SonarQube import

**Acceptance Criteria:**
- [ ] Valid SARIF 2.1.0
- [ ] Locations map to actual files
- [ ] Fix information included
- [ ] Works with VS Code extension

---

#### FR-5.1.4: HTML Report
**Priority:** P1  
**Category:** Reporting

**Requirements:**
- Self-contained HTML file
- Works offline
- Responsive design
- Filterable/sortable
- Embeds CSS/JS (no external dependencies)

**Features:**
- Expand/collapse sections
- Filter by severity
- Search functionality
- Export timestamp
- Link to documentation

**Acceptance Criteria:**
- [ ] Opens in any browser
- [ ] Works without internet
- [ ] Responsive on mobile
- [ ] Print-friendly

---

### FR-5.2: Progress Indicators

#### FR-5.2.1: Progress Bar
**Priority:** P1  
**Category:** UX

**Requirements:**
- Show progress for operations >2 seconds
- Display: percentage, current/total, ETA
- Smooth updates (not jumpy)
- Clears when complete

**Example:**
```
Checking rules... ████████████░░░░░░░░ 65% (15/23) ETA: 2s
```

**Acceptance Criteria:**
- [ ] Appears after 2 second delay
- [ ] Updates smoothly
- [ ] ETA reasonably accurate
- [ ] Clears completely when done

---

## 6. Configuration System

### FR-6.1: Config File

#### FR-6.1.1: YAML Format
**Priority:** P0  
**Category:** Configuration

**Requirements:**
- Standard YAML syntax
- Schema version for compatibility
- Comments preserved on edit
- Supports includes (for shared configs)

**Validation:**
- Syntax check on load
- Type validation
- Unknown key warnings
- Required field checking

**Acceptance Criteria:**
- [ ] Valid YAML parsed correctly
- [ ] Clear error for invalid YAML
- [ ] Line numbers in errors
- [ ] Comments don't break parsing

---

#### FR-6.1.2: Config Discovery
**Priority:** P0  
**Category:** Configuration

**Search Order:**
1. `--config <file>` flag
2. `psx.yml` in current directory
3. `.psx.yml` in current directory
4. Walk up to git root
5. `~/.config/psx/psx.yml` (user default)
6. Built-in defaults

**Acceptance Criteria:**
- [ ] Finds config in all locations
- [ ] Respects priority order
- [ ] Shows which config is being used
- [ ] Works without config (uses defaults)

---

#### FR-6.1.3: Config Validation
**Priority:** P0  
**Category:** Configuration

**Validation Checks:**
- Schema version supported
- All rule IDs valid
- Severity values correct
- Paths exist (for templates)
- Regular expressions valid

**Command:**
```bash
psx config validate
psx config validate --config custom.yml
```

**Acceptance Criteria:**
- [ ] Catches all common errors
- [ ] Suggests fixes for errors
- [ ] Exit code 0 if valid, 1 if not
- [ ] Can validate without running check

---

### FR-6.2: Rule Configuration

#### FR-6.2.1: Enable/Disable Rules
**Priority:** P0  
**Category:** Configuration

**Syntax:**
```yaml
rules:
  general:
    readme_required: error    # enabled as error
    license_required: warning # enabled as warning
    changelog: false          # disabled
```

**Acceptance Criteria:**
- [ ] Rules can be turned off
- [ ] Severity can be changed
- [ ] Unknown rules warned about
- [ ] Category shortcuts work

---

#### FR-6.2.2: Custom Rules
**Priority:** P1  
**Category:** Configuration

**Syntax:**
```yaml
custom_rules:
  - name: "TERRAFORM_PRESENT"
    description: "Check for Terraform files"
    severity: info
    check: "file_exists:terraform/main.tf"
  
  - name: "DOCKER_PRESENT"
    description: "Check for Dockerfile"
    severity: warning
    check: "file_exists:Dockerfile"
```

**Check Types:**
- `file_exists:<path>`: File exists
- `folder_exists:<path>`: Folder exists
- `file_contains:<path>:<pattern>`: File contains text/regex
- `folder_not_empty:<path>`: Folder has files
- `shell:<command>`: Run shell command (exit 0 = pass)

**Acceptance Criteria:**
- [ ] All check types work
- [ ] Regex patterns supported
- [ ] Shell commands sandboxed
- [ ] Custom rules in reports

---

## 7. CLI Interface

### FR-7.1: Commands

#### FR-7.1.1: check Command
**Priority:** P0  
**Category:** CLI

**Usage:**
```bash
psx check [options]
```

**Options:**
- `--verbose, -v`: More detailed output
- `--level <error|warning|info>`: Filter by severity
- `--output <format>`: Output format
- `--config <file>`: Config file path
- `--fail-on <error|warning>`: When to exit 1

**Examples:**
```bash
psx check                          # basic check
psx check --verbose                # detailed
psx check --level error            # only errors
psx check --output json > out.json # json to file
psx check --fail-on warning        # fail on warnings too
```

**Acceptance Criteria:**
- [ ] All options work
- [ ] Help text clear
- [ ] Exit codes correct
- [ ] Errors to stderr, results to stdout

---

#### FR-7.1.2: fix Command
**Priority:** P0  
**Category:** CLI

**Usage:**
```bash
psx fix [options]
```

**Options:**
- `--interactive, -i`: Ask before each fix
- `--dry-run`: Show changes without applying
- `--rule <id>`: Fix specific rule
- `--all`: Fix all fixable issues
- `--create-backups`: Create backups before fixing

**Examples:**
```bash
psx fix                            # fix interactively
psx fix --dry-run                  # preview
psx fix --rule README_MISSING      # fix one rule
psx fix --all --no-interactive     # fix everything (CI mode)
```

**Acceptance Criteria:**
- [ ] Interactive mode is default
- [ ] Dry-run makes no changes
- [ ] Can select individual fixes
- [ ] Creates backups if requested

---

#### FR-7.1.3: init Command
**Priority:** P0  
**Category:** CLI

**Usage:**
```bash
psx init [options]
```

**Options:**
- `--template <lang>`: Use language template
- `--minimal`: Bare minimum config
- `--force`: Overwrite existing config

**Examples:**
```bash
psx init                    # interactive
psx init --template nodejs  # nodejs template
psx init --minimal          # minimal config
```

**Process:**
1. Detect project type (if not specified)
2. Ask configuration questions
3. Generate psx.yml
4. Show what was created

**Acceptance Criteria:**
- [ ] Creates valid config
- [ ] Doesn't overwrite without --force
- [ ] Templates work correctly
- [ ] Minimal config is truly minimal

---

#### FR-7.1.4: rules Command
**Priority:** P1  
**Category:** CLI

**Usage:**
```bash
psx rules [options]
```

**Options:**
- `--category <name>`: Filter by category
- `--json`: Output as JSON
- `--verbose`: Show full descriptions

**Output:**
```
Available Rules:

General (4 rules)
─────────────────
  readme_required        Checks for README file
  license_required       Checks for LICENSE file
  gitignore_required     Checks for .gitignore
  changelog_recommended  Recommends CHANGELOG.md

Structure (3 rules)
───────────────────
  src_folder_required         Checks for source folder
  tests_folder_required       Checks for tests folder
  docs_folder_recommended     Recommends docs folder

... more categories ...

Total: 23 rules
```

**Acceptance Criteria:**
- [ ] Shows all available rules
- [ ] Can filter by category
- [ ] JSON output valid
- [ ] Shows default severities

---

### FR-7.2: Global Flags

#### FR-7.2.1: Common Flags
**Priority:** P0  
**Category:** CLI

**Flags Available for All Commands:**
- `--help, -h`: Show help
- `--version`: Show version
- `--verbose, -v`: Verbose output
- `--quiet, -q`: Minimal output
- `--no-color`: Disable colors
- `--config <file>`: Config file

**Acceptance Criteria:**
- [ ] Work with all commands
- [ ] Help shows all options
- [ ] Flags can be combined
- [ ] Short and long forms work

---

#### FR-7.2.2: Environment Variables
**Priority:** P1  
**Category:** CLI

**Supported Variables:**
- `PSX_CONFIG`: Config file path
- `PSX_NO_COLOR`: Disable colors (any value)
- `PSX_LOG_LEVEL`: Log verbosity (debug, info, warn, error)
- `CI`: Auto-detect CI environment

**Precedence:**
1. CLI flags
2. Environment variables
3. Config file
4. Defaults

**Acceptance Criteria:**
- [ ] All variables work
- [ ] Precedence correct
- [ ] Can disable with empty value
- [ ] Documented in help

---

## Testing Requirements

### Test Coverage
- Unit tests: >80% coverage
- Integration tests: All major workflows
- E2E tests: One per supported platform
- Performance tests: Large project handling

### Test Scenarios
Each functional requirement must have:
- Happy path test
- Error case test
- Edge case tests
- Performance test (if applicable)

---

## Acceptance Checklist

Before v1.0 release, ALL P0 and P1 requirements must be:
- [ ] Implemented
- [ ] Tested
- [ ] Documented
- [ ] Reviewed

---

**Document Version:** 1.0.0  
**Status:** Draft  
**Next Review:** After initial implementation

*This is a living document and will be updated as requirements evolve.*
```
