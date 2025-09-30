# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Project parameters
BINARY_NAME=go-firestore-mock
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=.
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Default target
.PHONY: all
all: clean deps fmt vet test build

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Generate HTML coverage report
.PHONY: test-coverage-html
test-coverage-html: test-coverage
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# Run tests with race detection
.PHONY: test-race
test-race:
	$(GOTEST) -race -v ./...

# Run tests with benchmarks
.PHONY: test-bench
test-bench:
	$(GOTEST) -bench=. -v ./...

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

# Run linter
.PHONY: lint
lint:
	$(GOLINT) run

# Run linter with fix
.PHONY: lint-fix
lint-fix:
	$(GOLINT) run --fix

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GOMOD) get -u
	$(GOMOD) tidy

# Install dependencies
.PHONY: deps-install
deps-install:
	$(GOGET) -d -v ./...

# Generate mocks
.PHONY: generate
generate:
	$(GOCMD) generate ./...

# Run the application
.PHONY: run
run:
	$(GOCMD) run $(MAIN_PATH)

# Install the binary
.PHONY: install
install:
	$(GOCMD) install $(MAIN_PATH)

# Check for security vulnerabilities
.PHONY: security
security:
	$(GOCMD) list -json -deps ./... | nancy sleuth

# Check for outdated dependencies
.PHONY: deps-outdated
deps-outdated:
	$(GOCMD) list -u -m all

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Clean, deps, fmt, vet, test, and build"
	@echo "  build            - Build the application"
	@echo "  build-linux      - Build for Linux"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  test-race        - Run tests with race detection"
	@echo "  test-bench       - Run tests with benchmarks"
	@echo "  fmt              - Format code"
	@echo "  vet              - Run go vet"
	@echo "  lint             - Run linter"
	@echo "  lint-fix         - Run linter with fix"
	@echo "  deps             - Download dependencies"
	@echo "  deps-update      - Update dependencies"
	@echo "  deps-install     - Install dependencies"
	@echo "  generate         - Generate mocks"
	@echo "  run              - Run the application"
	@echo "  install          - Install the binary"
	@echo "  security         - Check for security vulnerabilities"
	@echo "  deps-outdated    - Check for outdated dependencies"
	@echo "  help             - Show this help message"

# CI/CD targets
.PHONY: ci
ci: deps fmt vet lint test-coverage

# Development targets
.PHONY: dev
dev: deps fmt vet test

# Release targets
.PHONY: release
release: clean deps fmt vet lint test-coverage build build-linux
