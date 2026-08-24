package mangodex

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*DexClient, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c := NewDexClient(WithBaseURL(srv.URL), WithTimeout(5*time.Second))
	return c, srv
}

func TestClient_Options(t *testing.T) {
	c := NewDexClient(WithBaseURL("https://example.com/"), WithUserAgent("test/1.0"), WithTimeout(2*time.Second))
	assert.Equal(t, "https://example.com", c.baseURL)
	assert.Equal(t, "test/1.0", c.userAgent)
	assert.Equal(t, 2*time.Second, c.client.Timeout)

	c2 := NewDexClient(WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
	assert.Equal(t, 10*time.Second, c2.client.Timeout)
}

func TestClient_Request_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "mangodex-go/2.0 (+https://github.com/yukiteruamano/mangodex)", r.Header.Get("User-Agent"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	})
	defer srv.Close()

	resp, err := c.Request(t.Context(), http.MethodGet, srv.URL+"/ping", nil)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestClient_Request_Non2xx_WithErrorResponse(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Result: "error",
			Errors: []Error{{Title: "Bad Request", Detail: "invalid param"}},
		})
	})
	defer srv.Close()

	_, err := c.Request(t.Context(), http.MethodGet, srv.URL+"/manga", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Bad Request")
}

func TestClient_Request_Non2xx_WithRequestID(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Not Found", Detail: "missing"}}})
	})
	defer srv.Close()
	_, err := c.Request(t.Context(), http.MethodGet, srv.URL+"/notfound", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "req-123")
}

func TestClient_RequestAndDecode_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"1","type":"manga","attributes":{"title":{"en":"Test"},"altTitles":{},"description":{},"isLocked":false,"links":{},"originalLanguage":"ja","state":"published","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","tags":[]}}}`))
	})
	defer srv.Close()
	var r SingleResponse[Manga]
	err := c.RequestAndDecode(t.Context(), http.MethodGet, srv.URL+"/manga/1", nil, &r)
	require.NoError(t, err)
	assert.Equal(t, "ok", r.Result)
	assert.Equal(t, "1", r.Data.ID)
}

func TestClient_RequestAndDecode_InvalidJSON(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	})
	defer srv.Close()
	var r Response
	err := c.RequestAndDecode(t.Context(), http.MethodGet, srv.URL+"/bad", nil, &r)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}

func TestClient_BuildURL(t *testing.T) {
	c := NewDexClient(WithBaseURL("https://api.example.com/"))
	u, err := c.buildURL("manga")
	require.NoError(t, err)
	assert.Equal(t, "https://api.example.com/manga", u.String())

	u2, err := c.buildURL("/manga/123")
	require.NoError(t, err)
	assert.Equal(t, "https://api.example.com/manga/123", u2.String())

	str, err := c.buildURLWithParams("manga", url.Values{"limit": {"10"}})
	require.NoError(t, err)
	assert.Contains(t, str, "limit=10")
}
