VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build info
BINARY_NAME := psx
BUILD_DIR := build
CMD_DIR := ./cmd/psx
 #===#

LDFLAGS= -ldflags "-X main.Version=$(VERSION)"

all: clean build
build:
	@echo "Building PSX $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "✓ Version: $(VERSION)"

dev:
	@echo "Running PSX in development mode..."
	@go run $(LDFLAGS) $(CMD_DIR)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

install: build
	@echo "Installing PSX to system..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

uninstall:
	@echo "Uninstalling PSX..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled"


# Build for multiple platforms
build-all:
	@echo "Building for all platforms(unix base system and windows)..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)

	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)

	@echo "Building for macOS (arm64/M1)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)

	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

	@echo "✓ All builds complete"

release: clean test build-all
	@echo "Creating release..."
	@cd $(BUILD_DIR) && sha256sum * > SHA256SUMS
	@echo "✓ Release ready in $(BUILD_DIR)/"

# watch:

check-deps:
	@echo "Checking dependencies..."
	@go mod verify
	@go mod tidy
	@echo "✓ Dependencies OK"

version:
	@echo "PSX Build Information"
	@echo "====================="
	@echo "Version:    $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"

help: ## Show this help message
	@echo "PSX Build System"
	@echo "================"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

