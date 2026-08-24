# Contributing

Thanks for contributing to `mangodex`! Setup <10 min.

## Prerequisites

- Go 1.26+ (`go version`), toolchain 1.27+ recommended (`GOTOOLCHAIN=auto`)
- `golangci-lint` v2.13+ (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

## Clone

```bash
git clone https://github.com/yukiteruamano/mangodex.git
cd mangodex
```

## Build & Test

```bash
go vet ./...
golangci-lint run ./...
go test -race -cover -pgo=auto ./...
# cover gate
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total # >=80%
# bench + PGO profile
go test -run=^$ -bench=BenchmarkMangaList -count=1 -cpuprofile /tmp/cpu.pprof
```

## PR Process

1. Branch from `main` (`feat/`, `fix/`, `docs/`).
2. Keep `go 1.26` in `go.mod` (toolchain may be 1.27+). Use `cmp.Or`, `maps.Copy`, `url.JoinPath`, `bytes.Equal`, `b.Loop()`, `t.Context()`, `iter.Seq2` where applicable.
3. Every exported func needs `// Foo ...` doc comment + `WithContext` variant already exists — keep pattern.
4. Run `gofmt -l` 0, `go vet`, `golangci-lint run` (modernize linter excludes `newexpr`).
5. Add `httptest` tests via `newTestClient(t)` (`client_test.go:16`) + `t.Context()`, target `80.6%` cover.
6. If touching `at_home.go` or `api.go`, test `X-Cache` and `forcePort443`.

## Commit

- `feat:`, `fix:`, `perf:`, `docs:`, `chore:` + detailed body + `benchstat` if perf.
- Breaking changes bump `v3` (`BaseAPI` removed, `6` `WithParams`).

## Release

- Tag `v3.x.x` matching `https://api.mangadex.org/docs/static/api.yaml` `v5.13.1` read coverage.
- `default.pgo` committed (generated from `BenchmarkMangaList`).

## Code of Conduct

Be respectful. Issues welcome.
