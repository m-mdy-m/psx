# Getting started

This walks through using psx on a project from scratch. Every command and every line of
output here is real — copy it as you go.

If you just want the short version, the [README](../readme.md) has a 60-second tour.

---

## What psx actually is

psx is a **linter for the shape of your project**, not your code.

It does not read your source files to find bugs. It looks at *which files exist* and asks:
does this project have a README? A license? Tests? A CI workflow? A `.gitignore`?

```
                    your project
                         │
        ┌────────────────┴────────────────┐
        │                                 │
   psx check                         psx fix
   "here's what's missing"          "here, I made it"
   (reads only)                      (creates files)
```

That distinction matters:

- It **never edits your code.** It creates files. That's all.
- It **never talks to the network.**
- It **never overwrites a file that already has content.**

---

## Before you start

You need a project to point it at. If you don't have one handy, make a throwaway:

```bash
mkdir ~/psx-demo && cd ~/psx-demo
go mod init example.com/demo      # or: npm init -y
```

Everything below works on any project — Go, Node, Rust, Python, or nothing in
particular.

---

## Step 1: Run it with no setup at all

```bash
psx check
```

On a fresh project with nothing in it, you get a lot of output. That's expected — with
no configuration psx checks everything it knows.

Two flags make the first look easier:

```bash
psx check --level warning    # hide the "info" suggestions
psx check --only readme,license   # just these two rules
```

Nothing is modified. `check` never writes to your project.

---

## Step 2: Write a configuration

That was a lot of noise, because psx was guessing. Tell it what matters.

```bash
psx init
```

psx inspects your files and writes a `psx.yml`. For a Go project with a `cmd/` directory
it produces:

```yaml
version: 1

project:
  type: go        # worked out from go.mod
  kind: cli        # worked out from the cmd/ directory

rules:
  readme: error
  license: warning
  gitignore: info
  # ...

ignore:
  - vendor/
  - dist/
  - build/
```

Three things to understand here:

**`rules` is an allow-list.** Anything not listed does not run. That single fact explains
most of the file's behaviour.

**Each rule has a severity**, and severity decides whether the build fails:

| Severity | Effect |
| --- | --- |
| `error` | `psx check` exits non-zero |
| `warning` | Reported, exit stays zero |
| `info` | Reported only |
| `false` | Not checked at all |

**`ignore` stops psx descending into those directories.** Without it, a copy of a library
in `vendor/` could make your tests look like they exist.

### If the generated file isn't right

Edit it. That is the normal workflow, not an escape hatch. Two common changes:

```yaml
# I don't care about containers at all
rules:
  ...
  dockerfile: false
  docker_compose: false
  kubernetes: false

# These two are non-negotiable for us
rules:
  readme: error
  tests_folder: error
```

Verify what psx thinks your project is:

```bash
psx detect
```

```
type       go
kind       cli
confidence 97%
signals    go.mod, cmd/
```

If that is wrong, set `type` yourself in `psx.yml`. An explicit choice always wins over
detection.

---

## Step 3: Look at one finding in detail

Every finding names a rule. To see everything about it:

```bash
psx explain tests_folder
```

```
tests_folder
  category    structure
  severity    error
  description Projects need automated tests

  message: No tests found

  checks for:
    [generic]
      tests/
      test/
    [go]
      **/*_test.go
    [nodejs]
      test/
      **/*.test.ts
    ...
  fix: none (manual)

  hint: psx fix --rule tests_folder
```

The `checks for` section is the important part: it lists exactly which paths satisfy the
rule, per language. For a Go project, **any** file matching `**/*_test.go` anywhere below
the root is enough — it does not have to live in a `tests/` directory. The first entry
under `[generic]` applies when psx has not recognised the language.

`fix: none (manual)` means psx will report this rule but will not create the file: there is
no sensible way to generate real tests.

To see every rule:

```bash
psx rules
psx rules --fixable      # only the ones that can create their own files
```

---

## Step 4: Create the missing files

```bash
psx fix --dry-run
```

```
would create  .gitignore
would create  LICENSE
would create  Makefile
would create  README.md
would create  scripts/build.sh
would create  scripts/clean.sh
would create  scripts/setup.sh
would create  scripts/test.sh
would create  SECURITY.md
...
Run without --dry-run to apply
```

Nothing was written. When the list looks right:

```bash
psx fix
```

Without `--yes` this asks before each change. Commit your work first — psx has no undo.

### What it will not do

- **It will not overwrite.** A file with content is skipped. Use `--force` if you really
  mean it.
- **It will not invent tests or source code.** `tests_folder` and `src_folder` are
  reported but not fixed, because there is no sensible way to generate them. Write those
  yourself.
- **It will not delete anything.**

### The one-time questions

