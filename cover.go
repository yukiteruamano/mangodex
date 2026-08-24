package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	CoverListPath = "cover"
	CoverPath     = "cover/%s"
)

// CoverService provides cover art services.
type CoverService service

// CoverList is a paginated list of covers.
type CoverList = ListResponse[Cover]

// Cover holds cover art information.
type Cover struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    CoverAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

// CoverAttributes holds cover attributes.
type CoverAttributes struct {
	Volume      *string `json:"volume"`
	FileName    string  `json:"fileName"`
	Description *string `json:"description"`
	Locale      *string `json:"locale"`
	Version     int     `json:"version"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// SearchCovers searches cover art.
func (s *CoverService) SearchCovers(params url.Values) (*CoverList, error) {
	return s.SearchCoversContext(context.Background(), params)
}

// SearchCoversContext searches covers with context.
func (s *CoverService) SearchCoversContext(ctx context.Context, params url.Values) (*CoverList, error) {
	urlStr, err := s.client.buildURLWithParams(CoverListPath, params)
	if err != nil {
		return nil, err
	}
	var l CoverList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetCover returns a single cover by ID.
func (s *CoverService) GetCover(id string, params url.Values) (*SingleResponse[Cover], error) {
	return s.GetCoverContext(context.Background(), id, params)
}

// GetCoverContext returns cover with context.
func (s *CoverService) GetCoverContext(ctx context.Context, id string, params url.Values) (*SingleResponse[Cover], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(CoverPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Cover]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetCoverImageURL builds a direct cover image URL (not API).
// Example: https://uploads.mangadx.org/covers/{mangaId}/{fileName}.256.jpg
func GetCoverImageURL(mangaID, fileName, size string) string {
	return fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s%s", mangaID, fileName, size)
}
