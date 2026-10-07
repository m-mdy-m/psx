# GitHub Actions templates

psx bundles 21 workflow templates, grouped so a project adopts only what it needs. Each one
is a standalone file: copy it, or let `psx fix` write it, and it works without depending on
any other workflow existing.

```bash
psx workflows                        # the catalogue
psx workflows --group docker         # one group
psx workflows --show                 # print the templates
psx workflows --json                 # for scripting
```

## Groups

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
| `NPM_TOKEN` | `release_nodejs` |
| `DEPLOY_HOST`, `DEPLOY_SSH_KEY` | `deploy_container` |

`docker_ghcr`, `codeql`, `security_secret_scan` and `release_aggregate` use the automatic
`GITHUB_TOKEN` and need nothing.

## Reviewing before you commit

A generated workflow runs with your repository's permissions. Read it before pushing,
particularly anything that publishes:

```bash
psx workflows --group docker --show > review.yml
```

The bundled workflows pin actions to major versions and grant the minimum `permissions`
each job needs. Pinning to commit SHAs is stricter still; the templates do not do it because
it makes every security update a manual edit.