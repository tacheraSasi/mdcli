# mdcli Makefile
# Advanced Markdown CLI Processor

VERSION := 2.0.0
BINARY_NAME := mdcli
PACKAGE := github.com/tacheraSasi/mdcli
MAIN_PACKAGE := .

# Build information
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Linker flags for version information
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

# Directories
BIN_DIR := bin
DIST_DIR := dist

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

.PHONY: help build build_all clean test lint fmt vet deps install uninstall release dev

# Default target
help: ## Show this help message
	@echo "$(CYAN)mdcli v$(VERSION) - Makefile Help$(NC)"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "$(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
dev: ## Build development version
	@echo "$(BLUE)Building development version...$(NC)"
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./...

bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	go test -bench=. ./...

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; exit 1; }
	golangci-lint run

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	go vet ./...

deps: ## Download and tidy dependencies
	@echo "$(BLUE)Managing dependencies...$(NC)"
	go mod download
	go mod tidy
	go mod verify

# Build targets
build: dev ## Build for current platform

build_all: build_linux build_mac build_windows build_android ## Build for all platforms

build_linux: ## Build for Linux AMD64
	@echo "$(GREEN)Building Linux AMD64...$(NC)"
	@mkdir -p $(BIN_DIR)
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)_linux_amd64
	@if [ -f "./upx" ]; then \
		echo "$(YELLOW)Compressing binary...$(NC)"; \
		./upx --best --lzma $(BIN_DIR)/$(BINARY_NAME)_linux_amd64; \
	fi
	@echo "$(YELLOW)Creating archive...$(NC)"
	tar -czf $(BIN_DIR)/$(BINARY_NAME)_linux_amd64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)_linux_amd64
	@echo "$(GREEN)✅ Linux build complete: $(BIN_DIR)/$(BINARY_NAME)_linux_amd64.tar.gz$(NC)"

build_mac: ## Build for macOS AMD64
	@echo "$(GREEN)Building macOS AMD64...$(NC)"
	@mkdir -p $(BIN_DIR)
	env GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64
	@if [ -f "./upx" ]; then \
		echo "$(YELLOW)Compressing binary...$(NC)"; \
		./upx --best --lzma $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64; \
	fi
	@echo "$(YELLOW)Creating archive...$(NC)"
	tar -czf $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)_darwin_amd64
	@echo "$(GREEN)✅ macOS build complete: $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64.tar.gz$(NC)"

build_mac_arm: ## Build for macOS ARM64 (Apple Silicon)
	@echo "$(GREEN)Building macOS ARM64...$(NC)"
	@mkdir -p $(BIN_DIR)
	env GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)_darwin_arm64
	@echo "$(YELLOW)Creating archive...$(NC)"
	tar -czf $(BIN_DIR)/$(BINARY_NAME)_darwin_arm64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)_darwin_arm64
	@echo "$(GREEN)✅ macOS ARM64 build complete: $(BIN_DIR)/$(BINARY_NAME)_darwin_arm64.tar.gz$(NC)"

build_windows: ## Build for Windows AMD64
	@echo "$(GREEN)Building Windows AMD64...$(NC)"
	@mkdir -p $(BIN_DIR)
	env GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)_windows_amd64.exe
	@if [ -f "./upx" ]; then \
		echo "$(YELLOW)Compressing binary...$(NC)"; \
		./upx --best --lzma $(BIN_DIR)/$(BINARY_NAME)_windows_amd64.exe; \
	fi
	@echo "$(YELLOW)Creating archive...$(NC)"
	cd $(BIN_DIR) && zip $(BINARY_NAME)_windows_amd64.zip $(BINARY_NAME)_windows_amd64.exe
	@echo "$(GREEN)✅ Windows build complete: $(BIN_DIR)/$(BINARY_NAME)_windows_amd64.zip$(NC)"

build_android: ## Build for Android ARM64
	@echo "$(GREEN)Building Android ARM64...$(NC)"
	@mkdir -p $(BIN_DIR)
	env GOOS=android GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)_android_arm64
	@echo "$(YELLOW)Creating archive...$(NC)"
	tar -czf $(BIN_DIR)/$(BINARY_NAME)_android_arm64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)_android_arm64
	@echo "$(GREEN)✅ Android build complete: $(BIN_DIR)/$(BINARY_NAME)_android_arm64.tar.gz$(NC)"

# Install/Uninstall
install: build ## Install binary to $GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	go install -ldflags="$(LDFLAGS)"
	@echo "$(GREEN)✅ $(BINARY_NAME) installed successfully$(NC)"

uninstall: ## Remove binary from $GOPATH/bin
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ $(BINARY_NAME) uninstalled$(NC)"

# Release targets
release: clean build_all ## Create release builds
	@echo "$(PURPLE)Creating release v$(VERSION)...$(NC)"
	@mkdir -p $(DIST_DIR)
	@cp $(BIN_DIR)/*.tar.gz $(BIN_DIR)/*.zip $(DIST_DIR)/ 2>/dev/null || true
	@echo "$(GREEN)✅ Release builds created in $(DIST_DIR)/$(NC)"

checksums: ## Generate checksums for release files
	@echo "$(BLUE)Generating checksums...$(NC)"
	@cd $(DIST_DIR) && find . -type f \( -name "*.tar.gz" -o -name "*.zip" \) -exec shasum -a 256 {} \; > checksums.txt
	@echo "$(GREEN)✅ Checksums generated: $(DIST_DIR)/checksums.txt$(NC)"

# Utility targets
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf $(BIN_DIR) $(DIST_DIR)
	rm -f $(BINARY_NAME)
	go clean -cache -testcache -modcache
	@echo "$(GREEN)✅ Clean complete$(NC)"

size: build ## Show binary size information
	@echo "$(CYAN)Binary size information:$(NC)"
	@ls -lh $(BINARY_NAME) | awk '{print "Size: " $$5}'
	@file $(BINARY_NAME)

run: build ## Build and run with test file
	@echo "$(BLUE)Building and running with test.md...$(NC)"
	./$(BINARY_NAME) render test.md

demo: build ## Run demo commands
	@echo "$(PURPLE)mdcli Demo$(NC)"
	@echo "============"
	./$(BINARY_NAME) --version
	@echo ""
	./$(BINARY_NAME) themes
	@echo ""
	@echo "$(CYAN)Rendering test.md:$(NC)"
	./$(BINARY_NAME) render test.md

# Docker targets (if needed)
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run: docker-build ## Run in Docker container
	@echo "$(BLUE)Running in Docker...$(NC)"
	docker run --rm -v $(PWD):/workspace $(BINARY_NAME):$(VERSION)

# Info targets
info: ## Show build information
	@echo "$(CYAN)Build Information$(NC)"
	@echo "=================="
	@echo "Version:     $(VERSION)"
	@echo "Package:     $(PACKAGE)"
	@echo "Build Time:  $(BUILD_TIME)"
	@echo "Git Commit:  $(GIT_COMMIT)"
	@echo "Git Branch:  $(GIT_BRANCH)"
	@echo "Go Version:  $(shell go version)"
	@echo "Platform:    $(shell go env GOOS)/$(shell go env GOARCH)"

deps-update: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

# Legacy compatibility
build_linux: build_linux
build_mac: build_mac  
build_windows: build_windows
build_android: build_android
dependencies: deps
	go clean
