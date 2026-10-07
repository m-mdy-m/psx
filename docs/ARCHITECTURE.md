# Architecture

psx is a single static binary with all rules and templates embedded. There is no config
schema to fetch, no plugin registry, and no network access at runtime.

## The load-bearing decisions

### Rules are data, not code

A rule lives entirely in `internal/config/embedded/rules.yml`:

```yaml
tests_folder:
  id: TESTS_FOLDER_REQUIRED
  category: structure
  severity: error
  patterns:
    go: ["**/*_test.go"]
    nodejs: ["test/", "**/*.test.ts"]
    generic: ["tests/"]
  message: "No tests found"
  fix:
    path: tests/
    template: docs_index
    mode: 493
  fix_hint: "psx fix --rule tests_folder"
```

Adding a rule means editing YAML. The Go packages contain no `switch` on a rule id, which
is what previously made every new rule a three-file change.

Two constraints keep this honest:

- **`fix.path` is always a literal.** A glob describes something to detect, not something
  safe to create. A rule that detects `**/*_test.go` must not try to create a directory
  literally named `**`.
- **A template must exist.** `TestEveryFixTemplateExists` fails the build if a rule names a
  template that does not exist, so a rename cannot silently break every fix.

### One tree walk, then in-memory matching

`tree.Scan` walks the project once, pruning ignored directories, and builds an immutable
index. Every rule then matches against that index in memory.

This replaced per-pattern `os.Stat` and `filepath.Glob` calls, which had two problems beyond
speed:

- `filepath.Glob` does not support `**`. `**/*_test.go` returned **zero** matches, so
  `tests_folder` failed on every Go project regardless of its tests.
- Glob matching now handles `**` and brace expansion (`**/*.{test,spec}.{js,ts}`) via
  `doublestar`.

A literal prefix is extracted from each pattern to narrow the candidate set, so a rule
looking for `docs/**` does not walk every file in the repository.

### Tri-state results

A rule is `passed`, `failed`, or `skipped`. Skipped means the rule's patterns did not resolve
for this project type — a Go-only rule on a Python project.

Previously "no patterns matched" was reported as a **pass**, which silently hid real
violations: with `project.type: ts` the type was never normalised, no patterns matched, and
every Go-specific rule reported success on a project with no tests. Skips are counted
separately and shown with `--show-skipped`.

### Read-only versus interactive is explicit

`cmdctx.Mode` states what a command may do:

| Mode | Prompts | Writes |
| --- | --- | --- |
| `ModeReadOnly` | no | no |
| `ModeDryRun` | yes | no |
| `ModeInteractive` | yes | yes |

`check` uses `ModeReadOnly`. It used to read `fix`'s flags, and since `--interactive`
defaulted to true, every `psx check` opened a prompt form and wrote `.psx-project.yml`.

Prompting additionally requires a real terminal. In CI, in a pipe, or with
`PSX_NON_INTERACTIVE` set, prompts return their default rather than blocking on stdin.

### stdout carries only the report

Every diagnostic goes to stderr. `psx check -o json | jq` used to fail because
"Configuration loaded and validated" was printed to stdout first.

### Fixes are declarative and idempotent

A fix resolves its rule's `fix` block to a list of literal paths and writes only what is
missing. A file with content is never overwritten without `--force`.

`TestFixIsIdempotent` runs `fix` twice and requires the second run to report nothing, so
`psx fix` always converges.

Paths from `custom` entries are resolved against the project root and rejected if they
escape it, symlinks included.

### One template resolution path

`resources.Template(name, projectType, vars)` is the only way a template body is
obtained. It covers documents, container files, developer scripts and GitHub Actions
workflows alike. A rule names a template; nothing else resolves content.

`TestEveryTemplateRenders` walks every template for every language and fails on any
unresolved `{{placeholder}}`. That test is why generated scripts can no longer contain a
literal `{{test_command}}`.

Some tokens are intentionally left for a human — `{{number}}` and `{{title}}` in the ADR
template. `resources.IsHumanPlaceholder` marks them so they are not treated as bugs.

## Testing

| Package | Covers |
| --- | --- |
| `tree` | Glob semantics, ignore rules, negation, snapshot behaviour |
| `rules` | Engine verdicts, ignore handling, idempotence, generated content |
| `resources` | Template rendering, placeholder safety, workflow validity |
| `report` | Every output format, JSON purity, CI formats |
| `config` | Generated rule reference |

The suites worth knowing about, because they each caught a class of bug:

- `TestEveryTemplateRenders` — no template may reference a variable the renderer lacks.
- `TestGeneratedYAMLIsValid` — every generated `.yml`/`.yaml` must parse. Found
  `FUNDING.yml` written as prose, which GitHub silently ignores.
- `TestWorkflowTemplatesAreValidYAML` — a broken workflow fails at push time otherwise.
- `TestGeneratedContentHasNoRawPlaceholders` — end-to-end, through the real fixer.
- `TestFixIsIdempotent` — `fix` must converge.
- `TestEveryFixTemplateExists` — a template rename cannot silently break fixes.

## Extending

**A rule** — add it to `rules.yml`. If it needs new content, add a template to the relevant
embedded file; no Go change.

**A workflow** — add it to `github-actions.yml` and register the name in
`internal/resources/actions.go`, then add it to a group in `WorkflowGroups`.

**A language** — add a block to `languages.yml` with its manifests, patterns and commands.
Generated scripts, `.gitignore` bodies and detection all read from it.

**An output format** — implement `report.Renderer` and register it in `report.New`.

## Deliberate non-goals

- **No filesystem watcher.** `watch` polls. It behaves identically on Linux, macOS, Windows,
  WSL and network shares, with no `fsnotify` dependency. Polling cost is proportional to
  file count, not to the number of rules.
- **No LSP.** Structured output exists (`-o ndjson`) so an editor plugin can consume it
  without psx hosting a language server.
- **No workspace-wide rule scoping yet.** `detect` reports `monorepo` and the workspace
  globs, but rules still evaluate the project root.