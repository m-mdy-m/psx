# Configuration

**Most people never need this page.** Run `psx init` and psx writes a file that already
suits your project. Come back here when you want to change what it checks.

---

## Do I need a config file?

No. Without one, psx runs all 43 rules at their default severity:

```bash
psx check          # works with no configuration at all
```

That is useful for a one-off look, but noisy on day one: it reports rules you have no
intention of satisfying, like ADR directories for a script.

A config file is how you say **which rules matter to you**. That is its whole job.

---

## The smallest useful file

This is a complete, working configuration. Every line is explained underneath.

```yaml
version: 1

project:
  type: go

rules:
  readme: error
  tests_folder: error
  license: warning
```

Three things to notice:

**`rules` is an allow-list.** Only these three run. Anything you leave out is not checked.
That single fact explains most of psx's behaviour.

**Severity is the whole point.** `error` fails your build, `warning` is advice, `info` is
a suggestion, and `false` means psx will not mention the rule at all.

**`project.type` is optional.** Leave it out and psx works it out from your files. Set it
only when the guess is wrong.

---

## Getting one written for you

```bash
psx init                        # detect the type, enable a sensible subset
psx init --minimal              # only the essentials
psx init --type nodejs --kind library
psx init --force                # overwrite an existing file
```

`psx init` reads your layout, decides what kind of project this is, and writes a `psx.yml`
listing the rules that make sense for it. Read it afterwards — editing it is the normal
workflow, not a fallback.

Check what psx thinks about your project:

```bash
psx detect
```

```
type       go
kind       cli
confidence 97%
signals    go.mod, cmd/
```

If that guess is wrong, set `project.type` yourself. An explicit choice always beats
detection.

---

## Where the file goes

psx accepts `psx.yml`, `.psx.yml`, `psx.yaml` and `.psx.yaml`, and looks in this order:

1. the path you passed to `--config`
2. those four names in the project directory
3. the same four names at the git root, walking up at most 10 levels
4. `~/.config/psx/`
5. built-in defaults

Step 3 means a `psx.yml` at the top of a monorepo covers every package inside it.

---

## `project`

```yaml
project:
  type: go            # nodejs | go | rust | python | generic
  kind: cli           # app | library | cli | monorepo | microservice | plugin
```

`type` picks the language tier. Short names resolve: `ts`, `js`, `node` and `typescript`
all mean `nodejs`; `golang` means `go`.

Leave it empty, or set `auto`, to detect from the layout — `go.mod`, `package.json`,
`Cargo.toml`, `pyproject.toml`.

This matters because rules define their patterns **per language**. A rule that only has Go
patterns is reported as **skipped** for a Python project rather than quietly passing. Skips
show up in the summary and in `--show-skipped`, so a green result never hides a rule that
never actually ran.

`kind` is looser. It mostly affects which templates get picked.

---

## `rules`

```yaml
rules:
  readme: error       # fail the build
  license: warning    # mention it, don't fail
  changelog: info     # a suggestion
  dockerfile: false   # don't mention it at all
```

A rule id you misspell is reported as a warning and skipped, so a typo produces a message
rather than a silently smaller check.

See [Rules](RULES.md) for the full list of ids.

---

## `ignore`

Directories and files psx should not look at.

```yaml
ignore:
  - dist/               # any directory named dist
  - /build              # only at the project root
  - "**/*.snap"         # nested snapshots
  - "logs/*"            # the contents of logs
  - "!logs/keep.txt"    # but keep this one
```

The syntax is gitignore's, so it should feel familiar: a trailing `/` means directory, a
leading `/` anchors to the root, a leading `!` re-includes.

Ignored directories are **pruned during the scan**, not filtered afterwards. Without this, a
copy of a library in `vendor/` could satisfy `tests_folder` and convince you your project
has tests.

---

## `fix`

```yaml
fix:
  interactive: true   # confirm each change
  backup: false       # copy the original into .psx/backup before overwriting
```

