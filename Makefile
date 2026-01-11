.PHONY: build build-all test test-coverage lint fmt vet clean install setup-hooks

# Build variables
BINARY_NAME=updoc
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build for current platform
build:
	@echo "==> Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/updoc

# Build for all platforms
build-all:
	@echo "==> Building for all platforms..."
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/updoc
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/updoc
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/updoc
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/updoc
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/updoc

# Run tests
test:
	@echo "==> Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "==> Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "==> Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Format code
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "==> Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "==> Cleaning..."
	rm -rf bin/ dist/ coverage.out coverage.html

# Install binary locally
install: build
	@echo "==> Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) /usr/local/bin/

# Setup git hooks
setup-hooks:
	@echo "==> Setting up git hooks..."
	git config core.hooksPath .githooks
	chmod +x .githooks/*
	@echo "Git hooks configured to use .githooks/"

# Run all checks (same as pre-commit)
check: fmt vet lint test
	@echo "==> All checks passed!"

# Development setup
dev-setup: setup-hooks
	@echo "==> Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment ready!"
