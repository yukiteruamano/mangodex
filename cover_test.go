package mangodex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCover_SearchAndGet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cover" {
			list := ListResponse[Cover]{Result: "ok", Data: []Cover{{ID: "c1", Attributes: CoverAttributes{FileName: "cover.jpg"}}}}
			_ = json.NewEncoder(w).Encode(list)
			return
		}
		assert.Equal(t, "/cover/c1", r.URL.Path)
		_ = json.NewEncoder(w).Encode(SingleResponse[Cover]{Result: "ok", Data: Cover{ID: "c1"}})
	})
	defer srv.Close()
	list, err := c.Cover.SearchCoversContext(t.Context(), url.Values{"limit": {"10"}})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)

	single, err := c.Cover.GetCoverContext(t.Context(), "c1", nil)
	require.NoError(t, err)
	assert.Equal(t, "c1", single.Data.ID)
}

func TestCover_ImageURL(t *testing.T) {
	u := GetCoverImageURL("manga-id", "file.jpg", ".256.jpg")
	assert.Contains(t, u, "manga-id")
	assert.Contains(t, u, "file.jpg")
	u2 := GetCoverImageURL("mid", "f.jpg", "")
	assert.Equal(t, "https://uploads.mangadex.org/covers/mid/f.jpg", u2)
}
