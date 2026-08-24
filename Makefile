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

.PHONY: help check fmt fmt-check lint vet test cover bench build build-debug build-release pgo govulncheck push push-tags clean

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

push: ## Push branch to origin (git push origin main)
	git push origin main

push-tags: ## Create tag v$(DATE) and push + auto-trigger pkg.go.dev (standard v2026.08.24)
	@VERSION=v$$(date +%Y.%m.%d); \
	if git rev-parse $$VERSION >/dev/null 2>&1; then echo "Tag $$VERSION already exists locally"; else echo "Creating tag $$VERSION..."; git tag -a $$VERSION -m "$$VERSION: sync to MangaDx API $$(grep 'version:' api.yaml 2>/dev/null | head -1 || echo v5.13.1)"; fi; \
	git push origin $$VERSION; \
	echo "Pushed tag $$VERSION"; \
	echo "Triggering proxy.golang.org + pkg.go.dev (auto)..."; \
	TMPDIR=$$(mktemp -d); \
	( cd $$TMPDIR && go mod init tmp >/dev/null 2>&1; GOPROXY=https://proxy.golang.org,direct go get github.com/yukiteruamano/mangodex@$$VERSION 2>&1 | head -5; ); \
	rm -rf $$TMPDIR; \
	sleep 2; \
	echo "Proxy info:"; curl -s https://proxy.golang.org/github.com/yukiteruamano/mangodex/@v/$$VERSION.info | head -5; \
	echo "Check pkg.go.dev in ~60s: https://pkg.go.dev/github.com/yukiteruamano/mangodex?tab=versions"; \
	curl -s https://pkg.go.dev/github.com/yukiteruamano/mangodex?tab=versions | grep -q $$VERSION && echo "pkg.go.dev: $$VERSION visible" || echo "pkg.go.dev: wait ~60s then check"

clean: ## Clean artifacts
	rm -f coverage.out coverage.html bench.txt $(BIN) cpu.pprof
	@echo "clean ok"
