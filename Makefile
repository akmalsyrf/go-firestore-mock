GOCMD=go
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

.PHONY: all
all: deps fmt vet test

.PHONY: clean
clean:
	$(GOCMD) clean
	rm -f coverage.out coverage.html

.PHONY: test
test:
	$(GOTEST) -count=1 ./...

.PHONY: test-race
test-race:
	$(GOTEST) -race -count=1 ./...

.PHONY: test-integration
test-integration:
	$(GOTEST) -tags=integration -count=1 ./fstest/...

.PHONY: accuracy
accuracy:
	./scripts/check-accuracy.sh

.PHONY: fmt
fmt:
	$(GOFMT) ./...

.PHONY: vet
vet:
	$(GOVET) ./...

.PHONY: lint
lint:
	$(GOLINT) run

.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: generate
generate:
	$(GOCMD) generate ./...

.PHONY: generate-check
generate-check: generate
	git diff --exit-code -- mocks/

.PHONY: gate
gate: deps fmt vet lint test-race generate-check accuracy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test / test-race     - Unit + apicheck"
	@echo "  test-integration     - All fstest (needs FIRESTORE_EMULATOR_HOST)"
	@echo "  accuracy             - PR accuracy gate (apicheck + required parity; skip=fail)"
	@echo "  generate-check       - Fail if mocks stale"
	@echo "  gate                 - Full local PR gate (needs emulator for accuracy)"
