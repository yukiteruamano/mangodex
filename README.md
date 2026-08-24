# mangodex

[![Go Version](https://img.shields.io/github/go-mod/go-version/yukiteruamano/mangodex)](https://go.dev/) [![License](https://img.shields.io/github/license/yukiteruamano/mangodex)](./LICENSE) [![Build Status](https://img.shields.io/github/actions/workflow/status/yukiteruamano/mangodex/ci.yml?branch=main)](https://github.com/yukiteruamano/mangodex/actions) [![Coverage](https://img.shields.io/badge/coverage-80.6%25-brightgreen)](https://github.com/yukiteruamano/mangodex) [![Go Report Card](https://goreportcard.com/badge/github.com/yukiteruamano/mangodex)](https://goreportcard.com/report/github.com/yukiteruamano/mangodex) [![Go Reference](https://pkg.go.dev/badge/github.com/yukiteruamano/mangodex.svg)](https://pkg.go.dev/github.com/yukiteruamano/mangodex)

Go client for [MangaDx API v5.13.1](https://api.mangadex.org/docs/static/api.yaml) — typed, thread-safe, and PGO-optimized. Covers all 58 READ endpoints with generic pagination, context propagation, and scrapper-tuned transport.

> Fork of [`darylhjd/mangodex`](https://github.com/darylhjd/mangodex) — original abandoned — rewritten for Go 1.26+.

## Demo

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/url"

    "github.com/yukiteruamano/mangodex"
)

func main() {
    c := mangodex.NewDexClient()
    list, err := c.Manga.GetMangaListContext(context.Background(), url.Values{
        "title": {"One Piece"},
        "limit": {"5"},
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, m := range list.Data {
        fmt.Printf("%s — %s\n", m.ID, m.GetTitle("en"))
    }

    // Paginate with iterator (Go 1.23)
    p := mangodex.NewPaginator(func(ctx context.Context, limit, offset int) (*mangodex.ListResponse[mangodex.Manga], error) {
        return c.Manga.GetMangaListContext(ctx, url.Values{
            "limit": {fmt.Sprint(limit)}, "offset": {fmt.Sprint(offset)},
        })
    })
    for data, err := range p.All(context.Background()) {
        if err != nil { break }
        fmt.Println("page:", len(data))
    }
}
```

## Getting Started

**Prerequisites:** Go 1.26+, toolchain 1.27+ recommended (`GOTOOLCHAIN=auto`).

**Install:**

```bash
go get github.com/yukiteruamano/mangodex@v3
```

**Minimal example — fetch manga by ID:**

```go
c := mangodex.NewDexClient()
manga, err := c.Manga.GetMangaContext(ctx, "f9c33607-9180-4ba6-b85c-e4b5fa63ff33", url.Values{"includes[]": {"author"}})
```

**Authenticated:**

```go
c := mangodex.NewDexClient()
if err := c.Auth.LoginContext(ctx, "user", "pass"); err != nil { log.Fatal(err) }
me, _ := c.User.GetLoggedUserContext(ctx)
```

**Images (At-Home):**

```go
ah, _ := c.AtHome.NewMDHomeClientContext(ctx, chapterID, "data", false)
data, _ := ah.GetChapterPageWithContext(ctx, ah.Pages[0])
```

## Features / Specification

- **58 READ endpoints** — Manga, Chapter, Author, Cover, ScanlationGroup, User, CustomList, Feed, Statistics, AtHome, Infrastructure (`/ping` `text/plain`), Rating, Report, Relation, Settings, Upload, ApiClient, Tag
- **Generic responses:** `ListResponse[T]` / `SingleResponse[T]` with `GetResult()`
- **Every method has `WithContext` variant** (`GetMangaListContext(ctx, params)`), `context.Background` wrapper for convenience
- **Params via `url.Values`** — supports `includes[]`, `manga[]`, `group[]`, `contentRating[]`, `order[title]=asc` etc. 6 former `GetX(id)` now `GetX(id, params)` for `includes[]` correctness
- **Follow checks:** `CheckIfMangaFollowed`, `CheckFollowedGroup/User/CustomList` map 404 → `false, nil`
- **At-Home:** `MDHomeClient` with `25MiB` limit, `X-Cache` reporting via `api.mangadex.network/report`, `forcePort443`
- **Thread-safe client:** `RWMutex` header/baseURL, `header.Clone()` per request, tuned `Transport` `MaxIdleConnsPerHost=20` via `sync.OnceValue`
- **Observability:** `WithLogger(*slog.Logger)` for non-2xx, `X-Request-ID` in errors, `PGO` `default.pgo` (50% faster `MangaList` unmarshal)
- **Breaking v3:** removed `var BaseAPI` (use `DefaultBaseAPI` const or `WithBaseURL`), removed `SeinenDemograpic` typo alias, 6 `GetX` now require `params`

OpenAPI: `https://api.mangadex.org/docs/static/api.yaml` (v5.13.1, 113 ops, 5 deprecated -> 108, 58 GET)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md). Setup <10 min: clone, `go vet`, `golangci-lint run`, `go test -race -cover`.

## Contributors

Thanks to `darylhjd` (original) and all contributors.

[![Contributors](https://img.shields.io/github/contributors/yukiteruamano/mangodex)](https://github.com/yukiteruamano/mangodex/graphs/contributors)

## License

[MIT](./LICENSE) — Copyright (c) 2021 darylhjd, 2026 yukiteruamano
