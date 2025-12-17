VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build info
BINARY_NAME := psx
BUILD_DIR := build
CMD_DIR := ./cmd/psx

# LDFLAGS
LDFLAGS = -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)"

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m 

all: clean build

build:
	@echo "$(YELLOW)Building PSX $(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@echo "$(GREEN)✓ Version: $(VERSION)$(NC)"

dev:
	@echo "$(YELLOW)Building development version...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -race -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)✓ Dev build complete$(NC)"

clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ Clean complete$(NC)"

install: build
	@echo "$(YELLOW)Installing PSX...$(NC)"
	@./scripts/install.sh local $(BUILD_DIR)/$(BINARY_NAME)

uninstall:
	@echo "$(YELLOW)Uninstalling PSX...$(NC)"
	@./scripts/install.sh uninstall

build-all:
	@echo "$(YELLOW)Building for all platforms...$(NC)"
	@./scripts/build.sh all

release:
	@echo "$(YELLOW)Creating release packages...$(NC)"
	@./scripts/build.sh release

test:
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	@golangci-lint run --timeout=5m
	@echo "$(GREEN)✓ Lint passed$(NC)"

fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -s -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

check-deps:
	@echo "$(YELLOW)Checking dependencies...$(NC)"
	@go mod verify
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies OK$(NC)"

version:
	@echo "PSX Build Information"
	@echo "====================="
	@echo "Version:    $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Go Version: $(shell go version)"

docker:
	@echo "$(YELLOW)Building Docker image (standard)...$(NC)"
	@docker build -t psx:latest -t psx:$(VERSION) -f Dockerfile .
	@echo "$(GREEN)✓ Docker image built: psx:latest$(NC)"

docker-alpine:
	@echo "$(YELLOW)Building Docker image (Alpine)...$(NC)"
	@docker build -t psx:alpine -t psx:$(VERSION)-alpine -f infra/Dockerfile.alpine .
	@echo "$(GREEN)✓ Docker image built: psx:alpine$(NC)"

docker-scratch:
	@echo "$(YELLOW)Building Docker image (Scratch)...$(NC)"
	@docker build -t psx:scratch -t psx:$(VERSION)-scratch -f infra/Dockerfile.scratch .
	@echo "$(GREEN)✓ Docker image built: psx:scratch$(NC)"

docker-all: docker docker-alpine docker-scratch
	@echo "$(GREEN)✓ All Docker images built$(NC)"

docker-run:
	@docker run --rm -v $(PWD):/project psx:latest check

docker-run-alpine:
	@docker run --rm -v $(PWD):/project psx:alpine check

docker-run-scratch:
	@docker run --rm -v $(PWD):/project psx:scratch check

docker-compose-build:
	@docker-compose build

docker-compose-up:
	@docker-compose up psx

help:
	@echo "PSX Build System"
	@echo "================"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main targets:"
	@echo "  build         - Build for current platform"
	@echo "  dev           - Build development version (with race detector)"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install to system"
	@echo "  uninstall     - Remove from system"
	@echo ""
	@echo "Cross-platform:"
	@echo "  build-all     - Build for all platforms"
	@echo "  release       - Create release packages"
	@echo ""
	@echo "Testing:"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  benchmark     - Run benchmarks"
	@echo "  lint          - Run linters"
	@echo "  fmt           - Format code"
	@echo ""
	@echo "Docker:"
	@echo "  docker        - Build standard Docker image"
	@echo "  docker-alpine - Build Alpine Docker image"
	@echo "  docker-scratch- Build Scratch Docker image"
	@echo "  docker-all    - Build all Docker images"
	@echo ""
	@echo "Utilities:"
	@echo "  check-deps    - Verify and tidy dependencies"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build              # Build binary"
	@echo "  make test               # Run tests"
	@echo "  make docker-all         # Build all Docker images"
	@echo "  make release            # Create release packages"