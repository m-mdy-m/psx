# psx

**Check that a project has the files it should have — and create the ones it doesn't.**

Projects drift. The README describes a tool that no longer exists. There are no tests.
Nobody remembers why `Dockerfile` is there. A new contributor has no idea how to build it.

psx finds that drift and fixes it.

---

## Install

```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex
```

One static binary. Nothing else to install.

---

## 60-second tour

### 1. Say what kind of project this is

```bash
psx init
```

psx looks at your files, works out that this is a Go command-line tool, and writes a
`psx.yml`:

```yaml
version: 1
project:
  type: go
  kind: cli
rules:
  readme: error
  license: warning
  gitignore: info
  # ... and ten more
```

One rule is an error, one is a warning, the rest are suggestions. `rules` is an
allow-list, so this file is the list of what gets checked.

**Skipping this step is fine** — psx checks all 43 rules instead, and reports 41 findings on
a fresh project. Most are rules you will never care about.

### 2. See what's missing

```bash
psx check
```

```
Errors (1)

  error  readme
        No README file found in the project root
        fix: psx fix --rule readme

Warnings (1)

  warn   license
        No LICENSE file found
        fix: psx fix --rule license

Info (11)

  info   adr
        No ADR directory found
        fix: psx fix --rule adr
  ...

Result: 1 errors, 1 warnings, 11 info
Status: FAILED
```

Nothing was changed. `check` only reads — it never writes to your project and never asks
you anything, so it is safe in scripts and CI.

### 3. Create the missing files

```bash
psx fix --dry-run     # preview — changes nothing
psx fix               # do it
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

Then `psx check` again:

```
Result: 2 info
Status: PASSED
```

The two leftovers are advisory: a lockfile, which only your package manager can produce, and
a release workflow, which you only need once you publish something.

Two things worth knowing:

- **It will not overwrite your work.** A file that already has content is left alone.
- **It always converges.** Run `fix` twice and the second run finds nothing, because it
  only ever creates what is missing.

---

## Reading the output

| Part | Meaning |
| --- | --- |
| `error` | Serious. Fails the build. |
| `warn` | Should be fixed, but won't fail the build. |
| `info` | Nice to have. Purely advisory. |
| The rule name | Which rule failed. `psx explain <rule>` describes it. |
| `fix:` | The exact command that resolves it. |
| `Result:` / `Status:` | Totals, and whether the run passed. |

Every finding is a command you can copy and run.

---

## Choosing what matters

`psx init` gives you a reasonable set, but you decide. Edit `psx.yml`:

```yaml
rules:
  readme: error       # fail the build
  tests_folder: error # fail the build
  license: warning    # mention it, don't fail
  dockerfile: false   # don't mention it at all
```

**Anything you don't list is switched off.** That file *is* your rule list.

Ready-made configs for Go CLIs, npm packages, Rust binaries, Python packages, monorepos
and microservices are in [examples/](examples/README.md).

Full reference: [Configuration](docs/CONFIGURATION.md) · [All rules](docs/RULES.md)

---

## While you work

```bash
psx watch            # re-check as you edit
psx watch --fix      # and create missing files as you go
```

It prints only what changed, so a long session stays readable.

---

## In CI

```yaml
- run: psx check
```

That is the whole job definition — `check` already exits non-zero on any error-severity
failure. Add `--fail-on warning` if you want warnings to block too.

To annotate the pull request instead of dumping text into the log:

```bash
psx check -o github     # inline annotations on the diff
psx check -o sarif      # uploads to GitHub code scanning
```

Turning it on for a repo that already has problems? Run it once and psx writes down what it
found, so only *new* problems break the build:

```bash
psx check --baseline .psx-baseline.txt    # writes the file, then passes
```

---

## Commands

| Command | What it does |
| --- | --- |
| `check` | Report problems. Read-only. |
| `fix` | Create the missing files. |
| `watch` | Re-check as you edit, optionally fixing. |
| `init` | Write a `psx.yml` suited to the project. |
| `rules` | List everything psx can check. |
| `explain <rule>` | Describe one rule in detail. |
| `workflows` | List the bundled GitHub Actions templates. |
| `detect` | Show what psx thinks your project is. |

---

## Documentation

**Start here**
- [Getting started](docs/GETTING-STARTED.md) — a walkthrough from an empty folder

**When you need detail**
- [Configuration](docs/CONFIGURATION.md) — every setting, explained
- [Rules](docs/RULES.md) — all 43, and what each one creates
- [GitHub Actions templates](docs/WORKFLOWS.md) — 21 ready-made workflows
- [Examples](examples/README.md) — configs for different project types

**If you're contributing**
- [Architecture](docs/ARCHITECTURE.md) · [Contributing](docs/CONTRIBUTING.md)
- [Installation](docs/INSTALLATION.md) · [Verification log](docs/VERIFICATION.md)

---

## License

[MIT](LICENSE) © m-mdy-m