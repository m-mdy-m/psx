# Build and Installation Scripts

PSX build and installation scripts for Linux/macOS and Windows.

## Install

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

The installer downloads the platform release archive, extracts the binary, verifies its SHA-256 checksum against `checksums.txt`, installs it, and updates the appropriate user PATH file when needed.

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex
```

The PowerShell installer downloads `psx-windows-x64.exe`, verifies its SHA-256 checksum, installs it, and updates the user/system PATH. It uses `throw` instead of `exit`, so it is safe to run through `Invoke-Expression`.

## Build

```bash
./scripts/build.sh current
./scripts/build.sh all
./scripts/build.sh release
```

```powershell
.\scripts\build.ps1 current
.\scripts\build.ps1 all
.\scripts\build.ps1 release
```

Release platform names are standardized to:
- `linux-x64`
- `linux-arm64`
- `darwin-x64`
- `darwin-arm64`
- `windows-x64`
