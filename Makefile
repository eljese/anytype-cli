.PHONY: all build install install-local uninstall uninstall-local clean clean-tantivy download-tantivy lint lint-fix install-linter

all: download-tantivy build

VERSION ?= $(shell git describe --tags 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_STATE ?= $(shell git diff --quiet 2>/dev/null && echo "clean" || echo "dirty")

LDFLAGS := -s -w \
           -X 'github.com/anyproto/anytype-cli/core.Version=$(VERSION)' \
           -X 'github.com/anyproto/anytype-cli/core.Commit=$(COMMIT)' \
           -X 'github.com/anyproto/anytype-cli/core.BuildTime=$(BUILD_TIME)' \
           -X 'github.com/anyproto/anytype-cli/core.GitState=$(GIT_STATE)'

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUTPUT ?= dist/anytype

TANTIVY_VERSION := v1.0.4
TANTIVY_LIB_PATH := dist/tantivy
CGO_LDFLAGS := -L$(TANTIVY_LIB_PATH)

GOLANGCI_LINT_VERSION := v2.2.1

##@ Build

build: download-tantivy ## Build the CLI binary
	@echo "Building Anytype CLI with embedded server..."
	@CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)
	@echo "Built successfully: $(OUTPUT)"

cross-compile: ## Build for all platforms in parallel
	@echo "Cross-compiling for all platforms..."
	@$(MAKE) -j \
		build-darwin-amd64 \
		build-darwin-arm64 \
		build-windows-amd64 \
		build-linux-amd64 \
		build-linux-arm64
	@echo "All platforms built successfully!"

build-darwin-amd64:
	@$(MAKE) clean-tantivy
	@GOOS=darwin GOARCH=amd64 OUTPUT=dist/anytype-darwin-amd64 $(MAKE) build

build-darwin-arm64:
	@$(MAKE) clean-tantivy
	@GOOS=darwin GOARCH=arm64 OUTPUT=dist/anytype-darwin-arm64 $(MAKE) build

build-windows-amd64:
	@$(MAKE) clean-tantivy
	@GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc OUTPUT=dist/anytype-windows-amd64.exe $(MAKE) build

build-linux-amd64:
	@$(MAKE) clean-tantivy
	@GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc OUTPUT=dist/anytype-linux-amd64 $(MAKE) build

build-linux-arm64:
	@$(MAKE) clean-tantivy
	@GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc OUTPUT=dist/anytype-linux-arm64 $(MAKE) build

download-tantivy: ## Download tantivy library for current platform
	@if [ ! -f "$(TANTIVY_LIB_PATH)/libtantivy_go.a" ]; then \
		echo "Downloading tantivy library $(TANTIVY_VERSION) for $(GOOS)/$(GOARCH)..."; \
		mkdir -p $(TANTIVY_LIB_PATH); \
		if [ "$(GOOS)" = "darwin" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/darwin-amd64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			elif [ "$(GOARCH)" = "arm64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/darwin-arm64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for macOS"; \
				exit 1; \
			fi; \
		elif [ "$(GOOS)" = "linux" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/linux-amd64-musl.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			elif [ "$(GOARCH)" = "arm64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/linux-arm64-musl.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for Linux"; \
				exit 1; \
			fi; \
		elif [ "$(GOOS)" = "windows" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/windows-amd64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for Windows"; \
				exit 1; \
			fi; \
		else \
			echo "Unsupported OS: $(GOOS)"; \
			exit 1; \
		fi; \
		echo "Tantivy library downloaded successfully"; \
	else \
		echo "Tantivy library already exists"; \
	fi

##@ Installation

install: build ## Install to /usr/local/bin (may require sudo)
	@echo "Installing Anytype CLI..."
	@cp dist/anytype /usr/local/bin/anytype 2>/dev/null || sudo cp dist/anytype /usr/local/bin/anytype
	@echo "Installed to /usr/local/bin/"
	@echo ""
	@echo "Usage:"
	@echo "  anytype serve              # Run server in foreground"
	@echo "  anytype service install    # Install as system service"

install-local: build ## Install to ~/.local/bin (user installation)
	@mkdir -p $$HOME/.local/bin
	@cp dist/anytype $$HOME/.local/bin/anytype
	@echo "Installed to $$HOME/.local/bin/"
	@echo "Make sure $$HOME/.local/bin is in your PATH"
	@echo ""
	@echo "Usage:"
	@echo "  anytype serve              # Run server in foreground"
	@echo "  anytype service install    # Install as system service"

uninstall: ## Uninstall from /usr/local/bin
	@echo "Uninstalling Anytype CLI..."
	@rm -f /usr/local/bin/anytype 2>/dev/null || sudo rm -f /usr/local/bin/anytype
	@echo "Uninstalled from /usr/local/bin/"

uninstall-local: ## Uninstall from ~/.local/bin
	@echo "Uninstalling Anytype CLI from local..."
	@rm -f $$HOME/.local/bin/anytype
	@echo "Uninstalled from $$HOME/.local/bin/"

##@ Development

lint: ## Run linters
	@golangci-lint run ./...

lint-fix: ## Run linters with auto-fix
	@golangci-lint run --fix ./...

install-linter: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "golangci-lint installed successfully"

##@ Cleanup

clean: clean-tantivy ## Clean all build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf dist/
	@echo "Build artifacts cleaned"

clean-tantivy: ## Clean tantivy libraries
	@echo "Cleaning tantivy libraries..."
	@rm -rf $(TANTIVY_LIB_PATH)
	@echo "Tantivy libraries cleaned"

##@ Other

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)