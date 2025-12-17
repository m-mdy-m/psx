# Build and Installation Scripts

This folder contains scripts for building and installing PSX.

## Quick Start

### Unix/Linux/macOS

```bash
# Build
./scripts/build.sh

# Install from local build
./scripts/install.sh local

# Or install from GitHub
./scripts/install.sh github
```

### Windows

```powershell
# Build
.\scripts\build.ps1

# Install from local build
.\scripts\install.ps1 local

# Or install from GitHub
.\scripts\install.ps1 github
```

---

## Build Scripts

### Unix: `build.sh`

Builds PSX binaries for various platforms.

**Usage:**
```bash
./scripts/build.sh [command]
```

**Commands:**
- `current` - Build for current platform (default)
- `all` - Build for all supported platforms
- `release` - Create release packages with checksums
- `clean` - Remove build directory

**Examples:**
```bash
# Build for current system
./scripts/build.sh

# Build for all platforms
./scripts/build.sh all

# Create release packages
./scripts/build.sh release

# Clean build artifacts
./scripts/build.sh clean
```

**Output:**
- Binaries in `build/` directory
- Release builds include SHA256 checksums

---

### Windows: `build.ps1`

PowerShell equivalent of the Unix build script.

**Usage:**
```powershell
.\scripts\build.ps1 [command]
```

**Commands:**
Same as Unix version:
- `current`
- `all`
- `release`
- `clean`

**Examples:**
```powershell
# Build for Windows
.\scripts\build.ps1

# Build for all platforms
.\scripts\build.ps1 -Command all

# Create release
.\scripts\build.ps1 -Command release
```

---

## Installation Scripts

### Unix: `install.sh`

Installs PSX binary to system or user directory.

**Usage:**
```bash
./scripts/install.sh [command] [options]
```

**Commands:**
- `github [version]` - Install from GitHub releases (default)
- `local [path]` - Install from local build
- `uninstall` - Remove PSX from system

**Install Locations:**
- System (with sudo): `/usr/local/bin`
- User (without sudo): `~/.local/bin`

**Examples:**
```bash
# Install latest from GitHub
./scripts/install.sh github

# Install specific version
./scripts/install.sh github v1.0.0

# Install from local build
./scripts/install.sh local

# Install from custom path
./scripts/install.sh local /path/to/psx

# Uninstall
./scripts/install.sh uninstall
```

**Features:**
- Auto-detects platform (Linux/macOS, amd64/arm64)
- Automatically adds to PATH
- Works without sudo (user install)
- Downloads from GitHub releases
- Verifies installation

---

### Windows: `install.ps1`

PowerShell installation script for Windows.

**Usage:**
```powershell
.\scripts\install.ps1 [command] [options]
```

**Commands:**
Same as Unix version:
- `github [version]`
- `local [path]`
- `uninstall`

**Install Locations:**
- System (Admin): `C:\Program Files\PSX`
- User (Non-admin): `%LOCALAPPDATA%\PSX`

**Examples:**
```powershell
# Install latest from GitHub
.\scripts\install.ps1 github

# Install specific version
.\scripts\install.ps1 github v1.0.0

# Install from build
.\scripts\install.ps1 local

# Install from custom path
.\scripts\install.ps1 local C:\path\to\psx.exe

# Uninstall
.\scripts\install.ps1 uninstall
```

**Features:**
- Detects Administrator privileges
- Automatically adds to PATH
- Works for both user and system installs
- Downloads from GitHub releases