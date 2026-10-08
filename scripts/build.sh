#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
BUILD_DATE="$(date -u '+%Y-%m-%d_%H:%M:%S')"
BINARY_NAME='psx'
BUILD_DIR='build'
CMD_DIR='./cmd/psx'

# Version is defined in internal/command, matching the release workflow.
LDFLAGS="-s -w -X github.com/m-mdy-m/psx/internal/command.Version=${VERSION}"

sha256_all() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$@"
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$@"
    else
        printf '%b\n' "${RED}sha256sum or shasum is required${NC}" >&2
        exit 1
    fi
}

build_platform() {
    local os="$1" arch="$2" output="$3"

    printf '%b\n' "${YELLOW}Building for ${os}/${arch}...${NC}"

    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
        -trimpath \
        -ldflags "$LDFLAGS" \
        -o "$output" \
        "$CMD_DIR"

    local size
    size="$(du -h "$output" | cut -f1)"
    printf '%b\n' "${GREEN}✓ Built: ${output} (${size})${NC}"
}

case "${1:-current}" in
    current)
        mkdir -p "$BUILD_DIR"
        build_platform "$(go env GOOS)" "$(go env GOARCH)" "$BUILD_DIR/$BINARY_NAME"
        ;;

    all)
        mkdir -p "$BUILD_DIR"
        build_platform linux amd64 "$BUILD_DIR/$BINARY_NAME-linux-x64"
        build_platform linux arm64 "$BUILD_DIR/$BINARY_NAME-linux-arm64"
        build_platform darwin amd64 "$BUILD_DIR/$BINARY_NAME-darwin-x64"
        build_platform darwin arm64 "$BUILD_DIR/$BINARY_NAME-darwin-arm64"
        build_platform windows amd64 "$BUILD_DIR/$BINARY_NAME-windows-x64.exe"
        ;;

    release)
        rm -rf "$BUILD_DIR"
        mkdir -p "$BUILD_DIR"

        printf '%b\n' "${YELLOW}Running tests...${NC}"
        go test ./...

        "$0" all

        cd "$BUILD_DIR"

        tar -czf "$BINARY_NAME-$VERSION-linux-x64.tar.gz" "$BINARY_NAME-linux-x64"
        tar -czf "$BINARY_NAME-$VERSION-linux-arm64.tar.gz" "$BINARY_NAME-linux-arm64"
        tar -czf "$BINARY_NAME-$VERSION-darwin-x64.tar.gz" "$BINARY_NAME-darwin-x64"
        tar -czf "$BINARY_NAME-$VERSION-darwin-arm64.tar.gz" "$BINARY_NAME-darwin-arm64"
        zip -q "$BINARY_NAME-$VERSION-windows-x64.zip" "$BINARY_NAME-windows-x64.exe"

        # Keep checksums for raw binaries, matching the GitHub release workflow.
        sha256_all \
            "$BINARY_NAME-linux-x64" \
            "$BINARY_NAME-linux-arm64" \
            "$BINARY_NAME-darwin-x64" \
            "$BINARY_NAME-darwin-arm64" \
            "$BINARY_NAME-windows-x64.exe" \
            > checksums.txt

        cd ..
        printf '%b\n' "${GREEN}✓ Release build complete${NC}"
        ;;

    clean)
        rm -rf "$BUILD_DIR"
        printf '%b\n' "${GREEN}✓ Clean complete${NC}"
        ;;

    *)
        echo "Usage: $0 {current|all|release|clean}"
        exit 1
        ;;
esac
