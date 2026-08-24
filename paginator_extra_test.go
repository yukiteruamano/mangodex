package mangodex

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginator_All(t *testing.T) {
	t.Parallel()
	calls := 0
	fetch := func(ctx context.Context, limit, offset int) (*ListResponse[Manga], error) {
		calls++
		if offset >= 2 {
			return &ListResponse[Manga]{Result: "ok", Data: []Manga{}, Total: 2}, nil
		}
		return &ListResponse[Manga]{Result: "ok", Data: []Manga{{ID: "m1"}, {ID: "m2"}}, Total: 2}, nil
	}
	p := NewPaginator[Manga](fetch)
	var collected []Manga
	for data, err := range p.All(t.Context()) {
		require.NoError(t, err)
		collected = append(collected, data...)
		if !p.HasMore() {
			break
		}
	}
	assert.Len(t, collected, 2)
	assert.Equal(t, 1, calls) // first call returns 2 items, total 2, so only 1 fetch
	p.Reset()
	assert.True(t, p.HasMore())
}

func TestWithLogger(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c := NewDexClient(WithLogger(logger), WithBaseURL("https://example.com"))
	assert.NotNil(t, c.logger)
	assert.Equal(t, "https://example.com", c.BaseURL())
	// ensure logger doesn't break request
	c2, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"result":"error","errors":[{"title":"e","detail":"d"}]}`))
	})
	defer srv.Close()
	c2.logger = logger
	_, err := c2.Manga.GetMangaListContext(t.Context(), nil)
	require.Error(t, err)
}

func TestDefaultTransportOnce(t *testing.T) {
	t.Parallel()
	tr1 := defaultTransport()
	tr2 := defaultTransport()
	assert.Same(t, tr1, tr2)
	assert.Equal(t, 100, tr1.MaxIdleConns)
}

func TestBuildURL_JoinPath(t *testing.T) {
	t.Parallel()
	c := NewDexClient(WithBaseURL("https://api.example.com///"))
	u, err := c.buildURL("manga")
	require.NoError(t, err)
	assert.Equal(t, "https://api.example.com/manga", u.String())
	u2, err := c.buildURL("/manga/123")
	require.NoError(t, err)
	assert.Equal(t, "https://api.example.com/manga/123", u2.String())
}

func TestClient_HeaderClone(t *testing.T) {
	t.Parallel()
	c2, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value", r.Header.Get("X-Custom"))
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	})
	defer srv.Close()
	c2.header.Set("X-Custom", "value")
	// ensure header not mutated by request
	_, _ = c2.Request(t.Context(), http.MethodGet, srv.URL+"/ping", nil)
	assert.Equal(t, "value", c2.header.Get("X-Custom"))
}
