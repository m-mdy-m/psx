# Security Policy

## supported versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| 0.x     | :x:                |

basically just use the latest version.

## reporting a vulnerability

found a security issue? please don't open a public issue.

email me instead: bitsgenix@gmail.com

put "PSX SECURITY" in the subject so i don't miss it.

### what to include

- description of the issue
- steps to reproduce
- affected versions
- possible fix if you have one

### what happens next

i'll try to respond within 48 hours (might be longer on weekends).

if it's legit:
1. i'll confirm it
2. fix it
3. release a patch
4. credit you if you want

if it's not really a security issue, i'll let you know and maybe open a regular issue for it.

## security best practices

for users:
- always use the latest version
- don't run psx with sudo unless you have to
- check checksums when downloading binaries

for contributors:
- don't commit secrets or keys
- use `go mod` for dependencies
- run security scanners if you can

## known issues

none right now (hopefully)

if there are any, they'll be listed here with workarounds.

---

thanks for helping keep psx secure 🔒
```