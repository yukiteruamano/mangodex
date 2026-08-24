package mangodex

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	MangaListPath            = "manga"
	MangaPath                = "manga/%s"
	MangaAggregatePath       = "manga/%s/aggregate"
	MangaFeedPath            = "manga/%s/feed"
	CheckIfMangaFollowedPath = "user/follows/manga/%s"
	ToggleMangaFollowPath    = "manga/%s/follow"
	MangaStatusPath          = "manga/status"
	MangaRandomPath          = "manga/random"
	MangaRecommendationPath  = "manga/%s/recommendation"
	MangaDraftsPath          = "manga/draft"
	MangaDraftPath           = "manga/draft/%s"
	MangaIdStatusPath        = "manga/%s/status"
)

// MangaService provides manga services.
type MangaService service

// MangaList is a paginated list of manga (alias for generic).
type MangaList = ListResponse[Manga]

// Manga holds information on a manga.
type Manga struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

// GetTitle returns title for langCode with fallback to AltTitles.
func (m *Manga) GetTitle(langCode string) string {
	return cmp.Or(m.Attributes.Title.GetLocalString(langCode), m.Attributes.AltTitles.GetLocalString(langCode))
}

// GetDescription returns description for langCode.
func (m *Manga) GetDescription(langCode string) string {
	return m.Attributes.Description.GetLocalString(langCode)
}

// MangaAttributes holds attributes for a manga.
type MangaAttributes struct {
	Title                        LocalisedStrings  `json:"title"`
	AltTitles                    LocalisedStrings  `json:"altTitles"`
	Description                  LocalisedStrings  `json:"description"`
	IsLocked                     bool              `json:"isLocked"`
	Links                        map[string]string `json:"links"`
	OriginalLanguage             string            `json:"originalLanguage"`
	LastVolume                   *string           `json:"lastVolume"`
	LastChapter                  *string           `json:"lastChapter"`
	PublicationDemographic       *string           `json:"publicationDemographic"`
	Status                       *string           `json:"status"`
	Year                         *int              `json:"year"`
	ContentRating                *string           `json:"contentRating"`
	Tags                         []Tag             `json:"tags"`
	State                        string            `json:"state"`
	Version                      int               `json:"version"`
	CreatedAt                    string            `json:"createdAt"`
	UpdatedAt                    string            `json:"updatedAt"`
	AvailableTranslatedLanguages []string          `json:"availableTranslatedLanguages"`
	LatestUploadedChapter        *string           `json:"latestUploadedChapter"`
}

// MangaAggregate holds volume/chapter aggregate.
type MangaAggregate struct {
	Result  string                     `json:"result"`
	Volumes map[string]AggregateVolume `json:"volumes"`
}

func (r *MangaAggregate) GetResult() string { return r.Result }

type AggregateVolume struct {
	Volume   string                      `json:"volume"`
	Count    int                         `json:"count"`
	Chapters map[string]AggregateChapter `json:"chapters"`
}

type AggregateChapter struct {
	Chapter string   `json:"chapter"`
	ID      string   `json:"id"`
	Others  []string `json:"others"`
	Count   int      `json:"count"`
}

// GetMangaList searches manga.
func (s *MangaService) GetMangaList(params url.Values) (*MangaList, error) {
	return s.GetMangaListContext(context.Background(), params)
}

// GetMangaListContext searches manga with context.
func (s *MangaService) GetMangaListContext(ctx context.Context, params url.Values) (*MangaList, error) {
	urlStr, err := s.client.buildURLWithParams(MangaListPath, params)
	if err != nil {
		return nil, err
	}
	var l MangaList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetManga returns a single manga by ID.
func (s *MangaService) GetManga(id string, params url.Values) (*SingleResponse[Manga], error) {
	return s.GetMangaContext(context.Background(), id, params)
}

// GetMangaContext returns a single manga with context.
func (s *MangaService) GetMangaContext(ctx context.Context, id string, params url.Values) (*SingleResponse[Manga], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Manga]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetRandomManga returns a random manga.
func (s *MangaService) GetRandomManga(params url.Values) (*SingleResponse[Manga], error) {
	return s.GetRandomMangaContext(context.Background(), params)
}

// GetRandomMangaContext returns random manga with context.
func (s *MangaService) GetRandomMangaContext(ctx context.Context, params url.Values) (*SingleResponse[Manga], error) {
	urlStr, err := s.client.buildURLWithParams(MangaRandomPath, params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Manga]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetAggregate returns volume/chapter aggregate for a manga.
func (s *MangaService) GetAggregate(id string, params url.Values) (*MangaAggregate, error) {
	return s.GetAggregateContext(context.Background(), id, params)
}

// GetAggregateContext returns aggregate with context.
func (s *MangaService) GetAggregateContext(ctx context.Context, id string, params url.Values) (*MangaAggregate, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaAggregatePath, id), params)
	if err != nil {
		return nil, err
	}
	var r MangaAggregate
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMangaFeed returns chapters for a manga (alias to ChapterService but kept for convenience).
func (s *MangaService) GetMangaFeed(id string, params url.Values) (*ChapterList, error) {
	return s.GetMangaFeedContext(context.Background(), id, params)
}

// GetMangaFeedContext returns feed with context.
func (s *MangaService) GetMangaFeedContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaFeedPath, id), params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// CheckIfMangaFollowed checks if the user follows a manga.
func (s *MangaService) CheckIfMangaFollowed(id string) (bool, error) {
	return s.CheckIfMangaFollowedContext(context.Background(), id)
}

