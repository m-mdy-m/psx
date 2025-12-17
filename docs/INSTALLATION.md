# PSX Installation Guide

This guide covers all methods for installing PSX on your system.

## Table of Contents

- [Quick Install](#quick-install)
- [Download Binary](#download-binary)
- [Build from Source](#build-from-source)
- [Docker](#docker)
- [Uninstallation](#uninstallation)

---

## Quick Install

### Linux / macOS

**One-line installation:**
```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

**Install specific version:**
```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash -s github v1.0.0
```

**What it does:**
- Detects your platform (Linux/macOS, amd64/arm64)
- Downloads the appropriate binary from GitHub releases
- Installs to `/usr/local/bin` (system) or `~/.local/bin` (user)
- Adds to PATH automatically

### Windows

**PowerShell (Run as Administrator for system-wide install):**
```powershell
irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex
```

**Or download and run:**
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1" -OutFile install.ps1
.\install.ps1 github
```

**What it does:**
- Downloads binary for Windows amd64
- Installs to `C:\Program Files\PSX` (Admin) or `%LOCALAPPDATA%\PSX` (User)
- Adds to PATH automatically

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

PSX is available as a Docker image with multiple variants.

### Standard Image

```bash
# Pull image
docker pull bitsgenix/psx:latest

# Run in current directory
docker run --rm -v $(pwd):/project bitsgenix/psx:latest check

# Fix mode
docker run --rm -v $(pwd):/project bitsgenix/psx:latest fix --dry-run
```

### Available Tags

- `latest` - Latest stable release
- `v1.0.0` - Specific version
- `alpine` - Alpine-based (smaller)

### Using Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  psx:
    image: bitsgenix/psx:latest
    volumes:
      - .:/project
    working_dir: /project
    command: check
```

Run:
```bash
docker-compose run psx check
docker-compose run psx fix --interactive
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

### Permission Denied

**Issue:** Cannot write to `/usr/local/bin`

**Solution:**
```bash
# Install to user directory instead
mkdir -p ~/.local/bin
curl -L -o ~/.local/bin/psx https://github.com/m-mdy-m/psx/releases/download/v1.0.0/psx-linux-amd64
chmod +x ~/.local/bin/psx
export PATH="$HOME/.local/bin:$PATH"
```

### Build Errors

**Issue:** Build fails with Go errors

**Solution:**
```bash
# Verify Go version
go version  # Should be 1.25 or higher

# Update dependencies
go mod download
go mod tidy

# Clean and rebuild
make clean
make build
```

---

## Verify Installation

After installation, verify PSX is working:

```bash
# Check version
psx --version

# Run help
psx --help

# Test check command
psx check --help
```

Expected output:
```
Detected: generic

ERRORS (...)
...

Summary: X errors, Y warnings
Status: FAILED ✗
```

---

## Getting Help

- **Issues:** https://github.com/m-mdy-m/psx/issues
- **Discussions:** https://github.com/m-mdy-m/psx/discussions
- **Email:** bitsgenix@gmail.com