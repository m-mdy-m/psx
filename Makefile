VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build info
BINARY_NAME := psx
BUILD_DIR := build
CMD_DIR := ./cmd/psx

# LDFLAGS
LDFLAGS = -ldflags "-s -w -X github.com/m-mdy-m/psx/internal/command.Version=$(VERSION)"

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: all build dev clean install uninstall build-all release \
        test test-unit test-race test-coverage lint fmt vet check \
        docs tidy verify version docker docker-all help

all: clean build

build:
	@echo "$(YELLOW)Building psx $(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

dev:
	@echo "$(YELLOW)Building development version...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -race -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)Dev build complete$(NC)"

clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)âœ“ Clean complete$(NC)"

install: build
	@echo "$(YELLOW)Installing PSX...$(NC)"
	@./scripts/install.sh local $(BUILD_DIR)/$(BINARY_NAME)

uninstall:
	@echo "$(YELLOW)Uninstalling PSX...$(NC)"
	@./scripts/install.sh uninstall

build-all:
	@echo "$(YELLOW)Building for all platforms...$(NC)"
	@./scripts/build.sh all

# Release target - creates a new release
release:
	@echo "$(YELLOW)Creating release $(VERSION)...$(NC)"
	@echo ""
	@# Check if on main branch
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then \
		echo "$(RED)Error: Must be on main branch to release$(NC)"; \
		exit 1; \
	fi
	@# Check if working directory is clean
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)Error: Working directory is not clean$(NC)"; \
		exit 1; \
	fi
	@# Check if version tag exists
	@if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo "$(RED)Error: Tag $(VERSION) already exists$(NC)"; \
		exit 1; \
	fi
	@# Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test ./... -v
	@echo ""
	@# Build for all platforms
	@echo "$(YELLOW)Building for all platforms...$(NC)"
	@./scripts/build.sh release
	@echo ""
	@# Update CHANGELOG
	@echo "$(YELLOW)Update CHANGELOG.md with release notes...$(NC)"
	@echo "Press Enter when done..."
	@read dummy
	@# Commit and tag
	@git add CHANGELOG.md
	@git commit -m "chore: release $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo ""
	@echo "$(GREEN)âœ“ Release $(VERSION) created!$(NC)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review the changes"
	@echo "  2. Push: git push origin main --tags"
	@echo "  3. Create GitHub release with build/$(BINARY_NAME)-*.tar.gz"

test:
	@echo "Running tests..."
	go test ./...

test-unit:
	@echo "Running unit tests..."
	go test -short ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

vet:
	@echo "$(YELLOW)Running go vet...$(NC)"
	go vet ./...

lint:
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run
	@echo "$(GREEN)Lint passed$(NC)"

# check is what CI runs; keep it fast enough to run before every commit.
check: fmt vet test
	@echo "$(GREEN)All checks passed$(NC)"

# Regenerate docs/RULES.md from rules.yml.
docs:
	@echo "$(YELLOW)Regenerating documentation...$(NC)"
	go test ./internal/config/ -run TestGenerateRuleReference
	@echo "$(GREEN)Docs regenerated$(NC)"

tidy:
	@echo "$(YELLOW)Tidying modules...$(NC)"
	go mod tidy

verify: tidy vet test
	@if [ -n "$$(git diff --name-only go.mod go.sum)" ]; then \
		echo "$(RED)go.mod or go.sum is out of date; run make tidy$(NC)"; exit 1; \
	fi
	@if [ -n "$$(gofmt -l . | grep -v '^build/')" ]; then \
		echo "$(RED)Unformatted files:$(NC)"; gofmt -l . | grep -v '^build/'; exit 1; \
	fi
	@echo "$(GREEN)Verified$(NC)"

fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -s -w .
	@echo "$(GREEN)Code formatted$(NC)"

check-deps:
	@echo "$(YELLOW)Checking dependencies...$(NC)"
	@go mod verify
	@echo "$(GREEN)Dependencies verified$(NC)"

version:
	@echo "psx $(VERSION)"
	@echo "go $(shell go version)"

docker:
	@echo "$(YELLOW)Building Docker image (standard)...$(NC)"
	@docker buildx build  -t psx:latest -t psx:$(VERSION) -f Dockerfile .
	@echo "$(GREEN)âœ“ Docker image built: psx:latest$(NC)"

docker-alpine:
	@echo "$(YELLOW)Building Docker image (Alpine)...$(NC)"
	@docker buildx build  -t psx:alpine -t psx:$(VERSION)-alpine -f infra/Dockerfile.alpine .
	@echo "$(GREEN)âœ“ Docker image built: psx:alpine$(NC)"

docker-scratch:
	@echo "$(YELLOW)Building Docker image (Scratch)...$(NC)"
	@docker buildx build  -t psx:scratch -t psx:$(VERSION)-scratch -f infra/Dockerfile.scratch .
	@echo "$(GREEN)âœ“ Docker image built: psx:scratch$(NC)"

docker-all: docker docker-alpine docker-scratch
	@echo "$(GREEN)âœ“ All Docker images built$(NC)"

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
	@echo "psx build system"
	@echo ""
	@echo "Build:"
	@echo "  build          Build for the current platform"
	@echo "  dev            Build with the race detector"
	@echo "  build-all      Cross-compile for every release platform"
	@echo "  install        Install to the system"
	@echo "  uninstall      Remove from the system"
	@echo "  clean          Remove build artifacts"
	@echo ""
	@echo "Check:"
	@echo "  check          Format, vet and test (what CI runs)"
	@echo "  test           Run tests"
	@echo "  test-race      Run tests with the race detector"
	@echo "  test-coverage  Run tests and write coverage.html"
	@echo "  vet            Run go vet"
	@echo "  lint           Run golangci-lint"
	@echo "  fmt            Format the code"
	@echo "  verify         Check modules are tidy and the tree is formatted"
	@echo ""
	@echo "Docs:"
	@echo "  docs           Regenerate docs/RULES.md from rules.yml"
	@echo ""
	@echo "Docker:"
	@echo "  docker         Build the standard image"
	@echo "  docker-alpine  Build the Alpine image"
	@echo "  docker-scratch Build the scratch image"
	@echo "  docker-all     Build every image variant"
	@echo ""
	@echo "Other:"
	@echo "  tidy           Tidy go.mod"
	@echo "  version        Show version information"
	@echo "  release        Create a release (requires a clean main branch)"
	@echo "  help           Show this help"