`interactive` is implied false when there is no terminal, so `fix` never hangs in CI.

---

## `custom`

Anything psx has no rule for. Use it for house conventions rather than inventing a rule id.

```yaml
custom:
  files:
    - path: .env.example
      content: |
        APP_ENV=development
        LOG_LEVEL=info
  folders:
    - path: config
      structure:
        base:
          production.yaml: {}
          staging.yaml: ""
```

A key **with children** becomes a directory. A key **with none** becomes a file — so
`production.yaml: {}` is an empty file, not a directory named like a config file. Give a
leaf a string to write content into it:

```yaml
custom:
  folders:
    - path: config
      structure:
        base:
          production.yaml: "replicas: 2\n"
```

Both are applied after the built-in rules, so `custom` files are created last and win.

Paths are resolved against the project root and **rejected if they escape it**. A `custom`
entry of `../../.bashrc` fails rather than writing outside the project — a configuration
file in a cloned repository is untrusted input.

---

## Command-line overrides

Useful for a one-off, without editing the file:

```bash
psx check --only readme,license        # run just these
psx check --category documentation     # run a whole group
psx check --level warning              # hide info-level findings
psx check --fail-on none               # report but never exit non-zero
psx check --show-skipped               # include rules that did not apply
psx check --all                        # include rules that passed
```

`--config` points at a different file entirely, which is how you try a config before
committing to it:

```bash
psx check --config ./psx.ci.yml
```

---

## Baselines

Turning psx on for a repository that already has problems should not fail the build on day
one. A baseline forgives problems you already have while still failing on new ones.

```bash
psx check --baseline .psx-baseline.txt
```

**If the file does not exist, psx writes it** from whatever the current run reports, and
that run is treated as already forgiven:

```bash
$ psx check --baseline .psx-baseline.txt
PASSED: 42 passed, 1 skipped
```

Now commit the file and keep using the same command in CI. Anything **not** listed in it
still fails, so only pre-existing debt is absorbed:

```bash
$ psx check --baseline .psx-baseline.txt
FAILED: 1 failed (1 errors)
```

That is what happens the moment a listed rule starts failing for a new reason — remove the
line only once you have genuinely fixed it.

The file can be plain lines, a JSON array, or JSON objects:

```
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

`psx watch --baseline` honours the same file, so a long editing session only reports
problems you have not already recorded.

---

## When psx pushes back

Setting `readme: false` or `license: false` prints a warning:

```
⚠ Critical rule 'readme' is disabled - this is not recommended
```

It is only a warning — your configuration wins, and the rules stay off. But those two
cover the two things a stranger looks for first when they open a repository, so it is
worth knowing you switched them off deliberately.

This applies to `false` only. Lowering `readme` to `info` is fine.

---

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Passed, or failed below the `--fail-on` threshold |
| 1 | Validation failed |
| 2 | Configuration or usage problem |
| 3 | Invalid arguments |

The distinction matters in CI: **1** means "your project has problems", **2** means "psx
could not run" and the report is not trustworthy. Alert on 2, ignore 1.

---

## Project metadata

Some templates need your name and repository — a README header, a `CODEOWNERS` file. psx
reads those from git, asks for anything it cannot infer the first time you run `fix`, and
caches the answers in `.psx-project.yml`.

`--answer` overrides one template question. Exactly one key is read today:

```bash
psx fix --rule docker_compose --answer with_database=yes
```

Without it the compose file has one service; with it, `app` plus a `postgres` service and a
named volume. Any other key is accepted and ignored, so a typo here fails silently.

Project metadata is **not** set through `--answer` — it comes from the prompt, then from git
and the environment.

Decide whether to commit that file or add it to `.gitignore` — your call, both work.

`check` never prompts and never writes it. When psx runs without a terminal (CI, a pipe, or
`PSX_NON_INTERACTIVE` set) it skips the questions entirely and uses the values it inferred,
so `fix` never blocks.