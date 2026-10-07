# Contributing

## Setup

Requires Go 1.25+ and make.

```bash
git clone https://github.com/m-mdy-m/psx
cd psx
make build
make test
```

## Style

```bash
make fmt     # gofmt -s -w .
make lint    # golangci-lint run
make check   # fmt, vet, test, lint
```

Comments explain why, not what. Keep them to a line or two. If a comment restates the code,
delete it.

## Adding a rule

Rules are declarative. You should not need to write Go.

1. Add an entry to `internal/config/embedded/rules.yml`.
2. If it needs generated content, add a template to the matching file in
   `internal/resources/embedded/`.
3. Add the rule id to `psx.default.yml` if it should be on by default.
4. Regenerate the docs: `make docs`.

```yaml
my_rule:
  id: MY_RULE_REQUIRED
  category: quality
  description: One line, shown in `psx rules`
  severity: warning
  patterns:
    go: [".golangci.yml"]
    generic: [".golangci.yml"]
  message: "No .golangci.yml found"      # shown when it fails
  fix:
    path: .golangci.yml                  # literal path, never a glob
    template: editorconfig               # a registered template
    safe: true                           # allows `watch --fix`
  fix_hint: "psx fix --rule my_rule"
  doc_url: ""
```

### Rules for writing one

- **`fix.path` must be literal.** A glob says what to detect, not what to create. Create a
  real starter file (a `.gitkeep`, a `tests/README.md`) rather than a directory a glob
  matched.
- **Give it a real `message`.** It is what the user reads when the rule fails.
- **`safe: true` only if it is non-destructive.** `watch --fix` applies safe rules
  unattended.
- **Omit `fix` if generation makes no sense.** The rule still reports; it is simply listed
  as `manual`.
- **Be honest about severity.** `error` fails CI. Reserve it for things that genuinely break
  a build.

### The tests that will check you

```bash
go test ./internal/config/ -run TestEveryFixTemplateExists
go test ./internal/rules/  -run 'TestEveryFixableRuleSucceeds|TestGeneratedYAMLIsValid'
go test ./internal/resources/ -run TestEveryTemplateRenders
```

Together these enforce that the template exists, renders, produces parseable YAML, contains
no unresolved `{{placeholder}}`, and can actually be written.

## Adding a template

Put it in the embedded file matching its kind:

| File | Contents |
| --- | --- |
| `templates.yml` | README, changelog, contributing, API docs |
| `docs-templates.yml` | Security, CoC, ADR, roadmap, runbook, OpenAPI |
| `quality-tools.yml` | editorconfig, pre-commit, gitattributes, Makefile |
| `devops.yml` | Docker, compose, CI, Dependabot, nginx, k8s, Helm |
| `github-actions.yml` | Every workflow template |
| `project-scripts.yml` | Developer scripts |
| `licenses.yml` | License texts |

Then register the name in `internal/resources/registry.go`.

Only reference variables from `ProjectInfo.ToVars()`. A new variable means adding it there;
a template using an undeclared one ships broken output.

Some tokens are deliberately left for a person to fill in — `{{number}}`, `{{title}}` in the
ADR template. Register them in `humanPlaceholders` so the placeholder checks skip them.

## Adding a GitHub Actions workflow

1. Add it to `github-actions.yml`.
2. Register the name in `internal/resources/actions.go`.
3. Add it to a group in `WorkflowGroups`.
4. Optionally expose it as a rule in `rules.yml` so `psx fix` can create it.

`go test ./internal/resources/ -run TestWorkflowTemplatesAreValidYAML` parses every workflow,
and `TestWorkflowTemplatesUseOnlyKnownVariables` catches leftover `{{...}}`.

Keep each workflow independently runnable. A project should be able to adopt one file
without another existing. `release_aggregate` is the deliberate exception: it is an
orchestrator whose only job is to call the others.

## Commit messages

```
fix: match ** globs in nested directories
feat: add lockfile rule
refactor: replace per-rule switches with declarative fix specs
docs: regenerate rule reference
```

## Pull requests

- One change per PR.
- Tests for behaviour changes. A bug fix should come with a test that failed before.
- Run `make check`.
- Update `CHANGELOG.md` under Unreleased.

## Regenerating docs

```bash
make docs    # rewrites docs/RULES.md from rules.yml
```

`docs/RULES.md` is generated. Do not edit it — the header says so, and the test will
overwrite your changes.