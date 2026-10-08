#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_NAME='psx'
REPO='m-mdy-m/psx'
INSTALL_DIR='/usr/local/bin'
USER_INSTALL_DIR="${HOME}/.local/bin"

say()  { printf '%b\n' "$*"; }
ok()   { say "${GREEN}$*${NC}"; }
warn() { say "${YELLOW}$*${NC}"; }
err()  { say "${RED}$*${NC}" >&2; }

command_exists() { command -v "$1" >/dev/null 2>&1; }

fail() {
    err "✗ $*"
    exit 1
}

sha256_file() {
    local file="$1"

    if command_exists sha256sum; then
        sha256sum "$file" | awk '{print $1}'
    elif command_exists shasum; then
        shasum -a 256 "$file" | awk '{print $1}'
    else
        fail 'sha256sum or shasum is required for release verification'
    fi
}

detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m | tr '[:upper:]' '[:lower:]')"

    case "$os" in
        linux*)  os='linux' ;;
        darwin*) os='darwin' ;;
        *) fail "Unsupported OS: $os" ;;
    esac

    case "$arch" in
        x86_64|amd64) arch='x64' ;;
        aarch64|arm64) arch='arm64' ;;
        *) fail "Unsupported architecture: $arch" ;;
    esac

    printf '%s-%s\n' "$os" "$arch"
}

get_latest_version() {
    local response version

    if command_exists curl; then
        response="$(curl -fsSL --retry 3 --retry-delay 1 "https://api.github.com/repos/${REPO}/releases/latest")" \
            || fail 'Could not fetch latest release information'
    else
        response="$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest")" \
            || fail 'Could not fetch latest release information'
    fi

    version="$(printf '%s\n' "$response" | sed -nE 's/^[[:space:]]*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*$/\1/p' | head -n1)"
    [ -n "$version" ] || fail 'Could not determine latest version'
    printf '%s\n' "$version"
}

normalize_version() {
    local version="$1"
    case "$version" in
        v*) printf '%s\n' "$version" ;;
        *)  printf 'v%s\n' "$version" ;;
    esac
}

download() {
    local url="$1" output="$2"

    if command_exists curl; then
        curl -fL --retry 3 --retry-delay 1 -o "$output" "$url"
    else
        wget -O "$output" "$url"
    fi
}

verify_binary_checksum() {
    local binary="$1" binary_name="$2" checksums="$3"
    local expected actual

    expected="$(awk -v file="$binary_name" '$2 == file {print $1; exit}' "$checksums")"
    [ -n "$expected" ] || fail "No checksum found for ${binary_name}"

    actual="$(sha256_file "$binary")"

    if [ "${actual,,}" != "${expected,,}" ]; then
        fail "Checksum verification failed for ${binary_name}"
    fi

    ok "✓ Checksum verified"
}

add_to_path() {
    local dir="$1"
    local shell_name shell_rc

    local shell_name="${SHELL##*/}"
    case "$shell_name" in
        zsh)  shell_rc="$HOME/.zshrc" ;;
        bash) shell_rc="$HOME/.bashrc" ;;
        fish) shell_rc="$HOME/.config/fish/config.fish" ;;
        *)    shell_name="profile"; shell_rc="$HOME/.profile" ;;
    esac

    case ":${PATH}:" in
        *":${dir}:"*)
            return
            ;;
    esac

    mkdir -p "$(dirname "$shell_rc")"

    if [ "$shell_name" = 'fish' ]; then
        if ! grep -Fq "$dir" "$shell_rc" 2>/dev/null; then
            {
                printf '\n# PSX\n'
                printf 'fish_add_path %s\n' "$dir"
            } >> "$shell_rc"
        fi
    else
        if ! grep -Fq "${dir}" "$shell_rc" 2>/dev/null; then
            {
                printf '\n# PSX\n'
                printf 'export PATH="%s:$PATH"\n' "$dir"
            } >> "$shell_rc"
        fi
    fi

    warn "${dir} was added to ${shell_rc}. Restart your shell (or source the file) to use psx."
}

