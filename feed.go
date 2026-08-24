package mangodex

import (
	"context"
	"net/http"
	"net/url"
)

const (
	UserFollowsMangaFeedPath = "user/follows/manga/feed"
	ReadingHistoryPath       = "user/history"
	MangaReadBatchPath       = "manga/read"
)

// FeedService provides feed services.
type FeedService service

// MangaReadResponse holds batch read markers.
type MangaReadResponse struct {
	Result string              `json:"result"`
	Data   map[string][]string `json:"data"`
}

func (r *MangaReadResponse) GetResult() string { return r.Result }

// GetFollowedMangaFeed returns chapters from followed manga.
func (s *FeedService) GetFollowedMangaFeed(params url.Values) (*ChapterList, error) {
	return s.GetFollowedMangaFeedContext(context.Background(), params)
}

// GetFollowedMangaFeedContext returns feed with context.
func (s *FeedService) GetFollowedMangaFeedContext(ctx context.Context, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(UserFollowsMangaFeedPath, params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetReadingHistory returns reading history. params may be nil; spec defines 0 query, extra values are ignored.
func (s *FeedService) GetReadingHistory(params url.Values) (*ChapterList, error) {
	return s.GetReadingHistoryContext(context.Background(), params)
}

// GetReadingHistoryContext returns history with context. params may be nil.
func (s *FeedService) GetReadingHistoryContext(ctx context.Context, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(ReadingHistoryPath, params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetMangaReadHistory returns grouped read markers for manga ids batch.
func (s *FeedService) GetMangaReadHistory(params url.Values) (map[string][]string, error) {
	return s.GetMangaReadHistoryContext(context.Background(), params)
}

// GetMangaReadHistoryContext returns with context.
func (s *FeedService) GetMangaReadHistoryContext(ctx context.Context, params url.Values) (map[string][]string, error) {
	urlStr, err := s.client.buildURLWithParams(MangaReadBatchPath, params)
	if err != nil {
		return nil, err
	}
	var r MangaReadResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}
