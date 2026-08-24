package mangodex

import (
	"context"
	"net/http"
	"net/url"
)

const RatingPath = "rating"

// RatingService provides rating endpoints.
type RatingService service

// RatingList holds ratings for manga (GET /rating?manga[]=id)
type RatingList struct {
	Result  string            `json:"result"`
	Ratings map[string]Rating `json:"ratings"`
}

func (r *RatingList) GetResult() string { return r.Result }

// Rating holds a single rating.
type Rating struct {
	Rating    int    `json:"rating"`
	CreatedAt string `json:"createdAt"`
}

// GetRatings returns ratings for manga ids.
func (s *RatingService) GetRatings(params url.Values) (*RatingList, error) {
	return s.GetRatingsContext(context.Background(), params)
}

// GetRatingsContext returns ratings with context.
func (s *RatingService) GetRatingsContext(ctx context.Context, params url.Values) (*RatingList, error) {
	urlStr, err := s.client.buildURLWithParams(RatingPath, params)
	if err != nil {
		return nil, err
	}
	var r RatingList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