install_binary() {
    local source="$1"
    local install_to="$INSTALL_DIR"
    local update_path=false

    if [ -d "$INSTALL_DIR" ] && [ -w "$INSTALL_DIR" ]; then
        install_to="$INSTALL_DIR"
    else
        install_to="$USER_INSTALL_DIR"
        update_path=true
        mkdir -p "$install_to"
        [ -w "$install_to" ] || fail "Cannot write to ${install_to}"
        warn "No write permission to ${INSTALL_DIR}; installing to ${install_to}"
    fi

    mkdir -p "$install_to"
    cp "$source" "${install_to}/${BINARY_NAME}"
    chmod 0755 "${install_to}/${BINARY_NAME}"

    ok '✓ PSX installed successfully!'
    say "Location: ${install_to}/${BINARY_NAME}"

    "${install_to}/${BINARY_NAME}" --version || fail 'Installed binary could not be executed'

    if [ "$update_path" = true ]; then
        add_to_path "$install_to"
    fi

    hash -r 2>/dev/null || true
}

install_from_github() {
    local platform version release_version archive_name binary_name
    local tmp_dir archive checksums extracted_binary checksum_binary

    platform="$(detect_platform)"
    version="${1:-latest}"

    say "Platform: ${platform}"
    say "Version: ${version}"
    say ''

    if ! command_exists curl && ! command_exists wget; then
        fail 'curl or wget is required'
    fi

    if [ "$version" = 'latest' ]; then
        warn 'Fetching latest version...'
        version="$(get_latest_version)"
    else
        version="$(normalize_version "$version")"
    fi

    release_version="$version"
    say "Release: ${release_version}"

    # Unix releases are distributed as tar.gz archives containing the raw binary.
    archive_name="psx-${release_version}-${platform}.tar.gz"
    binary_name="psx-${platform}"

    tmp_dir="$(mktemp -d)"
    trap 'rm -rf "$tmp_dir"' EXIT

    archive="${tmp_dir}/${archive_name}"
    checksums="${tmp_dir}/checksums.txt"
    extracted_binary="${tmp_dir}/${binary_name}"
    checksum_binary="${tmp_dir}/${binary_name}.checksum"

    warn "Downloading ${archive_name}..."
    download "https://github.com/${REPO}/releases/download/${release_version}/${archive_name}" "$archive" \
        || fail "Download failed: ${archive_name}"

    # A GitHub 404 can otherwise arrive as a tiny text file when a client follows a bad URL.
    [ -s "$archive" ] || fail 'Downloaded archive is empty'
    tar -tzf "$archive" >/dev/null 2>&1 || fail 'Downloaded release is not a valid gzip-compressed tar archive'

    warn 'Downloading checksums...'
    download "https://github.com/${REPO}/releases/download/${release_version}/checksums.txt" "$checksums" \
        || fail 'Could not download release checksums'

    tar -xzf "$archive" -C "$tmp_dir" \
        || fail 'Could not extract release archive'

    [ -f "$extracted_binary" ] || fail "Release archive does not contain ${binary_name}"
    chmod 0755 "$extracted_binary"

    # checksums.txt contains checksums for the raw binaries, not the archives.
    verify_binary_checksum "$extracted_binary" "$binary_name" "$checksums"

    install_binary "$extracted_binary"
}

install_from_local() {
    local binary_path="${1:-build/${BINARY_NAME}}"
    [ -f "$binary_path" ] || fail "Binary not found: ${binary_path}\nBuild first: make build"

    say "Installing from: ${binary_path}"
    install_binary "$binary_path"
}

uninstall() {
    local found=false location
    for location in "${INSTALL_DIR}/${BINARY_NAME}" "${USER_INSTALL_DIR}/${BINARY_NAME}"; do
        if [ -f "$location" ]; then
            rm -f "$location" 2>/dev/null || sudo rm -f "$location"
            ok "✓ Removed ${location}"
            found=true
        fi
    done

    if [ "$found" = false ]; then
        say 'PSX is not installed'
    else
        ok 'PSX uninstalled successfully'
        warn 'PATH entries are intentionally left untouched.'
    fi
}

say "${BLUE}PSX Installation Script${NC}"
say '======================='
say ''

case "${1:-github}" in
    github)    install_from_github "${2:-latest}" ;;
    local)     install_from_local "${2:-}" ;;
    uninstall) uninstall ;;
    help|--help|-h)
        cat <<'HELP'
Usage: install.sh {github|local|uninstall} [version/path]

Commands:
  github [version]  Install from GitHub releases (default: latest)
  local [path]     Install a local binary (default: build/psx)
  uninstall        Remove PSX

Examples:
  curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
  ./scripts/install.sh github v0.3.0
  ./scripts/install.sh local
HELP
        ;;
    *) fail "Unknown command: $1" ;;
esac
