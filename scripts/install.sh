#!/usr/bin/env bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
BINARY_NAME="psx"
REPO="m-mdy-m/psx"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"

echo -e "${BLUE}PSX Installation Script${NC}"
echo "======================="
echo ""

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        linux*)
            os="linux"
            ;;
        darwin*)
            os="darwin"
            ;;
        *)
            echo -e "${RED}Unsupported OS: $os${NC}"
            exit 1
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

add_to_path() {
    local dir="$1"
    
    # Detect shell
    local shell_rc=""
    if [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -f "$HOME/.profile" ]; then
        shell_rc="$HOME/.profile"
    fi
    
    if [ -z "$shell_rc" ]; then
        echo -e "${YELLOW}Could not detect shell config file${NC}"
        echo "Add this to your shell profile manually:"
        echo "  export PATH=\"${dir}:\$PATH\""
        return
    fi
    
    # Check if already in PATH
    if grep -q "export PATH.*${dir}" "$shell_rc" 2>/dev/null; then
        echo "Already in PATH config"
        return
    fi
    
    echo -e "${YELLOW}Adding ${dir} to PATH in ${shell_rc}...${NC}"
    echo "" >> "$shell_rc"
    echo "# PSX" >> "$shell_rc"
    echo "export PATH=\"${dir}:\$PATH\"" >> "$shell_rc"
    echo -e "${GREEN}✓ Added to PATH config${NC}"
    echo -e "${YELLOW}Run: source ${shell_rc}${NC}"
    echo "Or restart your terminal"
}

# Download and install from GitHub releases
install_from_github() {
    local platform=$(detect_platform)
    local version=${1:-latest}
    
    echo "Platform: ${platform}"
    echo "Version: ${version}"
    echo ""
    
    # Check for required tools
    if ! command_exists curl && ! command_exists wget; then
        echo -e "${RED}Error: curl or wget is required${NC}"
        exit 1
    fi
    
    # Get latest version if not specified
    if [ "$version" = "latest" ]; then
        echo -e "${YELLOW}Fetching latest version...${NC}"
        if command_exists curl; then
            version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        else
            version=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        fi
        
        if [ -z "$version" ]; then
            echo -e "${RED}Could not determine latest version${NC}"
            exit 1
        fi
        echo "Latest version: ${version}"
    fi
    
    # Build download URL
    local filename="${BINARY_NAME}-${platform}"
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    local tmp_file="/tmp/${BINARY_NAME}-${platform}"
    
    echo ""
    echo -e "${YELLOW}Downloading ${filename}...${NC}"
    
    if command_exists curl; then
        curl -L -o "${tmp_file}" "${url}"
    else
        wget -O "${tmp_file}" "${url}"
    fi
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Download failed${NC}"
        exit 1
    fi
    
    # Make executable
    chmod +x "${tmp_file}"
    
    # Determine install directory
    local install_to="${INSTALL_DIR}"
    local needs_path_update=false
    
    if [ ! -w "${INSTALL_DIR}" ]; then
        echo ""
        echo -e "${YELLOW}No write permission to ${INSTALL_DIR}${NC}"
        echo "Installing to user directory: ${USER_INSTALL_DIR}"
        install_to="${USER_INSTALL_DIR}"
        mkdir -p "${USER_INSTALL_DIR}"
        needs_path_update=true
    fi
    
    # Install
    echo ""
    echo -e "${YELLOW}Installing to ${install_to}...${NC}"
    
    if [ -w "${install_to}" ]; then
        mv "${tmp_file}" "${install_to}/${BINARY_NAME}"
    else
        sudo mv "${tmp_file}" "${install_to}/${BINARY_NAME}"
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ PSX installed successfully!${NC}"
        echo ""
        echo "Location: ${install_to}/${BINARY_NAME}"
        echo "Version: ${version}"
        
        # Check if in PATH
        if ! command_exists psx; then
            echo ""
            echo -e "${YELLOW}Warning: ${install_to} is not in your PATH${NC}"
            
            if [ "$needs_path_update" = true ]; then
                add_to_path "${install_to}"
            fi
        else
            echo ""
            echo -e "${GREEN}✓ psx command is available${NC}"
            echo "Try: psx --version"
        fi
    else
        echo -e "${RED}Installation failed${NC}"
        exit 1
    fi
}

