package mangodex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChapter_GetMangaChapters_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/feed")
		list := ListResponse[Chapter]{Result: "ok", Data: []Chapter{{ID: "ch1", Attributes: ChapterAttributes{Title: "Ch1"}}}}
		_ = json.NewEncoder(w).Encode(list)
	})
	defer srv.Close()
	list, err := c.Chapter.GetMangaChaptersContext(t.Context(), "manga-id", url.Values{"limit": {"10"}})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
}

func TestChapter_GetChapter_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chapter/ch-id", r.URL.Path)
		_ = json.NewEncoder(w).Encode(SingleResponse[Chapter]{Result: "ok", Data: Chapter{ID: "ch-id"}})
	})
	defer srv.Close()
	r, err := c.Chapter.GetChapterContext(t.Context(), "ch-id", nil)
	require.NoError(t, err)
	assert.Equal(t, "ch-id", r.Data.ID)
}

func TestChapter_SearchChapters(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chapter", r.URL.Path)
		_ = json.NewEncoder(w).Encode(ListResponse[Chapter]{Result: "ok", Data: []Chapter{}})
	})
	defer srv.Close()
	_, err := c.Chapter.SearchChaptersContext(t.Context(), url.Values{"title": {"test"}})
	require.NoError(t, err)
}

func TestChapter_ReadMarkers_GetAndSet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_ = json.NewEncoder(w).Encode(ChapterReadMarkers{Result: "ok", Data: []string{"ch1"}})
			return
		}
		assert.Equal(t, http.MethodPost, r.Method)
		var body map[string][]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		assert.Contains(t, body["chapterIdsRead"], "ch1")
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	markers, err := c.Chapter.GetReadMangaChaptersContext(t.Context(), "manga-id")
	require.NoError(t, err)
	assert.Contains(t, markers.Data, "ch1")

	resp, err := c.Chapter.SetReadUnreadMangaChaptersContext(t.Context(), "manga-id", []string{"ch1"}, []string{"ch2"})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Result)
}

func TestChapter_GetMangaChapters_Error(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Internal", Detail: "err"}}})
	})
	defer srv.Close()
	_, err := c.Chapter.GetMangaChaptersContext(t.Context(), "mid", nil)
	require.Error(t, err)
}
