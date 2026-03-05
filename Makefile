# mdcli Makefile
# Advanced Markdown CLI Processor

.PHONY: run build dev build-all build-linux build-darwin build-windows build-android clean test lint fmt vet deps deps-update install uninstall release checksums info demo

# ── Variables ──────────────────────────────────────────────────────────────────
APP_NAME    := mdcli
MODULE      := github.com/tacheraSasi/mdcli
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS     := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.buildTime=$(BUILD_DATE)' \
	-X 'main.gitCommit=$(COMMIT)'
BUILD_DIR   := ./build
CGO_ENABLED ?= 0

# ── Development ───────────────────────────────────────────────────────────────
run: build ## Build and run with test file
	./$(APP_NAME) render test.md

build: deps ## Build for current platform
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)

dev: ## Run directly without producing a binary
	go run . render test.md

test: ## Run all tests
	go test -v ./...

bench: ## Run benchmarks
	go test -bench=. ./...

lint: ## Run linter
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# ── Cross-platform builds (all) ──────────────────────────────────────────────
build-all: deps ## Build for all platforms
	@echo "==> Building for all platforms ($(VERSION))..."
	mkdir -p $(BUILD_DIR)
	# macOS Apple Silicon (M1/M2/M3/M4)
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64      .
	# macOS Intel
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64      .
	# Linux amd64
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64       .
	# Linux arm64
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64       .
	# Windows amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	# Windows arm64
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .
	# Android arm64
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-android-arm64     .
	@echo "==> Packaging archives..."
	cd $(BUILD_DIR) && cp $(APP_NAME)-darwin-arm64  $(APP_NAME) && tar czf $(APP_NAME)-darwin-arm64.tar.gz  $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-darwin-amd64  $(APP_NAME) && tar czf $(APP_NAME)-darwin-amd64.tar.gz  $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-linux-amd64   $(APP_NAME) && tar czf $(APP_NAME)-linux-amd64.tar.gz   $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-linux-arm64   $(APP_NAME) && tar czf $(APP_NAME)-linux-arm64.tar.gz   $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-android-arm64 $(APP_NAME) && tar czf $(APP_NAME)-android-arm64.tar.gz $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-windows-amd64.exe $(APP_NAME).exe && zip $(APP_NAME)-windows-amd64.zip $(APP_NAME).exe && rm $(APP_NAME).exe
	cd $(BUILD_DIR) && cp $(APP_NAME)-windows-arm64.exe $(APP_NAME).exe && zip $(APP_NAME)-windows-arm64.zip $(APP_NAME).exe && rm $(APP_NAME).exe
	@echo "==> Done! Binaries in $(BUILD_DIR)/"

# ── Individual platform builds ───────────────────────────────────────────────
build-darwin: deps ## Build for macOS (arm64 + amd64)
	@echo "==> Building for macOS..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	cd $(BUILD_DIR) && cp $(APP_NAME)-darwin-arm64 $(APP_NAME) && tar czf $(APP_NAME)-darwin-arm64.tar.gz $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-darwin-amd64 $(APP_NAME) && tar czf $(APP_NAME)-darwin-amd64.tar.gz $(APP_NAME) && rm $(APP_NAME)

build-linux: deps ## Build for Linux (amd64 + arm64)
	@echo "==> Building for Linux..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	cd $(BUILD_DIR) && cp $(APP_NAME)-linux-amd64 $(APP_NAME) && tar czf $(APP_NAME)-linux-amd64.tar.gz $(APP_NAME) && rm $(APP_NAME)
	cd $(BUILD_DIR) && cp $(APP_NAME)-linux-arm64 $(APP_NAME) && tar czf $(APP_NAME)-linux-arm64.tar.gz $(APP_NAME) && rm $(APP_NAME)

build-windows: deps ## Build for Windows (amd64 + arm64)
	@echo "==> Building for Windows..."
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .
	cd $(BUILD_DIR) && cp $(APP_NAME)-windows-amd64.exe $(APP_NAME).exe && zip $(APP_NAME)-windows-amd64.zip $(APP_NAME).exe && rm $(APP_NAME).exe
	cd $(BUILD_DIR) && cp $(APP_NAME)-windows-arm64.exe $(APP_NAME).exe && zip $(APP_NAME)-windows-arm64.zip $(APP_NAME).exe && rm $(APP_NAME).exe

build-android: deps ## Build for Android (arm64)
	@echo "==> Building for Android..."
	mkdir -p $(BUILD_DIR)
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-android-arm64 .
	cd $(BUILD_DIR) && cp $(APP_NAME)-android-arm64 $(APP_NAME) && tar czf $(APP_NAME)-android-arm64.tar.gz $(APP_NAME) && rm $(APP_NAME)

# ── SHA256 checksums ─────────────────────────────────────────────────────────
checksums: ## Generate checksums for release files
	@echo "==> Generating checksums..."
	cd $(BUILD_DIR) && shasum -a 256 *.tar.gz *.zip > checksums.txt
	@cat $(BUILD_DIR)/checksums.txt

# ── GitHub Release (requires gh CLI) ─────────────────────────────────────────
release: clean build-all checksums ## Create and push a GitHub release with binaries
	@echo "==> Creating GitHub release $(VERSION)..."
	gh release create $(VERSION) \
		--title "mdcli $(VERSION)" \
		--generate-notes \
		$(BUILD_DIR)/*.tar.gz \
		$(BUILD_DIR)/*.zip \
		$(BUILD_DIR)/checksums.txt

# ── Install / Uninstall ──────────────────────────────────────────────────────
install: build ## Install binary to $$GOPATH/bin
	go install -ldflags "$(LDFLAGS)"

uninstall: ## Remove binary from $$GOPATH/bin
	rm -f $(shell go env GOPATH)/bin/$(APP_NAME)

# ── Utilities ─────────────────────────────────────────────────────────────────
deps: ## Download and tidy dependencies
	go mod tidy

deps-update: ## Update all dependencies
	go get -u ./...
	go mod tidy

clean: ## Clean build artifacts
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)

info: ## Show build information
	@echo "Version:     $(VERSION)"
	@echo "Module:      $(MODULE)"
	@echo "Build Date:  $(BUILD_DATE)"
	@echo "Commit:      $(COMMIT)"
	@echo "Go Version:  $(shell go version)"
	@echo "Platform:    $(shell go env GOOS)/$(shell go env GOARCH)"

demo: build ## Run demo commands
	./$(APP_NAME) --version
	./$(APP_NAME) themes
	./$(APP_NAME) render test.md
