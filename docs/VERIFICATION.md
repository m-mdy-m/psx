# Verification

Every bug fixed in 3.0.0 was **reproduced first**, then fixed, then locked down with a
regression test. The list of what changed is in the
[changelog](../CHANGELOG.md); this page explains how each item was verified and how to
re-run the checks.

---

## Why this page exists

A changelog entry is a claim. This page is the evidence behind those claims, so you can
tell the difference between "we think this is fixed" and "we watched it fail, then watched
it pass".

---

## The rule we worked to

For each bug:

1. **Write a failing test.** If the bug cannot be expressed as a test, that is a signal the
   behaviour is not yet pinned down. We wrote the test anyway, or documented why not.
2. **Watch it fail.** The output must show the wrong behaviour, not a setup error.
3. **Fix it.** Usually the smallest change that makes the test pass.
4. **Watch the whole suite pass.** A fix that breaks something else is not a fix.

This is why several fixes came with new test files: `internal/command` had **no tests at
all** when the `--baseline` bug was found. The bug survived precisely because nothing
exercised it.

---

## Reproducing the fixes yourself

```bash
make verify
```

That runs, in order:

| Step | Command | What it proves |
| --- | --- | --- |
| Tidy | `go mod tidy` | Dependencies are correctly declared. |
| Vet | `go vet ./...` | No suspicious constructs. |
| Test | `go test ./...` | Every regression test passes. |
| Modules | `git diff go.mod go.sum` | `tidy` produced no uncommitted change. |
| Format | `gofmt -l .` | The tree is formatted. |

It also fails if `go.mod` or `go.sum` would change, which catches an out-of-date
dependency list before it reaches CI.

`docs/RULES.md` is checked separately, by a test that regenerates it and compares.

Individually:

```bash
go test ./internal/command/ -run TestMissingBaselineFileIsRecordedNotFatal -v
```

Each test name below is a regression test named after the behaviour it protects.

---

## The regression tests that guard the 3.0.0 fixes

### Generated files stay valid

| Test | Guards against |
| --- | --- |
| `TestEveryTemplateRenders` | A template with a typo silently producing an empty file |
| `TestGeneratedYAMLIsValid` | psx writing YAML that no tool can parse |
| `TestWorkflowTemplatesAreValidYAML` | The 21 bundled workflows being malformed |
| `TestExamplesAreValidConfigs` | A shipped example config that psx itself rejects |

### Fixes are safe and repeatable

| Test | Guards against |
| --- | --- |
| `TestEveryFixTemplateExists` | A rule advertising a fix that has no template |
| `TestFixIsIdempotent` | `fix` rewriting or duplicating files on a second run |
| `TestCustomFolderLeafIsAFileNotADirectory` | `custom.folders` creating a directory named like a config file |
| `TestCustomFolderCreationOrderIsStable` | Creation order varying with Go's map iteration |

### Configuration is not silently ignored

| Test | Guards against |
| --- | --- |
| `TestExamplesAreValidConfigs` | Example configs drifting from the schema |
| `TestRuleReferenceIsUpToDate` | `docs/RULES.md` drifting from `rules.yml` |
| `TestConfigFlagOverridesDiscovery` | An explicit `--config` losing to auto-discovery |
| `TestInvalidConfigFailsLoudly` | A broken config degrading to defaults and checking the wrong rules |

### Detection and reporting

| Test | Guards against |
| --- | --- |
| `TestDetectionRunsWhenTypeIsUnset` | Rules silently skipping when `project.type` is missing |
| `TestSummaryCountsOnlyVisibleFindings` | `--level` hiding findings the summary still counted |
| `TestCompactSummaryAccountsForEveryFailure` | Compact output reporting 12 failures but accounting for 1 |
| `TestEveryRelativeLinkResolves` | Documentation pointing at files that do not exist |

### Baselines

| Test | Guards against |
| --- | --- |
| `TestMissingBaselineFileIsRecordedNotFatal` | Baselines being impossible to create |
| `TestBaselineForgivesOnlyRecordedRules` | A baseline absorbing problems it does not list |
| `TestBaselineStillLetsNewWatchIssuesThrough` | A baseline hiding genuinely new failures |

---

## Reading a failure

A failing test is not an obstacle to work around. It is the specification telling you
something you did not know. Two rules:

- **Read the assertion before the code.** The test states the intended behaviour in
  plain terms; that sentence is the requirement.
- **Do not weaken a test to make it pass.** If a test is wrong, fix it deliberately and
  say why. A deleted assertion is a bug that comes back.

---

## What is still not covered

Being honest about the gaps:

- **Prompts are tested by their parsing, not by their interaction.** `ui` splits the
  answer parsing from the terminal handling so it can be tested at all — under `go test`
  there is no terminal, so `IsInteractive` is always false. The rendering of a prompt,
  and what a real user sees, is still unverified.
- **Templates are checked for validity, not for quality.** psx can prove a generated
  `Makefile` is syntactically correct, but not that its targets are the ones you want.
- **Fixes are not verified by running what they generate.** psx creates a
  `.github/workflows/ci.yml` and parses it back as YAML, but does not execute the
  workflow.
- **The watcher is tested against a real filesystem, not a real editor.** The timing
  tests use millisecond intervals, so a change that only appears under heavy filesystem
  contention would not be caught.
- **`--fail-on` and the exit codes are only tested at the unit level.** No test asserts
  the exit status of the built binary under every failure mode.

These are the first things to add tests for. See
[Contributing](CONTRIBUTING.md).