package mangodex

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtHome_NewMDHomeClient_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/at-home/server/")
		resp := MDHomeServerResponse{Result: "ok", BaseURL: "https://uploads.mangadex.org", Chapter: ChaptersData{Hash: "abc", Data: []string{"p1.jpg", "p2.jpg"}, DataSaver: []string{"p1s.jpg"}}}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()
	client, err := c.AtHome.NewMDHomeClientContext(t.Context(), "chapter-id", "data", false)
	require.NoError(t, err)
	assert.Equal(t, "https://uploads.mangadex.org", client.baseURL)
	assert.Len(t, client.Pages, 2)
	assert.Equal(t, "p1.jpg", client.Pages[0])
}

func TestAtHome_NewMDHomeClient_DataSaver(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		resp := MDHomeServerResponse{Result: "ok", BaseURL: "https://uploads.mangadex.org", Chapter: ChaptersData{Hash: "abc", Data: []string{"p1.jpg"}, DataSaver: []string{"p1s.jpg", "p2s.jpg"}}}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()
	client, err := c.AtHome.NewMDHomeClientContext(t.Context(), "ch", "data-saver", true)
	require.NoError(t, err)
	assert.Equal(t, []string{"p1s.jpg", "p2s.jpg"}, client.Pages)
	// check forcePort443 param was sent
	_ = srv
}

func TestAtHome_GetChapterPage_Success(t *testing.T) {
	// Mock image server
	imgSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Cache", "HIT")
		_, _ = w.Write([]byte("imagedata"))
	}))
	defer imgSrv.Close()
	// Mock report server via MDHomeClient reportURL override
	reportSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer reportSrv.Close()

	c := &MDHomeClient{
		client:    &http.Client{},
		baseURL:   imgSrv.URL,
		quality:   "data",
		hash:      "hash",
		Pages:     []string{"page.jpg"},
		reportURL: reportSrv.URL,
	}
	data, err := c.GetChapterPageWithContext(t.Context(), "page.jpg")
	require.NoError(t, err)
	assert.Equal(t, "imagedata", string(data))
}

func TestAtHome_GetChapterPage_ErrorStatus(t *testing.T) {
	imgSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer imgSrv.Close()
	reportSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer reportSrv.Close()
	c := &MDHomeClient{
		client:    &http.Client{},
		baseURL:   imgSrv.URL,
		quality:   "data",
		hash:      "h",
		reportURL: reportSrv.URL,
	}
	_, err := c.GetChapterPageWithContext(t.Context(), "p.jpg")
	require.Error(t, err)
}

func TestAtHome_ForcePort443_Query(t *testing.T) {
	var gotQuery string
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("forcePort443")
		_ = json.NewEncoder(w).Encode(MDHomeServerResponse{Result: "ok", BaseURL: "https://x", Chapter: ChaptersData{Hash: "h", Data: []string{"a"}}})
	})
	defer srv.Close()
	_, err := c.AtHome.NewMDHomeClientContext(t.Context(), "ch", "data", true)
	require.NoError(t, err)
	assert.Equal(t, "true", gotQuery)
}
