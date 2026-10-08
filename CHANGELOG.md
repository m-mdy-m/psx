# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.2] - 2026-10-08

### Fixed

- **Unix installer cleanup failed under set -u.** The temporary-directory EXIT trap could reference a local variable after its scope ended, producing `tmp_dir: unbound variable` after a successful installation.
- **PowerShell installer failed to parse.** Error handling used a variable immediately followed by a colon, which PowerShell interpreted as an invalid variable reference. The message now uses explicit formatting.
- **PSX dogfooding CI failed because psx.yml was missing.** Added a repository configuration used by the CI self-check job.
- **Go module verification failed in CI.** go.mod and go.sum are now synchronized with the actual doublestar dependency usage and current module graph.

### Changed

- Installer cleanup is now safe with strict Bash mode enabled.
- PowerShell installer error messages are compatible with both direct script execution and irm ... | iex.
- The repository now maintains an explicit PSX configuration for its own CI structure check.

## [0.3.1] - 2026-10-08

### Fixed

* **Linux installation failed with `Not Found`.** The Unix installer used outdated asset names such as `psx-linux-amd64`, while releases publish platform assets using the `x64` naming scheme.
* **Windows installation failed to download the binary.** The PowerShell installer requested `psx-windows-amd64.exe`, while the published release asset is `psx-windows-x64.exe`.
* **Installers could install invalid downloads as executables.** The Unix installer could save a GitHub `404 Not Found` response as `psx` and attempt to execute it.
* **Improved download validation.** Installers now verify that the downloaded release artifact is valid before installing it.
* **Release archive handling fixed.** Linux and macOS installers now download and extract the corresponding release archives instead of expecting a raw binary with an outdated filename.
* **Checksum verification improved.** Installers validate the downloaded binary against the release checksum information when available.
* **PowerShell installer no longer terminates the host session on normal installation errors.** This makes usage through `irm ... | iex` safer and more predictable.
* **PATH handling improved.** User-level installations correctly update the user's PATH and current PowerShell session where possible.
* **Build and release naming consistency.** Build scripts now use the same platform naming convention as the published GitHub release assets.

### Changed

* Standardized release platform names to:

  * `linux-x64`
  * `linux-arm64`
  * `darwin-x64`
  * `darwin-arm64`
  * `windows-x64`
* Updated Unix and Windows installation scripts to match the actual GitHub release asset layout.
* Improved installer error messages so failed downloads are reported clearly instead of appearing as successful installations.

### Installation

Install the latest release with:

