package mangodex

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseAPI   = "https://api.mangadex.org"
	defaultTimeout   = 15 * time.Second
	defaultUserAgent = "mangodex-go/2.0 (+https://github.com/yukiteruamano/mangodex)"
)

// DefaultBaseAPI is the default MangaDex API base URL.
const DefaultBaseAPI = defaultBaseAPI

// DexClient is the MangaDex API client.
type DexClient struct {
	mu           sync.RWMutex
	client       *http.Client
	header       http.Header
	baseURL      string
	refreshToken string
	userAgent    string
	logger       *slog.Logger

	common service

	// Services
	Auth            *AuthService
	Manga           *MangaService
	Chapter         *ChapterService
	User            *UserService
	AtHome          *AtHomeService
	Cover           *CoverService
	Author          *AuthorService
	ScanlationGroup *ScanlationGroupService
	CustomList      *CustomListService
	Feed            *FeedService
	Statistics      *StatisticsService
	Tag             *TagService
	Infrastructure  *InfrastructureService
	Rating          *RatingService
	Report          *ReportService
	Relation        *RelationService
	Settings        *SettingsService
	Upload          *UploadService
	ApiClient       *ApiClientService
}

// service is a wrapper for DexClient.
type service struct {
	client *DexClient
}

// Option configures a DexClient.
type Option func(*DexClient)

// WithBaseURL overrides the API base URL (useful for httptest).
func WithBaseURL(u string) Option {
	return func(c *DexClient) {
		c.baseURL = strings.TrimRight(u, "/")
	}
}

// WithHTTPClient injects a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *DexClient) {
		if hc != nil {
			c.client = hc
		}
	}
}

// WithTimeout sets the http.Client timeout (if no custom client provided, creates one).
func WithTimeout(d time.Duration) Option {
	return func(c *DexClient) {
		if c.client == nil {
			c.client = &http.Client{Timeout: d}
		} else {
			c.client.Timeout = d
		}
	}
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *DexClient) {
		c.userAgent = ua
		c.header.Set("User-Agent", ua)
	}
}

// WithLogger sets a structured logger (slog) for the client.
func WithLogger(l *slog.Logger) Option {
	return func(c *DexClient) {
		c.logger = l
	}
}

// defaultTransport returns a tuned transport for scrapper massive (desktop).
// perf: tuned for 50-200 concurrent goroutines - MaxIdleConnsPerHost 20 vs default 2 reduces TLS allocs.
func defaultTransport() *http.Transport {
	return defaultTransportOnce()
}

var defaultTransportOnce = sync.OnceValue(func() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}
})

// NewDexClient creates an anonymous client. Options are applied after defaults.
func NewDexClient(opts ...Option) *DexClient {
	hc := &http.Client{Timeout: defaultTimeout, Transport: defaultTransport()}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", defaultUserAgent)
	header.Set("Accept", "application/json")

	dex := &DexClient{
		client:    hc,
		header:    header,
		baseURL:   defaultBaseAPI,
		userAgent: defaultUserAgent,
	}
	for _, o := range opts {
		o(dex)
	}
	// Ensure baseURL fallback if empty after options
	dex.baseURL = cmp.Or(dex.baseURL, defaultBaseAPI)
	dex.common.client = dex

	dex.Auth = (*AuthService)(&dex.common)
	dex.Manga = (*MangaService)(&dex.common)
	dex.Chapter = (*ChapterService)(&dex.common)
	dex.User = (*UserService)(&dex.common)
	dex.AtHome = (*AtHomeService)(&dex.common)
	dex.Cover = (*CoverService)(&dex.common)
	dex.Author = (*AuthorService)(&dex.common)
	dex.ScanlationGroup = (*ScanlationGroupService)(&dex.common)
	dex.CustomList = (*CustomListService)(&dex.common)
	dex.Feed = (*FeedService)(&dex.common)
	dex.Statistics = (*StatisticsService)(&dex.common)
	dex.Tag = (*TagService)(&dex.common)
	dex.Infrastructure = (*InfrastructureService)(&dex.common)
	dex.Rating = (*RatingService)(&dex.common)
	dex.Report = (*ReportService)(&dex.common)
	dex.Relation = (*RelationService)(&dex.common)
	dex.Settings = (*SettingsService)(&dex.common)
	dex.Upload = (*UploadService)(&dex.common)
	dex.ApiClient = (*ApiClientService)(&dex.common)

	return dex
}

