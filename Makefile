# lazyvibe - btop-style CLI dashboard for Claude Code
# Usage: make <target>

.PHONY: all build clean test run install install-local help deps fmt lint

BINARY_NAME = lazyvibe
GO = go
BUILD_DIR = build

# Get version from git tag if available
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"

# Terminal sizes for captures
SIZE_SMALL := 80x24
SIZE_MEDIUM := 120x40
SIZE_LARGE := 160x50

# Capture output directory
CAPTURE_DIR := .dev-captures

all: build

#───────────────────────────────────────────────────────────────
# Help
#───────────────────────────────────────────────────────────────

help: ## Show this help message
	@echo "lazyvibe - btop-style CLI dashboard for Claude Code"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

#───────────────────────────────────────────────────────────────
# Build & Run
#───────────────────────────────────────────────────────────────

build: ## Build the binary
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/lazyvibe

run: build ## Build and run the TUI
	./$(BINARY_NAME)

install: ## Install to GOPATH/bin
	$(GO) install $(LDFLAGS) ./cmd/lazyvibe

install-local: build ## Install to /usr/local/bin (requires sudo)
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"
	@echo "  Run 'lazyvibe' from any terminal"

#───────────────────────────────────────────────────────────────
# Dev Workflow - Capture Commands
#───────────────────────────────────────────────────────────────

dump: build ## Dump raw JSON data to stdout
	./$(BINARY_NAME) --dump

capture: build ## Capture ASCII at medium size (120x40)
	./$(BINARY_NAME) --capture $(SIZE_MEDIUM)

capture-small: build ## Capture ASCII at small size (80x24)
	./$(BINARY_NAME) --capture $(SIZE_SMALL)

capture-large: build ## Capture ASCII at large size (160x50)
	./$(BINARY_NAME) --capture $(SIZE_LARGE)

capture-all: build ## Capture all sizes + JSON to timestamped directory
	@mkdir -p $(CAPTURE_DIR)
	@TIMESTAMP=$$(date +%Y%m%d-%H%M%S) && \
	DIR=$(CAPTURE_DIR)/$$TIMESTAMP && \
	mkdir -p $$DIR && \
	echo "Capturing to $$DIR..." && \
	./$(BINARY_NAME) --capture $(SIZE_SMALL) > $$DIR/capture-$(SIZE_SMALL).txt && \
	./$(BINARY_NAME) --capture $(SIZE_MEDIUM) > $$DIR/capture-$(SIZE_MEDIUM).txt && \
	./$(BINARY_NAME) --capture $(SIZE_LARGE) > $$DIR/capture-$(SIZE_LARGE).txt && \
	./$(BINARY_NAME) --dump > $$DIR/data.json && \
	echo "✓ Captures saved to $$DIR" && \
	ls -la $$DIR

dev: build ## Run dev captures (small + medium) - for quick iteration
	@echo ""
	@echo "─── 80x24 ───"
	@./$(BINARY_NAME) --capture $(SIZE_SMALL)
	@echo ""
	@echo "─── 120x40 ───"
	@./$(BINARY_NAME) --capture $(SIZE_MEDIUM)

#───────────────────────────────────────────────────────────────
# Cross-Platform Builds
#───────────────────────────────────────────────────────────────

build-all: clean ## Build for all platforms
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/lazyvibe
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/lazyvibe
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/lazyvibe
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/lazyvibe
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/lazyvibe
	@echo "✓ Built for all platforms in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

#───────────────────────────────────────────────────────────────
# Code Quality
#───────────────────────────────────────────────────────────────

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Lint code
	$(GO) vet ./...

test: ## Run tests
	$(GO) test -v ./...

test-quick: build ## Quick smoke test - verify app starts
	@echo "Testing --dump..."
	@./$(BINARY_NAME) --dump > /dev/null && echo "✓ --dump works"
	@echo "Testing --capture..."
	@./$(BINARY_NAME) --capture 80x24 > /dev/null && echo "✓ --capture works"
	@echo ""
	@echo "✓ All smoke tests passed"

#───────────────────────────────────────────────────────────────
# Utilities
#───────────────────────────────────────────────────────────────

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -rf $(CAPTURE_DIR)
	@echo "✓ Cleaned"

clean-captures: ## Remove only capture directory
	rm -rf $(CAPTURE_DIR)
	@echo "✓ Captures cleaned"

tree: ## Show project structure
	@echo "lazyvibe project structure:"
	@find cmd internal -type f -name "*.go" | sort | sed 's|^|  |'

watch: build ## Watch for changes and re-capture (requires entr)
	@command -v entr >/dev/null 2>&1 || { echo "Install entr: brew install entr"; exit 1; }
	@echo "Watching for changes... (Ctrl+C to stop)"
	@find cmd internal -name "*.go" | entr -c make capture

#───────────────────────────────────────────────────────────────
# Comparison (for dev workflow)
#───────────────────────────────────────────────────────────────

compare: ## Compare two capture directories (BEFORE=dir AFTER=dir)
	@if [ -z "$(BEFORE)" ] || [ -z "$(AFTER)" ]; then \
		echo "Usage: make compare BEFORE=<dir> AFTER=<dir>"; \
		echo ""; \
		echo "Available captures:"; \
		ls -1 $(CAPTURE_DIR) 2>/dev/null || echo "  (none - run 'make capture-all' first)"; \
		exit 1; \
	fi
	@echo "Comparing $(SIZE_MEDIUM)..."
	@diff --side-by-side --suppress-common-lines \
		$(BEFORE)/capture-$(SIZE_MEDIUM).txt \
		$(AFTER)/capture-$(SIZE_MEDIUM).txt || true

list-captures: ## List available capture directories
	@echo "Available captures in $(CAPTURE_DIR):"
	@ls -lt $(CAPTURE_DIR) 2>/dev/null || echo "  (none - run 'make capture-all' first)"
