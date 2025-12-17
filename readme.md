# PSX - Project Structure Checker

**PSX** is a command-line tool that validates and standardizes project structures across different programming languages. It helps maintain consistency in your projects by checking for essential files, proper folder organization, and development best practices.

## Features

- **Automatic Detection** - Identifies project type (Node.js, Go, etc.)
- **Auto-Fix** - Automatically creates missing files and folders
- **Multi-Language Support** - Supports Node.js, Go, and generic projects
- **Configurable** - Customize rules and severity levels via YAML
- **Fast** - Parallel rule execution for quick validation

## Installation

### Quick Install (Recommended)

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1" -OutFile install.ps1; .\install.ps1 github
```

### Download Binary

Download pre-built binaries from [Releases](https://github.com/m-mdy-m/psx/releases):

- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)

### Build from Source

Requirements: Go 1.25+

```bash
git clone https://github.com/m-mdy-m/psx
cd psx
make build
sudo make install
```

### Docker

```bash
docker pull bitsgenix/psx:latest
docker run --rm -v $(pwd):/project psx:latest check
```

See [INSTALLATION.md](docs/INSTALLATION.md) for detailed installation instructions.

## Quick Start

**Check your project:**
```bash
cd my-project
psx check
```

**Fix issues automatically:**
```bash
psx fix
```

**Fix with confirmation:**
```bash
psx fix --interactive
```

## Development

### Requirements

- Go 1.25+
- Make

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE) - Copyright (c) 2024 m-mdy-m

## Links

- **Repository:** https://github.com/m-mdy-m/psx
- **Issues:** https://github.com/m-mdy-m/psx/issues
- **Releases:** https://github.com/m-mdy-m/psx/releases
- **Documentation:** [docs/](docs/)

## Support

- Email: bitsgenix@gmail.com
- GitHub Discussions: [Discussions](https://github.com/m-mdy-m/psx/discussions)
- Bug Reports: [Issues](https://github.com/m-mdy-m/psx/issues)