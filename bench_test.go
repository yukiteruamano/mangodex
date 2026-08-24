package mangodex

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var benchMangaListJSON10 []byte
var benchMangaSingleJSON []byte

func init() {
	// Build realistic manga list with 10 items, each with tags and relationships
	m := Manga{
		ID:   "manga-id",
		Type: "manga",
		Attributes: MangaAttributes{
			Title:            LocalisedStrings{Values: map[string]string{"en": "Test Manga Title With Some Length"}},
			AltTitles:        LocalisedStrings{Values: map[string]string{"ja": "テスト"}},
			Description:      LocalisedStrings{Values: map[string]string{"en": "A very long description that simulates real manga description with many words and details about the story and characters and plot"}},
			OriginalLanguage: "ja",
			State:            "published",
			Version:          1,
			CreatedAt:        "2020-01-01T00:00:00",
			UpdatedAt:        "2020-01-01T00:00:00",
			Year:             intPtr(2020), //nolint:modernize
			Tags: []Tag{
				{ID: "tag1", Type: "tag", Attributes: TagAttributes{Name: LocalisedStrings{Values: map[string]string{"en": "Action"}}, Group: "genre", Version: 1}},
				{ID: "tag2", Type: "tag", Attributes: TagAttributes{Name: LocalisedStrings{Values: map[string]string{"en": "Adventure"}}, Group: "genre", Version: 1}},
			},
		},
		Relationships: []Relationship{
			{ID: "author1", Type: AuthorRel, Attributes: &AuthorAttributes{Name: "Author Name"}},
		},
	}
	list := ListResponse[Manga]{Result: "ok", Response: "collection", Data: make([]Manga, 10), Limit: 10, Offset: 0, Total: 100}
	for i := range list.Data {
		list.Data[i] = m
		list.Data[i].ID = "manga-" + string(rune('0'+i))
	}
	b, _ := json.Marshal(list)
	benchMangaListJSON10 = b

	single := SingleResponse[Manga]{Result: "ok", Response: "entity", Data: m}
	bb, _ := json.Marshal(single)
	benchMangaSingleJSON = bb
}

func BenchmarkGetMangaList_10(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(benchMangaListJSON10)
	}))
	defer srv.Close()
	client := NewDexClient(WithBaseURL(srv.URL))
	b.ReportAllocs()
	for b.Loop() {
		_, err := client.Manga.GetMangaListContext(b.Context(), url.Values{"limit": {"10"}})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetLocalString_Hit(b *testing.B) {
	ls := LocalisedStrings{Values: map[string]string{"en": "Hello World Hello World", "ja": "こんにちは", "fr": "Bonjour", "de": "Hallo"}}
	b.ReportAllocs()
	for b.Loop() {
		_ = ls.GetLocalString("en")
	}
}

func BenchmarkGetLocalString_Miss(b *testing.B) {
	ls := LocalisedStrings{Values: map[string]string{"en": "Hello", "ja": "こんにちは"}}
	b.ReportAllocs()
	for b.Loop() {
		_ = ls.GetLocalString("fr")
	}
}

func BenchmarkLocalisedStrings_Unmarshal_Map(b *testing.B) {
	data := []byte(`{"en":"Hello World","ja":"こんにちは","fr":"Bonjour"}`)
	b.ReportAllocs()
	for b.Loop() {
		var ls LocalisedStrings
		_ = json.Unmarshal(data, &ls)
	}
}

func BenchmarkLocalisedStrings_Unmarshal_Array(b *testing.B) {
	data := []byte(`[{"en":"Hello"},{"ja":"こんにちは"}]`)
	b.ReportAllocs()
	for b.Loop() {
		var ls LocalisedStrings
		_ = json.Unmarshal(data, &ls)
	}
}

func BenchmarkRelationship_Unmarshal_Manga(b *testing.B) {
	data := []byte(`{"id":"manga-id","type":"manga","attributes":{"title":{"en":"Title"},"altTitles":{},"description":{},"isLocked":false,"links":{},"originalLanguage":"ja","state":"published","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","tags":[]}}`)
	b.ReportAllocs()
	for b.Loop() {
		var rel Relationship
		_ = json.Unmarshal(data, &rel)
	}
}

func BenchmarkMangaList_Unmarshal_10(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var list ListResponse[Manga]
		_ = json.Unmarshal(benchMangaListJSON10, &list)
	}
}
