GOCMD=go
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint
MOCKGEN_VERSION=v0.6.0
GOLANGCI_LINT_VERSION=v2.13.2
EMULATOR_NAME=fsmock-firestore-emulator
EMULATOR_PORT=8080

.PHONY: all
all: deps fmt vet test

.PHONY: clean
clean:
	$(GOCMD) clean
	rm -f coverage.out cover.out coverage.html

.PHONY: bootstrap
bootstrap:
	$(GOCMD) install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
	$(GOCMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: emulator-up
emulator-up:
	@if docker ps --format '{{.Names}}' | grep -q '^$(EMULATOR_NAME)$$'; then \
		echo "emulator already running"; \
	else \
		docker rm -f $(EMULATOR_NAME) >/dev/null 2>&1 || true; \
		docker run -d --name $(EMULATOR_NAME) -p $(EMULATOR_PORT):8080 \
			gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators \
			gcloud beta emulators firestore start --host-port=0.0.0.0:8080; \
		for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15; do \
			if nc -z 127.0.0.1 $(EMULATOR_PORT) 2>/dev/null; then echo "emulator up"; exit 0; fi; \
			sleep 2; \
		done; \
		docker logs $(EMULATOR_NAME) || true; \
		exit 1; \
	fi

.PHONY: emulator-down
emulator-down:
	docker rm -f $(EMULATOR_NAME) >/dev/null 2>&1 || true

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
	@if [[ -z "$${FIRESTORE_EMULATOR_HOST:-}" ]]; then \
		$(MAKE) emulator-up; \
		export FIRESTORE_EMULATOR_HOST=127.0.0.1:$(EMULATOR_PORT); \
	fi; \
	FIRESTORE_EMULATOR_HOST=$${FIRESTORE_EMULATOR_HOST:-127.0.0.1:$(EMULATOR_PORT)} ./scripts/check-accuracy.sh

.PHONY: version-check
version-check:
	./scripts/check-version.sh

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
gate: deps fmt vet lint test-race generate-check version-check accuracy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  bootstrap            - Install mockgen + golangci-lint (CI versions)"
	@echo "  emulator-up/down     - Docker Firestore emulator"
	@echo "  test / test-race     - Unit + apicheck"
	@echo "  test-integration     - All fstest (needs FIRESTORE_EMULATOR_HOST)"
	@echo "  accuracy             - PR accuracy gate (starts emulator if needed)"
	@echo "  version-check        - Pairing: version.go ↔ go.mod ↔ COMPATIBILITY.md"
	@echo "  generate-check       - Fail if mocks stale"
	@echo "  gate                 - Full local PR gate"
