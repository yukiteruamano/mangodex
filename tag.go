package mangodex

import (
	"context"
	"net/http"
	"net/url"
)

const TagListPath = "manga/tag"

// TagService provides tag services.
type TagService service

// TagList is a list of tags.
type TagList struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     []Tag  `json:"data"`
}

func (l *TagList) GetResult() string { return l.Result }

// GetTags returns all tags.
func (s *TagService) GetTags() (*TagList, error) {
	return s.GetTagsContext(context.Background())
}

// GetTagsContext returns tags with context.
func (s *TagService) GetTagsContext(ctx context.Context) (*TagList, error) {
	u, err := s.client.buildURL(TagListPath)
	if err != nil {
		return nil, err
	}
	var l TagList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// SearchTags with params (includes grouping etc if needed).
func (s *TagService) SearchTags(params url.Values) (*TagList, error) {
	return s.SearchTagsContext(context.Background(), params)
}

func (s *TagService) SearchTagsContext(ctx context.Context, params url.Values) (*TagList, error) {
	urlStr, err := s.client.buildURLWithParams(TagListPath, params)
	if err != nil {
		return nil, err
	}
	var l TagList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}
