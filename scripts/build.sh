#!/usr/bin/env bash
# Build script for PSX
set -e
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
BINARY_NAME="psx"
BUILD_DIR="build"
CMD_DIR="./cmd/psx"

# Go build flags
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"

echo -e "${BLUE}PSX Build Script${NC}"
echo "================="
echo "Version: ${VERSION}"
echo "Date: ${BUILD_DATE}"
echo ""

# Function to build for specific platform
build_platform() {
    local os=$1
    local arch=$2
    local output=$3
    
    echo -e "${YELLOW}Building for ${os}/${arch}...${NC}"
    
    GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o "${output}" \
        ${CMD_DIR}
    
    if [ $? -eq 0 ]; then
        local size=$(du -h "${output}" | cut -f1)
        echo -e "${GREEN}✓ Built: ${output} (${size})${NC}"
    else
        echo -e "${RED}✗ Build failed for ${os}/${arch}${NC}"
        return 1
    fi
}

# Parse arguments
case "${1:-current}" in
    current)
        # Build for current platform
        echo "Building for current platform..."
        mkdir -p ${BUILD_DIR}
        build_platform $(go env GOOS) $(go env GOARCH) "${BUILD_DIR}/${BINARY_NAME}"
        echo ""
        echo -e "${GREEN}Build complete!${NC}"
        echo "Binary: ${BUILD_DIR}/${BINARY_NAME}"
        ;;
        
    all)
        # Build for all platforms
        echo "Building for all platforms..."
        mkdir -p ${BUILD_DIR}
        
        # Linux
        build_platform linux amd64 "${BUILD_DIR}/${BINARY_NAME}-linux-amd64"
        build_platform linux arm64 "${BUILD_DIR}/${BINARY_NAME}-linux-arm64"
        
        # macOS
        build_platform darwin amd64 "${BUILD_DIR}/${BINARY_NAME}-darwin-amd64"
        build_platform darwin arm64 "${BUILD_DIR}/${BINARY_NAME}-darwin-arm64"
        
        # Windows
        build_platform windows amd64 "${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe"
        
        # Create checksums
        echo ""
        echo -e "${YELLOW}Generating checksums...${NC}"
        cd ${BUILD_DIR}
        sha256sum * > SHA256SUMS
        cd ..
        
        echo ""
        echo -e "${GREEN}All builds complete!${NC}"
        echo "Binaries in: ${BUILD_DIR}/"
        ;;
        
    release)
        # Build for release (all platforms + checksums)
        echo "Building release..."
        
        # Clean first
        rm -rf ${BUILD_DIR}
        mkdir -p ${BUILD_DIR}
        
        # Run tests first
        echo -e "${YELLOW}Running tests...${NC}"
        go test ./... -v
        if [ $? -ne 0 ]; then
            echo -e "${RED}Tests failed! Aborting release build.${NC}"
            exit 1
        fi
        
        # Build all platforms
        $0 all
        
        # Create archives
        echo ""
        echo -e "${YELLOW}Creating archives...${NC}"
        cd ${BUILD_DIR}
        
        # Linux amd64
        tar czf ${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
        
        # Linux arm64
        tar czf ${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
        
        # macOS amd64
        tar czf ${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
        
        # macOS arm64
        tar czf ${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
        
        # Windows
        zip -q ${BINARY_NAME}-${VERSION}-windows-amd64.zip ${BINARY_NAME}-windows-amd64.exe
        
        # Update checksums
        sha256sum *.tar.gz *.zip > SHA256SUMS
        
        cd ..
        
        echo ""
        echo -e "${GREEN}Release build complete!${NC}"
        echo "Archives in: ${BUILD_DIR}/"
        ;;
        
    clean)
        echo "Cleaning build directory..."
        rm -rf ${BUILD_DIR}
        echo -e "${GREEN}Clean complete!${NC}"
        ;;
        
    *)
        echo "Usage: $0 {current|all|release|clean}"
        echo ""
        echo "Commands:"
        echo "  current  - Build for current platform (default)"
        echo "  all      - Build for all supported platforms"
        echo "  release  - Build release packages with checksums"
        echo "  clean    - Remove build directory"
        exit 1
        ;;
esac