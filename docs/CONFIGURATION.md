# Configuration

psx reads `psx.yml` (also accepted: `.psx.yml`, `psx.yaml`, `.psx.yaml`).

Lookup order:

1. the path given to `--config`
2. the four names above, in the project directory
3. the same four names at the git root, walking up at most 10 levels
4. `~/.config/psx/`
5. built-in defaults

With no configuration file, every rule is enabled at its default severity.

Generate a starting point with `psx init`, which infers the project type and writes a
sensible subset:

```bash
psx init
psx init --type nodejs --kind library
psx init --minimal
```

## Full schema

```yaml
version: 1

project:
  type: go            # nodejs | go | rust | python | generic
  kind: cli           # app | library | cli | monorepo | microservice | plugin

rules:
  readme: error       # error | warning | info | false
  license: false      # false disables a rule
  tests_folder: error

ignore:
  - node_modules/
  - vendor/
  - "**/*.snap"

fix:
  interactive: true
  backup: false

custom:
  files:
    - path: .env.example
      content: |
        APP_ENV=development
  folders:
    - path: config
      structure:
        base:
          production.yaml: {}
```

## project.type

Selects the language tier. Aliases resolve automatically: `ts`, `js`, `node` and
`typescript` all mean `nodejs`; `golang` means `go`.

Leave it empty (or set `auto`) to detect from the layout — `go.mod`, `package.json`,
`Cargo.toml`, `pyproject.toml`. A declared value always wins. Check what psx inferred
with `psx detect`.

Rules match patterns per type. A rule that only declares, say, Go patterns is reported as
**skipped** for a Python project rather than silently passing — a skip is visible in the
summary and in `--show-skipped`.

## rules

`rules` is an explicit allow-list. Listing only some rules means the others do not run.

```yaml
rules:
  readme: error
  license: warning
  adr: false      # recognised, but disabled
```

An unknown rule id is reported as a warning and skipped, so a typo degrades instead of
breaking the run.

## ignore

Patterns use gitignore syntax: trailing `/` matches directories, a leading `/` anchors to
the project root, and a leading `!` re-includes. Ignored directories are pruned during the
scan, so `node_modules` is never walked.

```yaml
ignore:
  - dist/               # any directory named dist
  - /build              # only at the project root
  - "**/*.snap"         # nested snapshots
  - "logs/*"            # contents of logs
  - "!logs/keep.txt"    # but keep this one
```

## fix

```yaml
fix:
  interactive: true   # confirm each change; implied false when there is no terminal
  backup: false       # copy the original to .psx/backup before overwriting
```

## custom

Anything psx has no rule for. Both keys are applied after the built-in rules.

`custom.files` creates files with exactly the content you give.

`custom.folders.structure` nests arbitrarily deep:

```yaml
custom:
  folders:
    - path: src
      structure:
        components:
          Button.tsx: {}
          index.ts: ""
        styles:
          tokens.css: ""
```

Paths are resolved against the project root and rejected if they escape it. A `custom`
entry of `../../.bashrc` fails rather than writing outside the project — a configuration
file in a cloned repository is untrusted input.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Passed, or failed below the `--fail-on` threshold |
| 1 | Validation failed |
| 2 | Configuration problem |
| 3 | Invalid arguments |

## Overriding on the command line

```bash
psx check --only readme,license   # a subset of rules
psx check --category documentation
psx check --level warning         # hide info-level findings
psx check --fail-on none          # never fail
psx check --show-skipped          # include rules that did not apply
psx check --all                   # include passing rules
```

## Adopting gradually

`--baseline` forgives problems that already exist, so a project can turn psx on without
failing its build on day one. New problems still fail.

```bash
psx check --baseline .psx-baseline.txt
```

Record the current state, then keep the file in version control. Any format works:

```yaml
# .psx-baseline.txt
readme
license
```

```json
["readme", "license"]
```

```json
[{ "rule_id": "readme", "severity": "error" }]
```

## Project metadata

Templates that mention the project name, author or repository read those from git and, on
first `fix`, from a short prompt. The result is cached in `.psx-project.yml`.

`check` never prompts and never writes this file.