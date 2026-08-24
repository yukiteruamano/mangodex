package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	StatisticsMangaPath        = "statistics/manga/%s"
	StatisticsChapterPath      = "statistics/chapter/%s"
	StatisticsGroupPath        = "statistics/group/%s"
	StatisticsMangaBatchPath   = "statistics/manga"
	StatisticsChapterBatchPath = "statistics/chapter"
	StatisticsGroupBatchPath   = "statistics/group"
)

// StatisticsService provides statistics services.
type StatisticsService service

// MangaStatistics holds manga statistics.
type MangaStatistics struct {
	Result     string                    `json:"result"`
	Statistics map[string]MangaStatistic `json:"statistics"`
}

func (r *MangaStatistics) GetResult() string { return r.Result }

type MangaStatistic struct {
	Comments *StatisticComments `json:"comments"`
	Rating   *StatisticRating   `json:"rating"`
	Follows  *int               `json:"follows"`
}

type StatisticComments struct {
	ThreadID     int `json:"threadId"`
	RepliesCount int `json:"repliesCount"`
}

type StatisticRating struct {
	Average      float64        `json:"average"`
	Bayesian     float64        `json:"bayesian"`
	Distribution map[string]int `json:"distribution"`
}

// GetMangaStatistics returns statistics for a manga.
func (s *StatisticsService) GetMangaStatistics(id string) (*MangaStatistics, error) {
	return s.GetMangaStatisticsContext(context.Background(), id)
}

// GetMangaStatisticsContext returns statistics with context.
func (s *StatisticsService) GetMangaStatisticsContext(ctx context.Context, id string) (*MangaStatistics, error) {
	u, err := s.client.buildURL(fmt.Sprintf(StatisticsMangaPath, id))
	if err != nil {
		return nil, err
	}
	var r MangaStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMangaStatisticsBatch returns statistics for multiple manga (comma separated ids in query).
func (s *StatisticsService) GetMangaStatisticsBatch(params url.Values) (*MangaStatistics, error) {
	return s.GetMangaStatisticsBatchContext(context.Background(), params)
}

func (s *StatisticsService) GetMangaStatisticsBatchContext(ctx context.Context, params url.Values) (*MangaStatistics, error) {
	urlStr, err := s.client.buildURLWithParams(StatisticsMangaBatchPath, params)
	if err != nil {
		return nil, err
	}
	var r MangaStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ChapterStatistics holds chapter statistics.
type ChapterStatistics struct {
	Result     string                      `json:"result"`
	Statistics map[string]ChapterStatistic `json:"statistics"`
}

func (r *ChapterStatistics) GetResult() string { return r.Result }

type ChapterStatistic struct {
	Comments *StatisticComments `json:"comments"`
}

// GetChapterStatistics returns statistics for a chapter.
func (s *StatisticsService) GetChapterStatistics(id string) (*ChapterStatistics, error) {
	return s.GetChapterStatisticsContext(context.Background(), id)
}

func (s *StatisticsService) GetChapterStatisticsContext(ctx context.Context, id string) (*ChapterStatistics, error) {
	u, err := s.client.buildURL(fmt.Sprintf(StatisticsChapterPath, id))
	if err != nil {
		return nil, err
	}
	var r ChapterStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetChapterStatisticsBatch returns batch chapter stats.
func (s *StatisticsService) GetChapterStatisticsBatch(params url.Values) (*ChapterStatistics, error) {
	return s.GetChapterStatisticsBatchContext(context.Background(), params)
}

// GetChapterStatisticsBatchContext returns batch with context.
func (s *StatisticsService) GetChapterStatisticsBatchContext(ctx context.Context, params url.Values) (*ChapterStatistics, error) {
	urlStr, err := s.client.buildURLWithParams(StatisticsChapterBatchPath, params)
	if err != nil {
		return nil, err
	}
	var r ChapterStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GroupStatistics holds group statistics.
type GroupStatistics struct {
	Result     string                    `json:"result"`
	Statistics map[string]GroupStatistic `json:"statistics"`
}

func (r *GroupStatistics) GetResult() string { return r.Result }

type GroupStatistic struct {
	Comments *StatisticComments `json:"comments"`
}

// GetGroupStatistics returns group stats.
func (s *StatisticsService) GetGroupStatistics(id string) (*GroupStatistics, error) {
	return s.GetGroupStatisticsContext(context.Background(), id)
}

// GetGroupStatisticsContext returns with context.
func (s *StatisticsService) GetGroupStatisticsContext(ctx context.Context, id string) (*GroupStatistics, error) {
	u, err := s.client.buildURL(fmt.Sprintf(StatisticsGroupPath, id))
	if err != nil {
		return nil, err
	}
	var r GroupStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetGroupStatisticsBatch returns batch group stats.
func (s *StatisticsService) GetGroupStatisticsBatch(params url.Values) (*GroupStatistics, error) {
	return s.GetGroupStatisticsBatchContext(context.Background(), params)
}

// GetGroupStatisticsBatchContext returns batch with context.
func (s *StatisticsService) GetGroupStatisticsBatchContext(ctx context.Context, params url.Values) (*GroupStatistics, error) {
	urlStr, err := s.client.buildURLWithParams(StatisticsGroupBatchPath, params)
	if err != nil {
		return nil, err
	}
	var r GroupStatistics
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
