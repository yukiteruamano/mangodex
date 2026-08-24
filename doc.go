// Package mangodex is a Go client for the MangaDex API v5.13.1.
//
// It provides typed access to all 58 read (GET) endpoints of
// https://api.mangadex.org/docs/static/api.yaml with thread-safe client,
// context propagation, and PGO-optimized JSON decoding.
//
// # Features
//
//   - 58 READ endpoints: Manga, Chapter, Author, Cover, ScanlationGroup,
//     User, CustomList, Feed, Statistics, AtHome, Infrastructure, Rating,
//     Report, Relation, Settings, Upload, ApiClient, Tag
//   - Generic paginated responses ListResponse[T] / SingleResponse[T]
//   - WithContext variants for every method + slog/PGO support
//   - Tuned transport for scrapper massive (MaxIdleConnsPerHost 20)
//
// # Usage
//
//	c := mangodex.NewDexClient(mangodex.WithBaseURL("https://api.mangadex.org"))
//	list, err := c.Manga.GetMangaListContext(ctx, url.Values{"title": {"One Piece"}})
//	if err != nil { log.Fatal(err) }
//	fmt.Println(list.Data[0].GetTitle("en"))
//
// For authenticated calls, set Authorization header or use AuthService:
//
//	c := mangodex.NewDexClient()
//	if err := c.Auth.LoginContext(ctx, "user", "pass"); err != nil { log.Fatal(err) }
//
// Pagination via Paginator[T]:
//
//	p := mangodex.NewPaginator(func(ctx context.Context, limit, offset int) (*mangodex.ListResponse[mangodex.Manga], error) {
//	    return c.Manga.GetMangaListContext(ctx, url.Values{"limit": {fmt.Sprint(limit)}, "offset": {fmt.Sprint(offset)}})
//	})
//	for data, err := range p.All(ctx) {
//	    if err != nil { break }
//	    fmt.Println(len(data))
//	}
//
// API docs: https://api.mangadex.org/docs/static/api.yaml (v5.13.1)
// Pkg docs: https://pkg.go.dev/github.com/yukiteruamano/mangodex
package mangodex
