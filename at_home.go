package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	GetMDHomeURLPath = "at-home/server/%s"
	MDHomeReportURL  = "https://api.mangadex.network/report"
)

// AtHomeService provides MangaDex@Home services.
type AtHomeService service

// MDHomeServerResponse is the response for getting a server URL.
type MDHomeServerResponse struct {
	Result  string       `json:"result"`
	BaseURL string       `json:"baseUrl"`
	Chapter ChaptersData `json:"chapter"`
}

func (r *MDHomeServerResponse) GetResult() string { return r.Result }

// ChaptersData holds chapter page data.
type ChaptersData struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// MDHomeClient is a client for MangaDex@Home image servers.
type MDHomeClient struct {
	client    *http.Client
	baseURL   string
	quality   string
	hash      string
	Pages     []string
	reportURL string
}

// NewMDHomeClient gets a MangaDex@Home client for a chapter.
func (s *AtHomeService) NewMDHomeClient(chapterID string, quality string, forcePort443 bool) (*MDHomeClient, error) {
	return s.NewMDHomeClientContext(context.Background(), chapterID, quality, forcePort443)
}

// NewMDHomeClientContext creates a MDHomeClient with custom context.
func (s *AtHomeService) NewMDHomeClientContext(ctx context.Context, chapterID string, quality string, forcePort443 bool) (*MDHomeClient, error) {
	u, err := s.client.buildURL(fmt.Sprintf(GetMDHomeURLPath, chapterID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("forcePort443", strconv.FormatBool(forcePort443))
	u.RawQuery = q.Encode()

	var r MDHomeServerResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}

	pages := r.Chapter.Data
	if quality == "data-saver" {
		pages = r.Chapter.DataSaver
	}

	// perf: reuse DexClient transport for scrapper massive to enable keep-alive
	var transport http.RoundTripper
	if s.client.client != nil && s.client.client.Transport != nil {
		transport = s.client.client.Transport
	} else {
		transport = defaultTransport()
	}
	return &MDHomeClient{
		client:    &http.Client{Timeout: 30 * time.Second, Transport: transport},
		baseURL:   r.BaseURL,
		quality:   quality,
		hash:      r.Chapter.Hash,
		Pages:     pages,
		reportURL: MDHomeReportURL,
	}, nil
}

// GetChapterPage returns page data for a chapter.
func (c *MDHomeClient) GetChapterPage(filename string) ([]byte, error) {
	return c.GetChapterPageWithContext(context.Background(), filename)
}

// GetChapterPageWithContext returns page data with custom context.
func (c *MDHomeClient) GetChapterPageWithContext(ctx context.Context, filename string) ([]byte, error) {
	path, err := url.JoinPath(c.baseURL, c.quality, c.hash, filename)
	if err != nil {
		return nil, fmt.Errorf("mangodex: build page URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("mangodex: create page request: %w", err)
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		// Report failure without resp
		go c.reportAsync(context.WithoutCancel(ctx), path, false, 0, time.Since(start).Milliseconds(), false)
		return nil, fmt.Errorf("mangodex: fetch page: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		go c.reportAsync(context.WithoutCancel(ctx), path, false, 0, time.Since(start).Milliseconds(), strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT"))
		return nil, fmt.Errorf("mangodex: page non-2xx %d", resp.StatusCode)
	}

	// Limit to 25 MiB per page to avoid OOM
	limited := io.LimitReader(resp.Body, 25<<20)
	fileData, readErr := io.ReadAll(limited)
	duration := time.Since(start).Milliseconds()
	cached := strings.HasPrefix(resp.Header.Get("X-Cache"), "HIT")

	go c.reportAsync(context.WithoutCancel(ctx), path, readErr == nil, len(fileData), duration, cached)

	if readErr != nil {
		return nil, fmt.Errorf("mangodex: read page: %w", readErr)
	}
	return fileData, nil
}

func (c *MDHomeClient) reportAsync(ctx context.Context, url string, success bool, bytes int, duration int64, cached bool) {
	r := &reportPayload{
		URL:      url,
		Success:  success,
		Bytes:    bytes,
		Duration: duration,
		Cached:   cached,
	}
	resp, _ := c.reportContext(ctx, r)
	if resp != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

// reportPayload holds fields for reporting download result.
type reportPayload struct {
	URL      string `json:"url"`
	Success  bool   `json:"success"`
	Bytes    int    `json:"bytes"`
	Duration int64  `json:"duration"`
	Cached   bool   `json:"cached"`
}

func (c *MDHomeClient) reportContext(ctx context.Context, r *reportPayload) (*http.Response, error) {
	// Short timeout for report
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.reportURL, bytes.NewReader(rBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}
