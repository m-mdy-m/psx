# Security Policy

## Supported versions

| Version | Supported |
| --- | --- |
| 3.x | Yes |
| 2.x | Security fixes only |
| < 2.0 | No |

Use the latest release.

## Reporting a vulnerability

Please do not open a public issue for a security problem.

Email **bitsgenix@gmail.com** with "PSX Security" in the subject.

Include:

- what the issue is and what an attacker gains
- steps to reproduce
- affected version and platform
- a suggested fix, if you have one

You will get an acknowledgement within 48 hours. If it is confirmed, a fix and a patch
release follow, and you are credited unless you prefer otherwise.

If it turns out not to be a security issue, you will be told why and the conversation may
move to a normal issue.

## What psx does with your project

psx reads your project and, when you run `fix`, writes into it. It makes no network
requests.

- `check`, `detect`, `rules`, `explain` and `workflows` never write to your project.
- `fix` writes only inside the project directory. Paths are resolved against the project
  root and rejected if they escape it, including through symlinks. A `psx.yml` from a
  cloned repository therefore cannot make `psx fix --yes` write elsewhere on disk.
- An existing file with content is never overwritten unless you pass `--force`.
- Nothing is deleted.

The trade-off: `fix` writes files into your working tree without a backup unless you set
`fix.backup: true` in `psx.yml`. Commit or stash first, or use `--dry-run`.

## Verifying a download

Release builds publish `checksums.txt`. Verify before running an installer:

```bash
sha256sum -c checksums.txt
```

`scripts/install.sh` does this for you.

## Hardening your own repository

psx can set up the checks that catch the common cases:

```bash
psx fix --rule workflow_secret_scan          # gitleaks on push and PR
psx fix --rule workflow_codeql               # CodeQL analysis
psx fix --rule workflow_dependency_updates   # Dependabot
psx fix --rule env_example                   # .env.example, so .env stays untracked
```

These are conveniences, not a security boundary. `psx` does not replace `gitleaks`,
`trivy`, or a dependency audit.

## For contributors

- Never commit a secret. `.gitattributes` and `.gitignore` templates mark common binary and
  lockfile paths, but a secret is your responsibility.
- Dependencies are pinned in `go.sum`. Run `make tidy` and commit the result.
- `make verify` checks the module graph is tidy and the tree is formatted.
- Report a vulnerable dependency through the Dependabot alert, or by email as above.