# GitHub Actions templates

CI configuration is the most tedious part of setting up a repository, and it is also the
part you are most likely to get wrong in a way nobody notices for months. psx bundles 21
ready-made workflows so you can start from a working one instead of a blank file.

**You do not need any of them.** psx checks that a workflow *exists*; it does not judge
whether yours is correct. Adopt these only if you want them.

---

## Start here

Most projects need one workflow, not twenty-one. This writes the CI for your language:

```bash
psx fix --rule workflow_ci
```

It picks `ci_go`, `ci_nodejs`, `ci_rust` or `ci_python` based on your project, so a Go
project gets `gofmt`, `go vet` and `go test -race`; a Node project gets pnpm, lint, typecheck
and test.

Then look at what it produced before you push it:

```bash
psx workflows --show
```

A generated workflow runs with your repository's permissions. Read it first — especially
anything that publishes.

---

## When you want them

```bash
psx fix --rule workflow_codeql       # CodeQL scanning
psx fix --rule workflow_secret_scan  # catch committed secrets
psx fix --rule workflow_docker       # build and publish images
psx fix --rule workflow_release      # publish from a version tag
```

Those are the four rules that produce a workflow. Each writes its own file, so you adopt
them one at a time rather than all at once.

Some templates are not reachable from a rule — the `deploy_*` group, for instance. Print one
and copy it into place by hand:

```bash
psx workflows --group deploy --show
```

---

## Looking around

```bash
psx workflows                        # every template, grouped
psx workflows --group docker         # just one group
psx workflows --show                 # print the full YAML
psx workflows --json                 # for scripting
```

`--show` prints a table **and** the YAML bodies, so redirecting it captures the table too.
Copy from the `───── name ─────` section rather than saving the whole output.

```bash
psx workflows --group deploy --show    # read it
psx workflows --json > templates.json  # or work with it programmatically
```

---

## The catalogue

### `ci` — validate every push and pull request

| Template | For |
| --- | --- |
| `ci_go` | `gofmt`, `go vet`, `go test -race`, coverage artifact |
| `ci_nodejs` | pnpm, lint, typecheck, test, build, frozen lockfile |
| `ci_rust` | `cargo fmt`, `clippy -D warnings`, `cargo test` |
| `ci_python` | ruff, pytest with coverage |
| `ci_monorepo` | `pnpm -r` across every workspace member |

```bash
psx fix --rule workflow_ci
```

### `release` — publish from a version tag

| Template | For |
| --- | --- |
| `release_go` | Go binaries, checksums, GitHub Release |
| `release_nodejs` | Test, `pnpm publish --provenance`, GitHub Release |
| `release_rust` | `cargo test`, goreleaser cross-compile and sign |
| `release_linux` | Linux x86_64 tarball, `workflow_call` worker |
| `release_windows` | Windows x86_64 zip, `workflow_call` worker |
| `release_macos` | macOS arm64 tarball, `workflow_call` worker |
| `release_aggregate` | Orchestrator: calls the three workers, publishes once |

```bash
psx fix --rule workflow_release      # writes release.yml
```

`release_aggregate` is the interesting one: it delegates to the per-platform workflows via
`workflow_call` and publishes a single release with notes from `CHANGELOG.md` and an
attribution list. Each worker also runs standalone, so a platform can be debugged on its own.

### `docker` — build and publish images

| Template | For |
| --- | --- |
| `docker_build` | Build on every push and PR, never push, GHA layer cache |
| `docker_hub` | Multi-arch Docker Hub publish from a tag |
| `docker_ghcr` | GitHub Container Registry, push on non-PR |

```bash
psx fix --rule workflow_docker
```

`docker_hub` follows the digest pattern, which is the part worth copying:

1. Each platform builds and pushes **by digest only** — no tags yet.
2. Digests are uploaded as artifacts.
3. A final job merges them into one manifest and applies the real tags.

Tags are applied exactly once, after every platform has succeeded. A half-finished release
can never leave `latest` pointing at a missing architecture.

It also uses `ubuntu-24.04-arm` for arm64, so there is no QEMU emulation and the build is
not slow.

### `deploy` — ship the project

| Template | For |
| --- | --- |
| `deploy_github_pages` | Vite/pnpm site to Pages, with an SPA fallback |
| `deploy_docs` | MkDocs to the `gh-pages` branch, built with `--strict` |
| `deploy_container` | Pull the image and `docker compose up` over SSH |

### `security` — scan

| Template | For |
| --- | --- |
| `codeql` | Weekly plus every PR, Go and JS/TS |
| `security_dependency_review` | Block a PR that adds a vulnerable dependency |
| `security_secret_scan` | gitleaks on push and PR |

```bash
psx fix --rule workflow_codeql
psx fix --rule workflow_secret_scan
```

## Variables

Templates are rendered with your project metadata:

| Variable | Source |
| --- | --- |
| `{{project_name}}` | Project name |
| `{{repo_name}}` | Repository name |
| `{{github_username}}` | GitHub user or organisation |
| `{{docker_image}}` | `<user>/<repo>`, lowercased |
| `{{domain}}` | GitHub Pages domain |

GitHub's own `${{ ... }}` expressions are left untouched.

## Secrets

The templates expect these repository secrets, and only the ones matching the workflows you
adopt:

| Secret | Used by |
| --- | --- |
| `DOCKER_USERNAME`, `DOCKER_TOKEN` | `docker_hub` |
| `REGISTRY_USERNAME`, `REGISTRY_PASSWORD`, `REGISTRY_URL` | `deploy_container`, to pull a private image |
| `NPM_TOKEN` | `release_nodejs` |
| `DEPLOY_HOST`, `DEPLOY_SSH_KEY` | `deploy_container` |

`docker_ghcr`, `codeql`, `security_secret_scan` and `release_aggregate` use the automatic
`GITHUB_TOKEN` and need nothing configured.

A workflow referencing a secret you have not set fails at the point of use, not at parse
time. Adopt the CI and security workflows first — they need nothing — and add publishing
workflows once the credentials exist.

## Reviewing before you commit

A generated workflow runs with your repository's permissions. Read it before pushing,
particularly anything that publishes:

```bash
psx workflows --group docker --show
```

The bundled workflows pin actions to major versions and grant the minimum `permissions`
each job needs. Pinning to commit SHAs is stricter still; the templates do not do it because
it makes every security update a manual edit.