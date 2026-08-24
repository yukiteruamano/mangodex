package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	CustomListPath        = "list/%s"
	CustomListListPath    = "list"
	CustomListFeedPath    = "list/%s/feed"
	UserCustomListPath    = "user/list"
	UserIdCustomListPath  = "user/%s/list"
	UserFollowsListPath   = "user/follows/list"
	UserFollowsListIdPath = "user/follows/list/%s"
)

// CustomListService provides custom list services.
type CustomListService service

// CustomListList is a paginated list of custom lists.
type CustomListList = ListResponse[CustomList]

// CustomList holds custom list information.
type CustomList struct {
	ID            string               `json:"id"`
	Type          string               `json:"type"`
	Attributes    CustomListAttributes `json:"attributes"`
	Relationships []Relationship       `json:"relationships"`
}

// CustomListAttributes holds custom list attributes.
type CustomListAttributes struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
	Version    int    `json:"version"`
}

// SearchCustomLists searches custom lists.
func (s *CustomListService) SearchCustomLists(params url.Values) (*CustomListList, error) {
	return s.SearchCustomListsContext(context.Background(), params)
}

// SearchCustomListsContext searches with context.
func (s *CustomListService) SearchCustomListsContext(ctx context.Context, params url.Values) (*CustomListList, error) {
	urlStr, err := s.client.buildURLWithParams(CustomListListPath, params)
	if err != nil {
		return nil, err
	}
	var l CustomListList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetCustomList returns a single custom list.
func (s *CustomListService) GetCustomList(id string) (*SingleResponse[CustomList], error) {
	return s.GetCustomListContext(context.Background(), id)
}

// GetCustomListContext returns custom list with context.
func (s *CustomListService) GetCustomListContext(ctx context.Context, id string) (*SingleResponse[CustomList], error) {
	u, err := s.client.buildURL(fmt.Sprintf(CustomListPath, id))
	if err != nil {
		return nil, err
	}
	var r SingleResponse[CustomList]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetCustomListFeed returns feed for a custom list.
func (s *CustomListService) GetCustomListFeed(id string, params url.Values) (*ChapterList, error) {
	return s.GetCustomListFeedContext(context.Background(), id, params)
}

// GetCustomListFeedContext returns feed with context.
func (s *CustomListService) GetCustomListFeedContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(CustomListFeedPath, id), params)
	if err != nil {
		return nil, err
	}
	var l ChapterList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLoggedUserCustomLists returns custom lists for logged user.
func (s *CustomListService) GetLoggedUserCustomLists(params url.Values) (*CustomListList, error) {
	return s.GetLoggedUserCustomListsContext(context.Background(), params)
}

// GetLoggedUserCustomListsContext returns with context.
func (s *CustomListService) GetLoggedUserCustomListsContext(ctx context.Context, params url.Values) (*CustomListList, error) {
	urlStr, err := s.client.buildURLWithParams(UserCustomListPath, params)
	if err != nil {
		return nil, err
	}
	var l CustomListList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetUserCustomLists returns custom lists for a user.
func (s *CustomListService) GetUserCustomLists(userID string, params url.Values) (*CustomListList, error) {
	return s.GetUserCustomListsContext(context.Background(), userID, params)
}

// GetUserCustomListsContext returns with context.
func (s *CustomListService) GetUserCustomListsContext(ctx context.Context, userID string, params url.Values) (*CustomListList, error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(UserIdCustomListPath, userID), params)
	if err != nil {
		return nil, err
	}
	var l CustomListList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetFollowedCustomLists returns followed custom lists.
func (s *CustomListService) GetFollowedCustomLists(params url.Values) (*CustomListList, error) {
	return s.GetFollowedCustomListsContext(context.Background(), params)
}

// GetFollowedCustomListsContext returns with context.
func (s *CustomListService) GetFollowedCustomListsContext(ctx context.Context, params url.Values) (*CustomListList, error) {
	urlStr, err := s.client.buildURLWithParams(UserFollowsListPath, params)
	if err != nil {
		return nil, err
	}
	var l CustomListList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// CheckFollowedCustomList checks if following a custom list.
func (s *CustomListService) CheckFollowedCustomList(id string) (bool, error) {
	return s.CheckFollowedCustomListContext(context.Background(), id)
}

// CheckFollowedCustomListContext checks with context.
func (s *CustomListService) CheckFollowedCustomListContext(ctx context.Context, id string) (bool, error) {
	u, err := s.client.buildURL(fmt.Sprintf(UserFollowsListIdPath, id))
	if err != nil {
		return false, err
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