// New is an alias for NewDexClient.
func New(opts ...Option) *DexClient { return NewDexClient(opts...) }

// BaseURL returns the client's base URL (thread-safe).
func (c *DexClient) BaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL
}

// buildURL joins baseURL with path, returning parsed URL.
func (c *DexClient) buildURL(path string) (*url.URL, error) {
	c.mu.RLock()
	base := c.baseURL
	c.mu.RUnlock()
	base = cmp.Or(base, defaultBaseAPI)
	joined, err := url.JoinPath(base, path)
	if err != nil {
		return nil, fmt.Errorf("mangodex: invalid URL %q: %w", base+"/"+path, err)
	}
	u, err := url.Parse(joined)
	if err != nil {
		return nil, fmt.Errorf("mangodex: invalid URL %q: %w", joined, err)
	}
	return u, nil
}

// Request sends a request to the MangaDex API.
func (c *DexClient) Request(ctx context.Context, method, urlStr string, body io.Reader) (*http.Response, error) {
	c.mu.RLock()
	hdr := c.header.Clone()
	client := c.client
	c.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("mangodex: create request: %w", err)
	}
	req.Header = hdr

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mangodex: do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
		// Try to decode ErrorResponse, fallback to raw body
		var er ErrorResponse
		if decErr := json.NewDecoder(resp.Body).Decode(&er); decErr == nil && len(er.Errors) > 0 {
			reqID := resp.Header.Get("X-Request-ID")
			if c.logger != nil {
				c.logger.Error("mangodex: non-2xx", "status", resp.StatusCode, "reqID", reqID, "errors", er.GetErrors(), "url", urlStr)
			}
			if reqID != "" {
				return nil, fmt.Errorf("mangodex: non-2xx status %d (req %s): %s: %w", resp.StatusCode, reqID, er.GetErrors(), ErrAPI)
			}
			return nil, fmt.Errorf("mangodex: non-2xx status %d: %s: %w", resp.StatusCode, er.GetErrors(), ErrAPI)
		}
		// Fallback: read body as string
		reqID := resp.Header.Get("X-Request-ID")
		if c.logger != nil {
			c.logger.Error("mangodex: non-2xx", "status", resp.StatusCode, "reqID", reqID, "url", urlStr)
		}
		if reqID != "" {
			return nil, fmt.Errorf("mangodex: non-2xx status %d (req %s): %w", resp.StatusCode, reqID, ErrAPI)
		}
		return nil, fmt.Errorf("mangodex: non-2xx status %d: %w", resp.StatusCode, ErrAPI)
	}
	return resp, nil
}

// RequestAndDecode is a convenience wrapper that decodes the response JSON into rt.
func (c *DexClient) RequestAndDecode(ctx context.Context, method, urlStr string, body io.Reader, rt ResponseType) error {
	resp, err := c.Request(ctx, method, urlStr, body)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if err := json.NewDecoder(resp.Body).Decode(rt); err != nil {
		return fmt.Errorf("mangodex: decode response: %w", err)
	}
	return nil
}

// ErrAPI is a sentinel for API errors.
var ErrAPI = errors.New("mangadex api error")

// buildURLWithParams is a helper for services.
func (c *DexClient) buildURLWithParams(path string, params url.Values) (string, error) {
	u, err := c.buildURL(path)
	if err != nil {
		return "", err
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	return u.String(), nil
}
