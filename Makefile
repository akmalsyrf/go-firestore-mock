GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all
all: deps fmt vet test

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

.PHONY: test-coverage-html
test-coverage-html: test-coverage
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

.PHONY: test-race
test-race:
	$(GOTEST) -race -v ./...

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: vet
vet:
	$(GOVET) ./...

.PHONY: lint
lint:
	$(GOLINT) run

.PHONY: lint-fix
lint-fix:
	$(GOLINT) run --fix

.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update:
	$(GOMOD) get -u
	$(GOMOD) tidy

.PHONY: generate
generate:
	$(GOCMD) generate ./...

.PHONY: deps-outdated
deps-outdated:
	$(GOCMD) list -u -m all

.PHONY: ci
ci: deps fmt vet test-coverage

.PHONY: dev
dev: deps fmt vet test

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all                - deps, fmt, vet, test"
	@echo "  clean              - Remove generated coverage files"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  test-race          - Run tests with race detection"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  lint               - Run linter (requires golangci-lint)"
	@echo "  lint-fix           - Run linter with --fix"
	@echo "  deps               - Download dependencies and run go mod tidy"
	@echo "  deps-update        - Update dependencies"
	@echo "  deps-outdated      - List outdated modules"
	@echo "  generate           - Regenerate mocks via go generate"
	@echo "  ci                 - deps, fmt, vet, test-coverage"
	@echo "  dev                - deps, fmt, vet, test"