```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

Windows:

```powershell
irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex
```

[0.3.2]: https://github.com/m-mdy-m/psx/releases/tag/v0.3.2
[0.3.1]: https://github.com/m-mdy-m/psx/releases/tag/v0.3.1
[0.3.0]: https://github.com/m-mdy-m/psx/releases/tag/v0.3.0


## [3.0.0] - Unreleased

Major release. The CLI, the rule system and the internals were reworked. Every item below
was reproduced before it was fixed; see [docs/VERIFICATION.md](docs/VERIFICATION.md).

### Fixed

Correctness bugs that produced wrong results rather than errors:

- **`**` globs matched nothing.** `filepath.Glob` does not support `**`, so
  `tests_folder` failed on every Go project regardless of its tests — a Go project with
  `internal/rules/engine_test.go` still reported no tests. Globbing now uses `doublestar`
  and handles brace expansion.
- **`project.type` was normalized then discarded.** Aliases such as `ts`, `js` and `golang`
  never reached the pattern lookup, so no patterns matched and language-specific rules
  silently reported success on projects that violated them.
- **A rule with no applicable patterns was reported as a pass.** Results are now tri-state:
  passed, failed, or skipped. Skips are counted separately and shown with `--show-skipped`.
- **`check` was interactive and wrote to the project.** It read `fix`'s flags, and
  `--interactive` defaulted to true, so every `psx check` opened a prompt form and wrote
  `.psx-project.yml`. Prompts now require a real terminal, and read-only commands never
  write.
- **`psx help` printed nothing.** The help text was fetched and discarded. Help is now
  produced by cobra.
- **JSON output was not parseable.** Diagnostics were printed to stdout, so
  `psx check -o json | jq` failed. All diagnostics moved to stderr.
- **`psx.rar` placeholder text in generated files.** `scripts/test.sh` contained a literal
  `{{test_command}}` and `clean.sh` contained `rm -rf {{build_dirs}}` — scripts that could
  not run. Command variables now come from the language profile.
- **`FUNDING.yml` was written as prose.** GitHub reads that file as structured data keyed
  by platform, so a markdown list parses as YAML but is silently ignored.
- **`ignore` was validated but never applied.** Vendored trees satisfied `tests_folder`.
  Ignored directories are now pruned during the scan.
- **`custom` paths could escape the project.** `path: ../../.bashrc` in a cloned repo's
  config would write outside the project on `psx fix --yes`. Paths are now resolved against
  the project root and rejected if they escape it, symlinks resolved.
- **Generated scripts were not executable.** Created 0644.
- **Empty files and directories failed their own rule.** `fix` filled an existing empty file
  or directory instead of leaving the rule unsatisfied.
- **`--config` was ignored.** The flag was registered on the root command and threaded into
  the options struct, but `cmdctx.Load` passed an empty path to the loader, so every
  invocation used whichever config it discovered instead.
- **Project detection never ran.** `NormalizeProjectType("")` returns `generic`, and that
  normalised value was written back onto the config. The command layer then saw a
  non-empty declared type, concluded the user had asked for `generic`, and skipped
  detection — so any project with an unset `project.type` was reported as `generic` no
  matter what manifests it had. `psx detect` and `psx check` disagreed as a result. The
  raw value is now kept separately from the normalised one.
- **Global flags were read before they were parsed.** `baseOptions()` ran while the
  command tree was being built, so it captured the flag defaults rather than the user's
  values. Persistent flags are now applied inside `RunE`.
- **Release binaries always reported "development".** The workflow passed
  `-X 'main.version=…'` but the variable lives in `internal/command`. The build also had a
  `|| go build` fallback that hid compile errors behind an unstripped, unversioned binary.
- **Release and Docker jobs ran on branch pushes.** `${GITHUB_REF#refs/tags/v}` was
  meaningless outside a tag build. Both are now tag-gated.
- **Non-deterministic output.** Template selection iterated a map, so a project without a
  `generic` tier could get different content each run.
- **`--baseline` could not be bootstrapped.** It returned an error when the file did not
  exist, but creating that file is the only way to adopt a baseline on an existing
  repository. A missing file is now written from the current run, with a header explaining
  how to shrink it.
- **`psx watch --baseline` was silently ignored.** The flag was registered but never read,
  so known issues were announced as fresh failures on every scan. Each iteration now filters
  against the file, without creating one as a side effect of watching.
- **`--level` filtered the report but not the summary.** `--level warning` printed an empty
  report above a line reading `Result: 2 info`. Counts now come from what was actually
  displayed.
- **Compact mode under-counted.** It printed every finding but only summarised errors and
  warnings, so a run with 12 failures read `12 failed (0 errors, 1 warnings)`. Both formats
  now share one severity counter, and the parts always add up to the total.
- **`psx fix` printed absolute paths.** Created files were listed by full path, which is
  noisy and unreadable. Paths are now relative to the project root.
- **`custom.folders` created every entry as a directory.** `production.yaml: {}` — an empty
  file — produced a directory with that name, so the documented structure was unusable. A
  key with children is a directory; a key with none is a file, and a string leaf has its
  content written.
- **`psx fix` counted directories as files.** The closing line read `Created 9 file(s)`
  when four of the nine were directories.
- **`psx detect` crashed on any project it could not recognise.** With no recognised
  manifest the type falls back to `generic`, whose profile declares an empty package
  manager list; reading a "preferred" manager out of that empty slice panicked. Every
  non-Go, Node, Rust or Python project crashed with a stack trace.
- **Detection picked the wrong language when two manifests were present.** The markers
  claimed to be ordered by specificity but were not: `package.json` (0.95) sat below
  `requirements.txt` (0.70), so a project with both was reported as Python.
- **A directory named `go.mod` was treated as a module manifest.** The path index mixed
  files and directories, so any manifest check could match a directory.
- **An npm project was offered a pnpm lockfile hint.** Every Node package manager shares
  `package.json` and differs only by its lockfile, so filtering on the manifest could not
  discriminate. Lockfiles are now matched first.
- **The watcher's ignore list was ignored.** It accepted the configured patterns and
  replaced them with a hardcoded `.git`/`node_modules`/`vendor` list, so a rebuild writing
  to an ignored directory re-triggered a full re-check. It now shares `tree.Scan` with the
  check itself, and the tree is walked once per change rather than twice.
- **`watch` compared against a nil result.** `Compare(nil, nil)` panicked; it guarded the
  previous result but not the current one.
- **`watch` panicked on a zero interval.** `time.NewTicker` panics on a non-positive
  interval, so an unpopulated `Options` crashed instead of reporting a problem.
- **`tree.Scan` returned an empty snapshot for a directory that did not exist**, and for a
  file passed as root, indexing the root entry as a file with its `IsDir` flag cleared.
- **`flags` carried three fields nothing could reach.** `Fix.Watch` was never set by a flag,
  `Fix.CreateBackups` had no `--backup` flag behind it, and `Watch.Format` was never read.
  Two `Defaults()` also disagreed on the watch timings and one silently overwrote the other.
- **`utils` declared a second, wrong exit-code table.** It assigned `ExitArgs = 4` where
  the command package uses 3, and nothing referenced it. Removed rather than corrected,
  since a second table that nobody reads is how the real one drifts.

### Added

- **Examples.** Thirteen runnable configurations in `examples/`, from `minimal.yml` through
  per-profile files to a `strict.yml` for gating a build. Each is also what `psx init`
  produces for that profile. Validated by `TestExamplesAreValidConfigs`, so none can rot.
- **`tree` package.** One `WalkDir` builds an immutable index; every rule then matches in
  memory instead of issuing its own `stat` or glob. Literal prefixes narrow the candidate
  set.
- **Declarative fixes.** Each rule declares its own fix in YAML, including a literal path,
  file mode, safety flag and template name. Adding a rule no longer requires writing Go,
  and all per-rule `switch` statements are gone.
- **Tri-state rule results.** `skipped` is distinct from `passed` and is visible in the
  summary.
- **Project detection.** Language and archetype are inferred from marker files
  (`go.mod`, `package.json`, `pnpm-workspace.yaml`, `go.work`, `services/`, …) with a
  confidence figure and the signals used. A declared `project.type` always wins.
  `psx init` uses it to pick a sensible rule set.
- **`psx init`** writes a `psx.yml` tuned to the detected profile.
- **`psx rules`, `psx explain <rule>`** — inspect the catalogue.
- **`psx detect`** — show the inferred language, kind, workspace globs and package managers.
- **`psx watch`** — re-check on change, reporting only the difference. `--fix` applies rules
  marked safe; `--once-clean` exits as soon as the project is clean.
- **21 GitHub Actions templates** across five groups — CI, release, Docker, deploy and
  security. Each is standalone. Includes the digest-push multi-arch Docker pattern and a
  release orchestrator that calls per-platform `workflow_call` workers.
- **Five workflow rules** (`workflow_ci`, `workflow_docker`, `workflow_release`,
  `workflow_codeql`, `workflow_secret_scan`) that detect and create those files.
- **`psx workflows`** to list and preview them.
- **Eight output formats**: `table`, `compact`, `json`, `ndjson`, `sarif`, `github`,
  `junit`, `markdown`.
- **`--baseline`** forgives already-existing violations so psx can be adopted without
  failing the build on day one.
- **Rule filters**: `--only`, `--category`, `--level`, `--show-skipped`, `--all`.
- **`--force`** to overwrite existing files; **`--yes`** to run unattended.
- **New rules**: `gitattributes`, `env_example`, `lockfile`, `funding`, `support`,
  `roadmap`, `architecture`, `runbook`, `openapi`, `makefile`, `kubernetes`, `nginx`.
  43 rules in total, 37 auto-fixable.
- **Rust and Python** language profiles, alongside Go and Node.js.
- **New templates**: gitattributes, Makefile, `.env.example`, nginx, Kubernetes, Helm,
  funding, support, roadmap, architecture, runbook, OpenAPI, release workflow, Dependabot
  and Renovate.
- **Tests.** Coverage across `tree`, `rules`, `resources`, `report` and `config`. The
  notable ones: `TestEveryTemplateRenders` (no template may reference an unsupplied
  variable), `TestGeneratedYAMLIsValid` (every generated YAML file must parse),
  `TestFixIsIdempotent`, `TestWorkflowTemplatesAreValidYAML`.
- **CI for psx itself**, with a job that runs psx against its own repository.
- **`.gitattributes`** and **`dependabot.yml`**.
- **Documentation**: `docs/RULES.md` (generated from `rules.yml`), `CONFIGURATION.md`,
  `WORKFLOWS.md`, `ARCHITECTURE.md`, and a rewritten `CONTRIBUTING.md` and SRS.

### Changed

- **Command surface.** `check`, `fix`, `watch`, `init`, `rules`, `explain`, `workflows`,
  `detect`. Flags now live on the command they affect: `fix` no longer accepts
  `--output`, `--baseline` or `--fail-on`, which did nothing there.
- **Fixes are literal.** `fix.path` must be a literal path. A glob describes something to
  detect, not something safe to create.
- **Fixes never overwrite silently.** A file with content is left alone unless `--force`.
  Running `fix` twice always reports nothing the second time.
- **Prompting requires a terminal.** In CI, in a pipe, or with `PSX_NON_INTERACTIVE`,
  prompts return defaults instead of blocking on stdin.
- **Diagnostics go to stderr**, always. stdout carries only the report.
- **One template resolution path.** `resources.Template` covers documents, container files,
  scripts and workflows alike.
- **Deterministic output.** Rules are evaluated in sorted order; templates are selected in
  sorted order.

### Removed

- **Global mutable flag state.** `flags.GetFlags()` returned the address of the defaults
  struct, which is how `check` ended up reading `fix` flags. Each command owns its options.
- **The `Rule` interface and `registry.go`.** Removed in 2.0.0 but still documented in
  `CONTRIBUTING.md`; rules are data now.
- **`pattern_resolver.go`, `content_gen.go`'s per-rule switches, `custom_rule.go`'s
  unsanitised paths, and the `reporter` package** — replaced by `tree`, `resources`, the
  declarative fix path and `report`.
- **Dead configuration.** `additional_checks` had no implementation; `--level` was parsed
  and ignored; `fix.backup` and `fix.interactive` were read and unused.

### Migration

- `project.type` is optional again: set it to `auto` or leave it empty to detect from the
  layout. An explicit value still wins.
- `ignore` now applies. Projects that relied on it being inert may see new failures in
  vendored directories; add those paths to `ignore`.
- `fix` no longer accepts `--output`, `--level`, `--baseline`, `--fail-on` or `--all`.
  Use `check` to report.
- Removed rules from 2.0.0 that never shipped as working rules (Renovate, Kubernetes, Nginx
  as rules) are replaced by the workflow template library. `psx rules` lists the current
  43.
- CI consumers: `-o json` output is now a single valid document with passing rules omitted
  unless `--all`. Use `--all` if you relied on the old shape.

## [2.0.0] - 2025-12-19

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

[3.0.0]: https://github.com/m-mdy-m/psx/compare/v2.0.0...HEAD