Templates that mention your name or repository need to know them. The first `fix` in a
project asks, offering a guess for each field and saving the answers to
`.psx-project.yml`:

```
Project Information:

Project name [notes]:
Description [A notes project]:
Author [Sam Rivera]:
Email [sam@example.com]:
GitHub username [samrivera]:
Repository name [notes]:
License (MIT/Apache-2.0/GPL-3.0/BSD-3-Clause) [MIT]:
```

Press Enter to accept the guess in brackets. Commit that file or add it to `.gitignore` —
either is fine, but if you ignore it you will be asked again next time.

Running in CI or through a pipe, psx detects that nobody is there to answer, skips the
questions entirely and uses the values it inferred. It never blocks.

To skip the questions on a terminal too, write the file yourself:

```yaml
# .psx-project.yml
name: notes
description: A tiny CLI for taking notes
author: Sam Rivera
email: sam@example.com
github_user: samrivera
repo_name: notes
license: MIT
```

---

## Step 5: Confirm

```bash
psx check
```

```
Info (2)

  info   lockfile
        No lockfile found; installs will not be reproducible
        fix: Commit the lockfile your package manager generates
  info   release_workflow
        No release automation found
        fix: Run 'psx fix --rule workflow_release'

Result: 2 info
Status: PASSED
```

Both are advisory and neither can be fixed for you: a lockfile comes out of your package
manager, and a release workflow only matters once you publish something. `info` never fails
a build, so leaving them is fine.

Run `fix` a second time and it finds nothing. `fix` only creates what is missing, so it
always reaches a fixed point.

---

## Step 6: Keep it honest while you work

```bash
psx watch
```

Re-checks as you save, and prints only what changed:

```
› Watching for changes. Press Ctrl+C to stop.
⚠ 1 errors · 1 warnings · 0 skipped · watching
✓ readme: fixed
⚠ 0 errors · 1 warnings · 0 skipped · watching
```

That last `⚠` is the same warning as the second line — it is still there because no
LICENSE file exists yet. Create one, and the line turns into `✓`.

Nothing is printed when a save does not change any rule's status, so a long editing
session stays readable. Ctrl+C stops it.

Add `--fix` and it will also create files as they become missing. Only rules marked
"safe" are applied, and only for files that do not exist yet.

---

## Step 7: Put it in CI

The point of all this is that the checks run without you remembering them.

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  structure:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: go build -o psx ./cmd/psx   # or download a release binary
      - run: ./psx check
```

`psx check` already exits non-zero when any `error`-severity rule fails, so nothing extra is
needed. That is the whole configuration.

If you want warnings to break the build too, say so explicitly:

```bash
psx check --fail-on warning
```

And to go the other way — report everything but never block:

```bash
psx check --fail-on none
```

### Making the output useful

```bash
psx check -o github    # annotate the changed lines in the pull request
psx check -o sarif     # feed GitHub code scanning
psx check -o json      # for your own scripts
```

`-o github` is usually the best choice, because the report lands on the diff instead of
buried in a log.

### Introducing it to an existing repository

Turning it on usually surfaces problems that have been there for years. Making all of
them fail immediately is how adoption gets reverted. Record the current state instead:

```bash
psx check --baseline .psx-baseline.txt . > /dev/null
```

Then commit that file and use the same command in CI. Violations listed in it are reported
as known issues but do not fail the build. Anything **not** listed still fails, so only
debt you already had is absorbed.

Remove lines from the file as you fix them.

---

## Common questions

**Can I just check some rules?**
Yes. `psx check --only readme,tests_folder`, or `--category documentation` to run a whole
group.

**Will it touch my code?**
No. It creates files and nothing else.

**Is the config optional?**
Yes, but without one you get all 43 rules including the ones you do not care about.

**Two rules both fail — which do I fix first?**
Any order. They are independent.

**How do I see what a rule would create without running fix?**
`psx explain <rule>` prints the `fix:` section. `psx fix --dry-run` shows it for the
whole project.

**Can I add my own files, not just psx's?**
Yes — `custom.files` and `custom.folders` in `psx.yml`. See
[Configuration](CONFIGURATION.md#custom).

**Does it work on Windows?**
Yes. The installer and every command are cross-platform.

---

## Where to next

| You want to | Go to |
| --- | --- |
| Tweak severities and ignore lists | [Configuration](CONFIGURATION.md) |
| See every rule and what it creates | [Rules](RULES.md) |
| Copy a config for your project type | [Examples](../examples/README.md) |
| Set up GitHub Actions, Docker, releases | [Workflows](WORKFLOWS.md) |
| Install it another way | [Installation](INSTALLATION.md) |
| Work on psx itself | [Contributing](CONTRIBUTING.md) |