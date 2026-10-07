# Installation

One static binary, no runtime dependencies. Pick whichever route suits you.

- [Quick install](#quick-install)
- [Download a binary](#download-a-binary)
- [Build from source](#build-from-source)
- [Docker](#docker)
- [Uninstall](#uninstall)
- [Troubleshooting](#troubleshooting)

---

## Quick install

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

A specific version:

```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash -s github v3.0.0
```

The script detects your platform, downloads the matching binary, verifies it against the
release `checksums.txt`, installs to `/usr/local/bin` or `~/.local/bin`, and updates PATH.

### Windows

```powershell
irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex
```

Or download and run:

```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1" -OutFile install.ps1
.\install.ps1
```

Installs to `%LOCALAPPDATA%\psx` (user) or `C:\Program Files\psx` (administrator) and
updates PATH. Restart your shell afterwards.

### Verify

```bash
psx --version
psx check
```

---

## Download Binary

Download pre-compiled binaries from [GitHub Releases](https://github.com/m-mdy-m/psx/releases).

### Available Platforms

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | amd64 | psx-linux-amd64 |
| Linux | arm64 | psx-linux-arm64 |
| macOS | amd64 (Intel) | psx-darwin-amd64 |
| macOS | arm64 (M1/M2) | psx-darwin-arm64 |
| Windows | amd64 | psx-windows-amd64.exe |

### Manual Installation

**Linux/macOS:**
```bash
# Download (replace VERSION and PLATFORM)
VERSION=v1.0.0
PLATFORM=linux-amd64
curl -L -o psx https://github.com/m-mdy-m/psx/releases/download/${VERSION}/psx-${PLATFORM}

# Make executable
chmod +x psx

# Move to PATH
sudo mv psx /usr/local/bin/

# Or install to user directory
mkdir -p ~/.local/bin
mv psx ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"  # Add to ~/.bashrc or ~/.zshrc
```

**Windows:**
```powershell
# Download
$VERSION = "v1.0.0"
Invoke-WebRequest -Uri "https://github.com/m-mdy-m/psx/releases/download/${VERSION}/psx-windows-amd64.exe" -OutFile psx.exe

# Move to Program Files (requires Admin)
New-Item -ItemType Directory -Force -Path "C:\Program Files\PSX"
Move-Item psx.exe "C:\Program Files\PSX\psx.exe"

# Add to PATH (System, requires Admin)
$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
[Environment]::SetEnvironmentVariable("Path", "$path;C:\Program Files\PSX", "Machine")

# Or install to user directory (no Admin needed)
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\PSX"
Move-Item psx.exe "$env:LOCALAPPDATA\PSX\psx.exe"
$path = [Environment]::GetEnvironmentVariable("Path", "User")
[Environment]::SetEnvironmentVariable("Path", "$path;$env:LOCALAPPDATA\PSX", "User")
```

### Verify Installation

```bash
psx --version
```

---

## Build from Source

### Prerequisites

- **Go:** 1.25 or higher
- **Git:** For cloning the repository
- **Make:** For using Makefile commands

### Clone and Build

```bash
# Clone repository
git clone https://github.com/m-mdy-m/psx.git && cd psx

# Build using Make
make build

# Install to system
sudo make install
```

### Build for Specific Platform

```bash
# Current platform
make build

# All platforms
make build-all
```

### Build Options

The Makefile provides several targets:

```bash
make build          # Build for current platform
make dev            # Build with race detector
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linters
make clean          # Remove build artifacts
```

### Development Build

```bash
# Build with debug info
make dev

# Run without installing
./build/psx check
```

---

## Docker

Images are published on Docker Hub as `bitsgenix/psx`.

```bash
docker pull bitsgenix/psx:latest

# Check the current directory
docker run --rm -v "$(pwd):/project" -w /project bitsgenix/psx:latest check

# Preview a fix
docker run --rm -v "$(pwd):/project" -w /project bitsgenix/psx:latest fix --dry-run

# Apply a fix (writes into the mounted directory)
docker run --rm -v "$(pwd):/project" -w /project bitsgenix/psx:latest fix --yes
```

`check` needs no write access. `fix` does, so only mount read-only when you are checking.

### Tags

| Tag | Image | Size |
| --- | --- | --- |
| `latest` | Debian | ~20 MB |
| `alpine` | Alpine | ~17 MB |
| `scratch` | No libc | ~5 MB |

The scratch image has no shell. Only the entrypoint works; do not run `sh` inside it.

### Compose

```yaml
services:
  psx:
    image: bitsgenix/psx:latest
    read_only: true
    volumes:
      - .:/project:ro
    working_dir: /project
    command: ["check", "--output", "github"]
```

The obsolete top-level `version` key is omitted; Compose v2 does not want it.

To use `fix`, drop `:ro`:

```bash
docker compose run --rm psx fix --yes
```

### Build Docker Image Locally

```bash
# Standard image
docker build -t psx:local .

# Alpine variant
docker build -t psx:alpine -f infra/Dockerfile.alpine .

# Minimal (scratch) variant
docker build -t psx:scratch -f infra/Dockerfile.scratch .
```

---

## Uninstallation

### Using Install Scripts

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash -s uninstall
```

**Windows:**
```powershell
.\install.ps1 uninstall
```

### Manual Removal

**Linux/macOS:**
```bash
# Remove binary
sudo rm /usr/local/bin/psx
# Or from user install
rm ~/.local/bin/psx

# Remove config (optional)
rm -rf ~/.config/psx
```

**Windows:**
```powershell
# Remove binary (System install)
Remove-Item "C:\Program Files\PSX\psx.exe"
Remove-Item "C:\Program Files\PSX"

# Or User install
Remove-Item "$env:LOCALAPPDATA\PSX\psx.exe"
Remove-Item "$env:LOCALAPPDATA\PSX"

# Remove from PATH manually if needed
```

**Docker:**
```bash
docker rmi bitsgenix/psx:latest
```

---

## Troubleshooting

### Command Not Found After Installation

**Issue:** `psx: command not found`

**Solution:**

**Linux/macOS:**
```bash
# Check if binary exists
ls -l /usr/local/bin/psx
# or
ls -l ~/.local/bin/psx

# Verify PATH
echo $PATH

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Reload shell config
source ~/.bashrc  # or source ~/.zshrc
```

**Windows:**
```powershell
# Check if binary exists
Test-Path "C:\Program Files\PSX\psx.exe"

# Verify PATH
$env:Path

# Restart terminal after installation
```

### Permission denied

Cannot write to `/usr/local/bin`. Install to your user directory instead:

```bash
mkdir -p ~/.local/bin
curl -fsSL -o ~/.local/bin/psx \
  https://github.com/m-mdy-m/psx/releases/latest/download/psx-linux-amd64
chmod +x ~/.local/bin/psx
export PATH="$HOME/.local/bin:$PATH"
```

Add that `export` to your shell profile to make it permanent.

### The installer downloaded a web page instead of a binary

An old `install.sh` used `curl -L` without `-f`, so a 404 wrote HTML into the binary. Pull
the current script:

```bash
curl -fsSL -O https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh
bash install.sh
```

### Build errors

```bash
go version          # must be 1.25+
go mod download
make clean && make build
```

If `make` is unavailable or you are on Windows:

```bash
go build -o build/psx ./cmd/psx
```

---

## Verify

```bash
psx --version
psx check
```

`check` should print grouped findings and a status line, and exit 0 or 1:

```
Errors (1)

  error  tests_folder
        No tests found
        fix: psx fix --rule tests_folder

Result: 1 errors, 0 warnings
Status: FAILED
```

If it hangs waiting for input, something is wrong: `check` never prompts. Run
`psx check --verbose` and look for a prompt, or pass `--yes`.

---

## Next steps

```bash
psx init       # write a psx.yml tuned to this project
psx check      # see what is missing
psx fix        # create it
```

See [CONFIGURATION.md](CONFIGURATION.md) for the options.

---

## Getting help

- Issues: https://github.com/m-mdy-m/psx/issues
- Discussions: https://github.com/m-mdy-m/psx/discussions
- Email: bitsgenix@gmail.com