package mangodex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockMangaListJSON() []byte {
	m := Manga{
		ID:   "manga-id",
		Type: "manga",
		Attributes: MangaAttributes{
			Title:            LocalisedStrings{Values: map[string]string{"en": "Test Manga"}},
			Description:      LocalisedStrings{Values: map[string]string{"en": "Desc"}},
			OriginalLanguage: "ja",
			State:            "published",
			Version:          1,
			CreatedAt:        "2020-01-01T00:00:00",
			UpdatedAt:        "2020-01-01T00:00:00",
			Year:             intPtr(2020), //nolint:modernize
		},
	}
	list := ListResponse[Manga]{Result: "ok", Response: "collection", Data: []Manga{m}, Limit: 10, Offset: 0, Total: 1}
	b, _ := json.Marshal(list)
	return b
}

func intPtr(i int) *int { return &i } //nolint:modernize

func TestManga_GetMangaList_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/manga", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(mockMangaListJSON())
	})
	defer srv.Close()
	list, err := c.Manga.GetMangaListContext(t.Context(), url.Values{"limit": {"10"}})
	require.NoError(t, err)
	require.Len(t, list.Data, 1)
	assert.Equal(t, "Test Manga", list.Data[0].GetTitle("en"))
}

func TestManga_GetMangaList_Error(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Bad", Detail: "bad"}}})
	})
	defer srv.Close()
	_, err := c.Manga.GetMangaListContext(t.Context(), url.Values{})
	require.Error(t, err)
}

func TestManga_GetManga_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/manga/manga-id", r.URL.Path)
		resp := SingleResponse[Manga]{Result: "ok", Data: Manga{ID: "manga-id", Type: "manga", Attributes: MangaAttributes{Title: LocalisedStrings{Values: map[string]string{"en": "Solo"}}}}}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()
	r, err := c.Manga.GetMangaContext(t.Context(), "manga-id", nil)
	require.NoError(t, err)
	assert.Equal(t, "manga-id", r.Data.ID)
}

func TestManga_GetRandom_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/manga/random", r.URL.Path)
		_ = json.NewEncoder(w).Encode(SingleResponse[Manga]{Result: "ok", Data: Manga{ID: "rand"}})
	})
	defer srv.Close()
	r, err := c.Manga.GetRandomMangaContext(t.Context(), nil)
	require.NoError(t, err)
	assert.Equal(t, "rand", r.Data.ID)
}

func TestManga_GetAggregate_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/aggregate")
		agg := MangaAggregate{Result: "ok", Volumes: map[string]AggregateVolume{"1": {Volume: "1", Count: 1}}}
		_ = json.NewEncoder(w).Encode(agg)
	})
	defer srv.Close()
	agg, err := c.Manga.GetAggregateContext(t.Context(), "manga-id", nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", agg.Result)
}

func TestManga_GetMangaFeed_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/feed")
		list := ListResponse[Chapter]{Result: "ok", Data: []Chapter{{ID: "ch1"}}}
		_ = json.NewEncoder(w).Encode(list)
	})
	defer srv.Close()
	list, err := c.Manga.GetMangaFeedContext(t.Context(), "manga-id", nil)
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
}

func TestManga_CheckIfMangaFollowed_True(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	ok, err := c.Manga.CheckIfMangaFollowedContext(t.Context(), "id")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestManga_CheckIfMangaFollowed_False404(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Not Found", Detail: "404"}}})
	})
	defer srv.Close()
	ok, err := c.Manga.CheckIfMangaFollowedContext(t.Context(), "id")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestManga_ToggleFollow_FollowAndUnfollow(t *testing.T) {
	cases := []struct {
		follow bool
		method string
	}{
		{true, http.MethodPost},
		{false, http.MethodDelete},
	}
	for _, tc := range cases {
		c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, tc.method, r.Method)
			_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
		})
		_, err := c.Manga.ToggleMangaFollowStatusContext(t.Context(), "id", tc.follow)
		require.NoError(t, err)
		srv.Close()
	}
}

func TestManga_GetReadingStatus(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/manga/status", r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{"result": "ok", "statuses": map[string]string{"manga-id": "reading"}})
	})
	defer srv.Close()
	statuses, err := c.Manga.GetReadingStatusContext(t.Context(), nil)
	require.NoError(t, err)
	assert.Equal(t, "reading", statuses["manga-id"])
}
