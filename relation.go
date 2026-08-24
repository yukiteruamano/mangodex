package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaRelationPath = "manga/%s/relation"
)

// RelationService provides manga relation endpoints.
type RelationService service

// MangaRelationList holds relations for a manga.
type MangaRelationList struct {
	Result   string          `json:"result"`
	Response string          `json:"response"`
	Data     []MangaRelation `json:"data"`
}

func (r *MangaRelationList) GetResult() string { return r.Result }

// MangaRelation holds a relation.
type MangaRelation struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    MangaRelationAttributes `json:"attributes"`
	Relationships []Relationship          `json:"relationships"`
}

type MangaRelationAttributes struct {
	Relation string `json:"relation"`
	Version  int    `json:"version"`
}

// GetMangaRelations returns relations for a manga.
func (s *RelationService) GetMangaRelations(mangaID string) (*MangaRelationList, error) {
	return s.GetMangaRelationsContext(context.Background(), mangaID)
}

// GetMangaRelationsContext returns relations with context.
func (s *RelationService) GetMangaRelationsContext(ctx context.Context, mangaID string) (*MangaRelationList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaRelationPath, mangaID), nil)
	if err != nil {
		return nil, err
	}
	// Allow params via nil, but also support with params variant
	return s.getRelationsWithURL(ctx, urlStr)
}

func (s *RelationService) getRelationsWithURL(ctx context.Context, urlStr string) (*MangaRelationList, error) {
	var r MangaRelationList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMangaRelationsWithParams returns relations with query params.
func (s *RelationService) GetMangaRelationsWithParams(mangaID string, params url.Values) (*MangaRelationList, error) {
	return s.GetMangaRelationsWithParamsContext(context.Background(), mangaID, params)
}

// GetMangaRelationsWithParamsContext returns relations with params and context.
func (s *RelationService) GetMangaRelationsWithParamsContext(ctx context.Context, mangaID string, params url.Values) (*MangaRelationList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(MangaRelationPath, mangaID), params)
	if err != nil {
		return nil, err
	}
	var r MangaRelationList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
