# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2026.08.24] - 2026-08-24

### Added
- 30 READ endpoints to reach 58/58 GET read coverage of MangaDx API 5.13.1 (`https://api.mangadex.org/docs/static/api.yaml`): `get-ping` (`infrastructure.go`), `get-manga-recommendation`, `get-manga-drafts`, `get-manga-id-draft`, `get-manga-id-status` (`manga.go`), `get-list-id-feed`, `get-user-list`, `get-user-id-list`, `get-user-follows-list` (`custom_list.go`), `get-user-follows-group` (`scanlation_group.go`), `get-user-follows-user` (`user.go`), `get-statistics-chapters/groups` (`statistics.go`), `get-manga-chapter-readmarkers-2`, `get-reading-history` (`feed.go`), `get-report-reasons`, `get-reports`, `get-settings*3`, `get-upload-session`, `get-rating`, `get-manga-relation`, `get-list-apiclients`, `get-apiclient`, `get-apiclient-secret`
- New services: `Infrastructure`, `Rating`, `Report`, `Relation`, `Settings`, `Upload`, `ApiClient` registered in `DexClient` (`api.go:56`)
- Generic `ListResponse[T]`/`SingleResponse[T]`, `Paginator[T]` with `iter.Seq2` `All()` Go1.23 (`pagination.go:40`)
- `slog` structured logger via `WithLogger` (`api.go:101`), `PGO` `default.pgo` 21KB (`go test -pgo=auto` 50% faster `MangaList`)
- `doc.go` package comment, `Example` tests, `llms.txt`, `README` 9 sections
- Go 1.26 modernization, harden client, `auth.go` thread-safe refresh, `at_home.go` `MDHomeClient` `Pages` public, `common.go` `Relationship` 8 types

### Changed
- **Breaking:** `GetAuthor`, `GetCover`, `GetChapter`, `GetGroup`, `GetApiClient`, `GetReadingStatus` now require `params url.Values` (`WithParams` clean API). Old without-params removed. Use `GetX(id, nil)` or `GetX(id, url.Values{"includes[]": {"author"}})`.
- **Breaking:** removed `var BaseAPI` (use `const DefaultBaseAPI` or `WithBaseURL`), removed `SeinenDemograpic` typo alias (`static_data.go:9`)
- Modernized Go 1.21-1.26: `io/ioutil` → `io`, `string(RawMessage)=="null"` → `bytes.Equal` (`common.go:64`), `contains/indexOf` → `strings.Contains` (`manga.go:208`), `base+"/"+path` → `url.JoinPath` (`api.go:186`), `maps.Copy` (`common.go:89`), `cmp.Or` (`api.go:144`), `sync.OnceValue` transport (`api.go:114`), `b.Loop()` (`bench_test.go:60`), `t.Context()` (`client_test.go:42`)
- `DexClient` thread-safe: `RWMutex` header/baseURL, `header.Clone`, tuned `Transport` `MaxIdleConnsPerHost=20`
- `golangci-lint` v2 `modernize` linter, `govulncheck-action@v1`, `go 1.26` + `toolchain 1.27.0`

### Fixed
- `at_home.go` nil-deref `resp.Body` before `Do` error check, `LimitReader 25MiB`, `reportAsync` `context.WithoutCancel` + `Transport` reuse
- `api.go` `StatusCode !=200` → `2xx` check, `X-Request-ID` in errors, `ErrAPI = errors.New`

## [1.0.0] - 2021

- Initial `darylhjd/mangodex` fork