// CheckIfMangaFollowedContext checks follow with context. Returns false on 404.
func (s *MangaService) CheckIfMangaFollowedContext(ctx context.Context, id string) (bool, error) {
	u, err := s.client.buildURL(fmt.Sprintf(CheckIfMangaFollowedPath, id))
	if err != nil {
		return false, err
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		if errors.Is(err, ErrAPI) && isNotFound(err) {
			return false, nil
		}
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "404") || strings.Contains(msg, "not found")
}

// ToggleMangaFollowStatus toggles follow status.
func (s *MangaService) ToggleMangaFollowStatus(id string, toFollow bool) (*Response, error) {
	return s.ToggleMangaFollowStatusContext(context.Background(), id, toFollow)
}

// ToggleMangaFollowStatusContext toggles with context.
func (s *MangaService) ToggleMangaFollowStatusContext(ctx context.Context, id string, toFollow bool) (*Response, error) {
	u, err := s.client.buildURL(fmt.Sprintf(ToggleMangaFollowPath, id))
	if err != nil {
		return nil, err
	}
	method := http.MethodPost
	if !toFollow {
		method = http.MethodDelete
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, method, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetReadingStatus returns manga reading status for current user.
func (s *MangaService) GetReadingStatus(params url.Values) (map[string]string, error) {
	return s.GetReadingStatusContext(context.Background(), params)
}

// statusResponse is internal for decoding manga/status.
type statusResponse struct {
	Result   string            `json:"result"`
	Statuses map[string]string `json:"statuses"`
}

func (s *statusResponse) GetResult() string { return s.Result }

// GetReadingStatusContext returns reading status with context.
func (s *MangaService) GetReadingStatusContext(ctx context.Context, params url.Values) (map[string]string, error) {
	urlStr, err := s.client.buildURLWithParams(MangaStatusPath, params)
	if err != nil {
		return nil, err
	}
	var sr statusResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &sr); err != nil {
		return nil, err
	}
	return sr.Statuses, nil
}

// MangaRecommendationList is alias for recommendations.
type MangaRecommendationList = ListResponse[Manga]

// GetMangaRecommendations returns recommendations for a manga.
func (s *MangaService) GetMangaRecommendations(id string, params url.Values) (*MangaRecommendationList, error) {
	return s.GetMangaRecommendationsContext(context.Background(), id, params)
}

// GetMangaRecommendationsContext returns recommendations with context.
func (s *MangaService) GetMangaRecommendationsContext(ctx context.Context, id string, params url.Values) (*MangaRecommendationList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaRecommendationPath, id), params)
	if err != nil {
		return nil, err
	}
	var l MangaRecommendationList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// MangaDraftList is draft list.
type MangaDraftList = ListResponse[Manga]

// GetMangaDrafts returns drafts. params may be nil; spec defines 0 query, extra values are ignored (future-proof).
func (s *MangaService) GetMangaDrafts(params url.Values) (*MangaDraftList, error) {
	return s.GetMangaDraftsContext(context.Background(), params)
}

// GetMangaDraftsContext returns drafts with context. params may be nil.
func (s *MangaService) GetMangaDraftsContext(ctx context.Context, params url.Values) (*MangaDraftList, error) {
	urlStr, err := s.client.buildURLWithParams(MangaDraftsPath, params)
	if err != nil {
		return nil, err
	}
	var l MangaDraftList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetMangaDraft returns a single draft. params may be nil; spec defines 0 query but accepts includes[] for future.
func (s *MangaService) GetMangaDraft(id string, params url.Values) (*SingleResponse[Manga], error) {
	return s.GetMangaDraftContext(context.Background(), id, params)
}

// GetMangaDraftContext returns draft with context. params may be nil.
func (s *MangaService) GetMangaDraftContext(ctx context.Context, id string, params url.Values) (*SingleResponse[Manga], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaDraftPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Manga]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// mangaIDStatusResponse is response for single manga status.
type mangaIDStatusResponse struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

func (r *mangaIDStatusResponse) GetResult() string { return r.Result }

// GetMangaIDStatus returns status for a single manga.
func (s *MangaService) GetMangaIDStatus(id string) (string, error) {
	return s.GetMangaIDStatusContext(context.Background(), id)
}

// GetMangaIDStatusContext returns single status with context.
func (s *MangaService) GetMangaIDStatusContext(ctx context.Context, id string) (string, error) {
	u, err := s.client.buildURL(fmt.Sprintf(MangaIdStatusPath, id))
	if err != nil {
		return "", err
	}
	var ss mangaIDStatusResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &ss); err != nil {
		return "", err
	}
	return ss.Status, nil
}
