# Examples

Runnable configurations. Copy one into your project as `psx.yml`, or point at one
to try it:

```bash
psx check --config examples/go-cli.yml .
```

Each file is validated by `TestExamplesAreValidConfigs`, so they all load. Several are
stricter than what `psx init` writes — `go-cli.yml` promotes tests and a lockfile to
`error`, where init leaves them out. Copy the closest one and trim it, or start from
`psx init` and tighten it.

| File | For |
| --- | --- |
| [`minimal.yml`](minimal.yml) | The smallest config worth keeping. Start here |
| [`strict.yml`](strict.yml) | Gate a build. Few rules, all `error` |
| [`nodejs-library.yml`](nodejs-library.yml) | Published npm package |
| [`go-cli.yml`](go-cli.yml) | Go command-line tool |
| [`rust-cli.yml`](rust-cli.yml) | Rust binary |
| [`python-library.yml`](python-library.yml) | Python package on PyPI |
| [`monorepo-node.yml`](monorepo-node.yml) | pnpm workspace with many packages |
| [`microservice-go.yml`](microservice-go.yml) | Containerised HTTP service |
| [`editor-plugin.yml`](editor-plugin.yml) | Editor or plugin host project |
| [`custom-scaffold.yml`](custom-scaffold.yml) | `custom.files` and `custom.folders` |
| [`ignore-patterns.yml`](ignore-patterns.yml) | Every `ignore` form |
| [`ci-output.yml`](ci-output.yml) | Tuned for a specific CI system |
| [`baseline.txt`](baseline.txt) | Accepted violations, for gradual adoption |
| [`psx.examples.yml`](psx.examples.yml) | The original example, kept for compatibility |

## A YAML trap worth knowing

Quote any `ignore` pattern containing `*`, `?` or `[`. Unquoted, YAML reads `*` as an
alias and the file fails to parse:

```yaml
ignore:
  - "*.tsbuildinfo"    # correct
  - *.tsbuildinfo      # parse error: "could not find alias"
```

## Rules

`rules` is an allow-list: name a rule and only that one runs. Omit the key entirely to
enable every rule at its default severity.

Severity is `error`, `warning`, `info`, or `false` to disable. See
[CONFIGURATION.md](../docs/CONFIGURATION.md) for the full schema and
[RULES.md](../docs/RULES.md) for what each rule checks.

## A note on severities

`error` fails CI. Reserve it for what the team has actually agreed to enforce — a config
full of `error` gets ignored, and then `--baseline` quietly absorbs everything.

A workable split:

- `error` — the project does not work without it (README, a dependency manifest, tests)
- `warning` — should be fixed, but will not block (LICENSE, `.env.example`, Dependabot)
- `info` — worth having, purely advisory (roadmap, CODEOWNERS)

## Ordering

`psx rules` lists the ids. Copy from there rather than typing them. [CONFIGURATION.md](../docs/CONFIGURATION.md)
explains what each severity does to the exit code.