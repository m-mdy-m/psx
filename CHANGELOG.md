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