# Install from local build
install_from_local() {
    local binary_path="${1:-build/${BINARY_NAME}}"
    
    if [ ! -f "${binary_path}" ]; then
        echo -e "${RED}Binary not found: ${binary_path}${NC}"
        echo "Build first: make build"
        exit 1
    fi
    
    echo "Installing from: ${binary_path}"
    
    local install_to="${INSTALL_DIR}"
    local needs_path_update=false
    
    if [ ! -w "${INSTALL_DIR}" ]; then
        echo ""
        echo -e "${YELLOW}No write permission to ${INSTALL_DIR}${NC}"
        echo "Installing to user directory: ${USER_INSTALL_DIR}"
        install_to="${USER_INSTALL_DIR}"
        mkdir -p "${USER_INSTALL_DIR}"
        needs_path_update=true
    fi
    
    echo -e "${YELLOW}Installing to ${install_to}...${NC}"
    
    if [ -w "${install_to}" ]; then
        cp "${binary_path}" "${install_to}/${BINARY_NAME}"
        chmod +x "${install_to}/${BINARY_NAME}"
    else
        sudo cp "${binary_path}" "${install_to}/${BINARY_NAME}"
        sudo chmod +x "${install_to}/${BINARY_NAME}"
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ PSX installed successfully!${NC}"
        echo ""
        echo "Location: ${install_to}/${BINARY_NAME}"
        
        "${install_to}/${BINARY_NAME}" --version
        
        if ! command_exists psx; then
            echo ""
            echo -e "${YELLOW}Warning: ${install_to} is not in your PATH${NC}"
            
            if [ "$needs_path_update" = true ]; then
                add_to_path "${install_to}"
            fi
        else
            echo ""
            echo -e "${GREEN}✓ psx command is available${NC}"
        fi
    else
        echo -e "${RED}Installation failed${NC}"
        exit 1
    fi
}

# Uninstall
uninstall() {
    local locations=(
        "${INSTALL_DIR}/${BINARY_NAME}"
        "${USER_INSTALL_DIR}/${BINARY_NAME}"
    )
    
    local found=0
    for location in "${locations[@]}"; do
        if [ -f "${location}" ]; then
            echo -e "${YELLOW}Removing ${location}...${NC}"
            
            if [ -w "$(dirname ${location})" ]; then
                rm "${location}"
            else
                sudo rm "${location}"
            fi
            
            found=1
            echo -e "${GREEN}✓ Removed${NC}"
        fi
    done
    
    if [ $found -eq 0 ]; then
        echo "PSX is not installed"
    else
        echo ""
        echo -e "${GREEN}PSX uninstalled successfully${NC}"
        echo ""
        echo "Note: PATH entries in shell configs remain."
        echo "Remove them manually if needed."
    fi
}

# Parse arguments
case "${1:-github}" in
    github)
        install_from_github "${2:-latest}"
        ;;
        
    local)
        install_from_local "${2}"
        ;;
        
    uninstall)
        uninstall
        ;;
        
    help|--help|-h)
        echo "Usage: $0 {github|local|uninstall} [options]"
        echo ""
        echo "Commands:"
        echo "  github [version]  - Install from GitHub releases (default: latest)"
        echo "  local [path]      - Install from local build (default: build/psx)"
        echo "  uninstall         - Remove PSX from system"
        echo ""
        echo "Examples:"
        echo "  $0 github              # Install latest from GitHub"
        echo "  $0 github v1.0.0       # Install specific version"
        echo "  $0 local               # Install from build/psx"
        echo "  $0 local /path/to/psx  # Install from custom path"
        echo "  $0 uninstall           # Remove PSX"
        ;;
        
    *)
        echo "Unknown command: $1"
        echo "Run '$0 help' for usage"
        exit 1
        ;;
esac