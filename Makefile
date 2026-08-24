GO := go
PKG := ./...
BIN := mangodex.test
LOCAL := github.com/yukiteruamano/mangodex
GOFLAGS ?=
DEBUG ?= 0
CGO_ENABLED ?= 1

# Release: stripped, PIE, netgo, PGO; Debug: symbols, no opt, race, netgo, PGO
LDFLAGS_RELEASE := -trimpath -ldflags="-s -w" -tags netgo -buildmode=pie -pgo=auto
LDFLAGS_DEBUG := -gcflags="all=-N -l" -race -tags debug,netgo -pgo=auto

.PHONY: help check fmt fmt-check lint vet test cover bench build build-debug build-release pgo govulncheck clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

check: vet lint test govulncheck ## Gate: vet + lint + test (80%) + govulncheck
	@echo "check ok"

fmt: ## Format with gofmt and goimports (local $(LOCAL))
	gofmt -w .
	goimports -w -local $(LOCAL) .

fmt-check: ## Check formatting (gofmt + goimports local)
	@test -z "$$(gofmt -l .)" || (echo "gofmt -l:" && gofmt -l . && exit 1)
	@test -z "$$(goimports -l -local $(LOCAL) .)" || (echo "goimports -l:" && goimports -l -local $(LOCAL) . && exit 1)

lint: ## golangci-lint v2 modernize
	golangci-lint run ./...

vet: ## go vet
	$(GO) vet $(PKG)

test: ## go test -race -cover -pgo=auto (CGO_ENABLED=1 for race)
	CGO_ENABLED=1 $(GO) test -race -cover -covermode=atomic -pgo=auto $(PKG)

cover: ## coverage html
	CGO_ENABLED=1 $(GO) test -coverprofile=coverage.out -covermode=atomic -pgo=auto $(PKG)
	$(GO) tool cover -func=coverage.out | grep total
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "coverage.html generated"

bench: ## benchmem + PGO
	$(GO) test -bench=. -benchmem -count=6 -pgo=auto $(PKG) | tee bench.txt
	@echo "bench.txt generated (use benchstat bench.txt)"

build: ## Default: release if DEBUG=0 else debug
ifeq ($(DEBUG),1)
	@$(MAKE) build-debug
else
	@$(MAKE) build-release
endif

build-debug: ## Debug build: symbols, no opt, race, netgo, PGO (CGO_ENABLED=1)
	CGO_ENABLED=1 $(GO) test -c -o $(BIN) $(LDFLAGS_DEBUG) $(PKG)
	@echo "debug build: $(BIN) (CGO_ENABLED=1, -N -l, race, netgo, PGO)"

build-release: ## Release build: stripped, PIE, netgo, PGO (CGO_ENABLED=0)
	CGO_ENABLED=0 $(GO) build $(LDFLAGS_RELEASE) $(PKG)
	@echo "release build: stripped PIE netgo PGO (CGO_ENABLED=0)"

pgo: ## Generate default.pgo from bench
	$(GO) test -run=^$$ -bench=BenchmarkMangaList_Unmarshal_10 -count=3 -cpuprofile /tmp/cpu.pprof $(PKG)
	cp /tmp/cpu.pprof default.pgo
	@echo "default.pgo generated (21KB)"

govulncheck: ## govulncheck via go tool (Go1.24 tool directive)
	$(GO) tool govulncheck ./... || go run golang.org/x/vuln/cmd/govulncheck@latest ./...

clean: ## Clean artifacts
	rm -f coverage.out coverage.html bench.txt $(BIN) cpu.pprof
	@echo "clean ok"
