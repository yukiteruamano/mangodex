package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaChaptersPath    = "manga/%s/feed"
	MangaReadMarkersPath = "manga/%s/read"
	ChapterPath          = "chapter/%s"
	ChapterListPath      = "chapter"
)

// ChapterService provides chapter services.
type ChapterService service

// ChapterList is a paginated list of chapters.
type ChapterList = ListResponse[Chapter]

// Chapter holds chapter information.
type Chapter struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
}

// GetTitle returns chapter title.
func (c *Chapter) GetTitle() string { return c.Attributes.Title }

// GetChapterNum returns chapter number or "-".
func (c *Chapter) GetChapterNum() string {
	if c.Attributes.Chapter != nil {
		return *c.Attributes.Chapter
	}
	return "-"
}

// ChapterAttributes holds chapter attributes.
type ChapterAttributes struct {
	Title              string  `json:"title"`
	Volume             *string `json:"volume"`
	Chapter            *string `json:"chapter"`
	TranslatedLanguage string  `json:"translatedLanguage"`
	Uploader           string  `json:"uploader"`
	ExternalURL        *string `json:"externalUrl"`
	Version            int     `json:"version"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
	PublishAt          string  `json:"publishAt"`
	Pages              int     `json:"pages"`
}

// GetMangaChapters returns chapters for a manga (feed).
func (s *ChapterService) GetMangaChapters(id string, params url.Values) (*ChapterList, error) {
	return s.GetMangaChaptersContext(context.Background(), id, params)
}

// GetMangaChaptersContext returns feed with context.
func (s *ChapterService) GetMangaChaptersContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaChaptersPath, id), params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetChapter returns a single chapter.
func (s *ChapterService) GetChapter(id string, params url.Values) (*SingleResponse[Chapter], error) {
	return s.GetChapterContext(context.Background(), id, params)
}

// GetChapterContext returns chapter with context.
func (s *ChapterService) GetChapterContext(ctx context.Context, id string, params url.Values) (*SingleResponse[Chapter], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(ChapterPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Chapter]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// SearchChapters searches chapters.
func (s *ChapterService) SearchChapters(params url.Values) (*ChapterList, error) {
	return s.SearchChaptersContext(context.Background(), params)
}

// SearchChaptersContext searches chapters with context.
func (s *ChapterService) SearchChaptersContext(ctx context.Context, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(ChapterListPath, params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// ChapterReadMarkers holds read marker response.
type ChapterReadMarkers struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (r *ChapterReadMarkers) GetResult() string { return r.Result }

// GetReadMangaChapters returns read chapter IDs for a manga.
func (s *ChapterService) GetReadMangaChapters(id string) (*ChapterReadMarkers, error) {
	return s.GetReadMangaChaptersContext(context.Background(), id)
}

// GetReadMangaChaptersContext returns read markers with context.
func (s *ChapterService) GetReadMangaChaptersContext(ctx context.Context, id string) (*ChapterReadMarkers, error) {
	u, err := s.client.buildURL(fmt.Sprintf(MangaReadMarkersPath, id))
	if err != nil {
		return nil, err
	}
	var r ChapterReadMarkers
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// SetReadUnreadMangaChapters marks chapters read/unread.
func (s *ChapterService) SetReadUnreadMangaChapters(id string, read, unRead []string) (*Response, error) {
	return s.SetReadUnreadMangaChaptersContext(context.Background(), id, read, unRead)
}

// SetReadUnreadMangaChaptersContext marks with context.
func (s *ChapterService) SetReadUnreadMangaChaptersContext(ctx context.Context, id string, read, unRead []string) (*Response, error) {
	u, err := s.client.buildURL(fmt.Sprintf(MangaReadMarkersPath, id))
	if err != nil {
		return nil, err
	}
	req := map[string][]string{
		"chapterIdsRead":   read,
		"chapterIdsUnread": unRead,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("mangodex: marshal read markers: %w", err)
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewReader(rBytes), &r); err != nil {
		return nil, err
	}
	return &r, nil
}
