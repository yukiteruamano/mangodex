package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	AuthorListPath = "author"
	AuthorPath     = "author/%s"
)

// AuthorService provides author/artist services.
type AuthorService service

// AuthorList is a paginated list of authors.
type AuthorList = ListResponse[Author]

// Author holds information on an author/artist.
type Author struct {
	ID            string           `json:"id"`
	Type          string           `json:"type"`
	Attributes    AuthorAttributes `json:"attributes"`
	Relationships []Relationship   `json:"relationships"`
}

// AuthorAttributes holds attributes for an author.
type AuthorAttributes struct {
	Name      string           `json:"name"`
	ImageURL  *string          `json:"imageUrl"`
	Biography LocalisedStrings `json:"biography"`
	Version   int              `json:"version"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
}

// SearchAuthors searches authors.
func (s *AuthorService) SearchAuthors(params url.Values) (*AuthorList, error) {
	return s.SearchAuthorsContext(context.Background(), params)
}

// SearchAuthorsContext searches authors with context.
func (s *AuthorService) SearchAuthorsContext(ctx context.Context, params url.Values) (*AuthorList, error) {
	urlStr, err := s.client.buildURLWithParams(AuthorListPath, params)
	if err != nil {
		return nil, err
	}
	var l AuthorList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetAuthor returns a single author by ID.
func (s *AuthorService) GetAuthor(id string, params url.Values) (*SingleResponse[Author], error) {
	return s.GetAuthorContext(context.Background(), id, params)
}

// GetAuthorContext returns author with context.
func (s *AuthorService) GetAuthorContext(ctx context.Context, id string, params url.Values) (*SingleResponse[Author], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(AuthorPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[Author]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
