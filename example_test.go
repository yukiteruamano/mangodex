package mangodex_test

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/yukiteruamano/mangodex"
)

// Example_mangaSearch shows searching manga.
func Example_mangaSearch() {
	c := mangodex.NewDexClient()
	list, err := c.Manga.GetMangaListContext(context.Background(), url.Values{
		"title": {"One Piece"},
		"limit": {"5"},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range list.Data {
		fmt.Println(m.GetTitle("en"))
	}
}

// Example_authorFetch shows fetching an author with includes.
func Example_authorFetch() {
	c := mangodex.NewDexClient()
	author, err := c.Author.GetAuthorContext(context.Background(), "f9c33607-9180-4ba6-b85c-e4b5fa63ff33", url.Values{
		"includes[]": {"manga"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(author.Data.Attributes.Name)
}

// Example_paginate shows paginating manga with Paginator and iter.Seq2.
func Example_paginate() {
	c := mangodex.NewDexClient()
	p := mangodex.NewPaginator(func(ctx context.Context, limit, offset int) (*mangodex.ListResponse[mangodex.Manga], error) {
		return c.Manga.GetMangaListContext(ctx, url.Values{
			"limit":  {fmt.Sprint(limit)},
			"offset": {fmt.Sprint(offset)},
		})
	})
	ctx := context.Background()
	for data, err := range p.All(ctx) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("page size:", len(data))
		if !p.HasMore() {
			break
		}
	}
}

// Example_atHome shows fetching a chapter image.
func Example_atHome() {
	c := mangodex.NewDexClient()
	// chapterID from feed
	chapterID := "b7cb4d72-04f1-4fc1-8402-6f63ed91cc09"
	ah, err := c.AtHome.NewMDHomeClientContext(context.Background(), chapterID, "data", false)
	if err != nil {
		log.Fatal(err)
	}
	if len(ah.Pages) > 0 {
		data, err := ah.GetChapterPageWithContext(context.Background(), ah.Pages[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("bytes:", len(data))
	}
}

// Example_infrastructurePing shows ping.
func Example_infrastructurePing() {
	c := mangodex.NewDexClient()
	pong, err := c.Infrastructure.GetPingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)
}
