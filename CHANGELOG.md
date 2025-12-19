# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-17

### Added

#### Core Features
- Project structure validation with configurable rules
- Auto-fix capability for common structural issues
- Multi-language project detection (Go, Node.js)
- Interactive and non-interactive modes for fixing issues
- Comprehensive configuration system via YAML files

#### Rules Engine
- 47 built-in validation rules across multiple categories:
  - General: README, LICENSE, .gitignore, CHANGELOG
  - Structure: src/, tests/, docs/, scripts/ folders
  - Documentation: ADR, CONTRIBUTING, API docs, SECURITY
  - CI/CD: GitHub Actions, Renovate, Dependabot
  - Quality: EditorConfig, pre-commit, Prettier, ESLint, Husky
  - DevOps: Docker, Kubernetes, Nginx configurations

#### Auto-Fix Capabilities
- Create missing files (README, LICENSE, etc.)
- Generate language-specific configurations
- Set up CI/CD workflows
- Configure code quality tools
- Create project documentation structure

#### Project Templates
- README templates for different languages
- Multiple LICENSE options (MIT, Apache-2.0, GPL-3.0, BSD-3-Clause)
- Language-specific .gitignore templates
- Docker and docker-compose configurations
- Kubernetes deployment templates
- GitHub Actions workflows
- Pre-commit hooks and quality tool configs

#### Installation & Distribution
- Single binary distribution for Linux, macOS, and Windows
- Docker images (standard, Alpine, scratch variants)
- Installation scripts for Unix and Windows
- Makefile for building and releasing

#### Developer Experience
- Verbose mode for detailed output
- Dry-run mode for previewing fixes
- Project information caching
- Configuration validation
- Shell completion support (bash, zsh, fish)

### Technical Details

#### Architecture
- Written in Go 1.25+
- Embedded configuration and templates
- Concurrent rule execution
- Modular rule system with registry pattern

#### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64/Apple Silicon)
- Windows (amd64)

#### Performance
- Fast project scanning (<1s for typical projects)
- Efficient parallel rule execution
- Low memory footprint (<50MB for most projects)

### Documentation
- Comprehensive README with examples
- Software Requirements Specification (SRS)
- Functional and Non-Functional Requirements documents
- Contributing guidelines
- Code of Conduct
- Security policy

[1.0.0]: https://github.com/m-mdy-m/psx/releases/tag/v1.0.0

## [1.0.1] - 2025-12-17

### Fixed
- Docker build and publish issues
- Minor fixes and updates

[1.0.1]: https://github.com/m-mdy-m/psx/releases/tag/v1.0.1

## [2.0.0] - 2025-12-19

### Changed
- Full codebase rewrite and large refactor: many internal modules and types were rewritten to simplify logic, improve maintainability, and reduce complexity.
- Refactored rules engine and loader to support multi-folder and multi-file handling and richer metadata (see [loader.go](./internal/resources/loader.go)). 
- Resources handling and config logic rewritten to better separate sources and configs and to support `languages.yml` metadata. 
- CLI surface simplified: consolidated commands and reduced surface area; improved CLI messages and the command reporter. Output formatting and reporter behavior was rewritten.
- Directory layout reorganized and simplified; `fixer` and `checker` directories and related internal complexity removed in favor of a streamlined structure. 
- Configuration handling rewritten: new robust handling of schema-less configuration and updated validation flows.
- Improved user-facing messages and Makefile updates. 
- Templates and resource placement reorganized: template files were moved into more appropriate locations (for example, ADR templates were moved from `templates.yml` into `docs-templates.yml`) and other template/resource files were relocated for clearer structure and discoverability.

### Removed (Breaking)
- Automatic project type detection removed — the `detector` directory and related auto-detection logic were deleted. PSX will **no longer** infer project type automatically. Users must explicitly set `project.type`.
- `project` CLI command removed — scripts and automation need to use the remaining top-level commands (`check`, `fix`). 
- Old `fixer`, `checker`, `detector`, and registry implementations removed; several legacy types and modules were deleted.
- Previously supported languages such as Python and Rust were removed from the default supported set to focus on **Go** and **Node.js** (so the product now ships narrower language support). 
- Removed multiple non-critical rules and CI/quality configurations, including Husky, commitlint, ESLint, Prettier, lint-staged, Git attributes, Nginx, ReVonk, Dependabot, and other optional rules. 
- Old CI/workflow files and quality-tool configs removed (e.g., Git hooks and some GitHub Actions were removed in favor of ADR/doc changes). 
- Removed many extra/unused files and useless helper functions that were introduced by duplicated flows in the old `fixer` and `checker` directories. Redundant flows were eliminated and the legacy noisy code was deleted to reduce maintenance burden.

### Fixed
- Fixed schema-less configuration validation issues.
- Fixed Docker build & publish issues.
- Multiple minor bug fixes across rules, handlers, resource logic, and config handling.

### Added
- Custom rule and validator for `psx.yml` configuration. 
- Linter integration and configuration: `.golangci.yml` added and repository now includes linting rules.
- `messages.yml` with improved and new user-facing messages — many new message entries and clearer wording were added for all YAML outputs and validations (ADR, API, README, etc.).
- `languages.yml` metadata: language-specific metadata added for Go and Node.js (including js/ts), centralizing language info used by resource handlers.
- Better YAML templates and message files for common artifacts (ADR, API docs, README, and others) — new/updated YML templates and message content were added.
- Consolidated fixer/checker implementation: removed duplicated flows and replaced them with simplified, centralized implementations located at `rules/checker.go` and `rules/fixer.go`.
- Example configuration file added at [psx.examples.yml](./examples/psx.examples.yml) showcasing both standard rules and custom files/folders setup. This provides a reference for users to define `custom.files` and `custom.folders` along with project rules, including Node.js or Go projects.

### Notes / Migration
- **Major / breaking release:** this is a breaking change release — bump to **v2.0.0** is required because of removed/renamed commands and changed behavior.
- **Project type**: Consumers must explicitly set `project.type` in `psx.yml` / `psx.yaml` / `.psx.yml` / `.psx.yaml`. Example:
  ```yaml
  project:
    type: "go"      
````

* **CLI:** Replace any use of `psx project` with `psx check` or `psx fix`. The reporter output and CLI flags have changed — update automation and CI accordingly.
* **Rules:** Rule metadata and rule handling were refactored; review `rules.yml` and any custom rules to ensure they match the new metadata shape and loader behavior.
* **Linter:** Add `golangci-lint` to your local dev flow / CI if you want to catch new lint rules.
* **Docs & ADR:** GitHub Actions removed in docs refactor — check `docs/ADR` for rationale and update CI if you relied on old workflows.
* **If you relied on auto-detection:** migrate to explicit `project.type` in user configs and onboarding docs.

[2.0.0]: https://github.com/m-mdy-m/psx/releases/tag/v2.0.0
