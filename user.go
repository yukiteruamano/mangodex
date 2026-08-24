package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetUserFollowedMangaListPath = "user/follows/manga"
	GetLoggedUserPath            = "user/me"
	UserPath                     = "user/%s"
	UserListPath                 = "user"
	UserFollowsUserPath          = "user/follows/user"
	UserFollowsUserIdPath        = "user/follows/user/%s"
)

// UserService provides user services.
type UserService service

// GetUserFollowedMangaList returns followed manga.
func (s *UserService) GetUserFollowedMangaList(limit, offset int, includes []string) (*MangaList, error) {
	return s.GetUserFollowedMangaListContext(context.Background(), limit, offset, includes)
}

// GetUserFollowedMangaListContext returns followed manga with context.
func (s *UserService) GetUserFollowedMangaListContext(ctx context.Context, limit, offset int, includes []string) (*MangaList, error) {
	u, err := s.client.buildURL(GetUserFollowedMangaListPath)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	for _, inc := range includes {
		q.Add("includes[]", inc)
	}
	u.RawQuery = q.Encode()
	var l MangaList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// SearchFollowedManga is a generic search over followed manga.
func (s *UserService) SearchFollowedManga(params url.Values) (*MangaList, error) {
	return s.SearchFollowedMangaContext(context.Background(), params)
}

func (s *UserService) SearchFollowedMangaContext(ctx context.Context, params url.Values) (*MangaList, error) {
	urlStr, err := s.client.buildURLWithParams(GetUserFollowedMangaListPath, params)
	if err != nil {
		return nil, err
	}
	var l MangaList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// UserResponse is a single user response.
type UserResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     User   `json:"data"`
}

func (r *UserResponse) GetResult() string { return r.Result }

// UserList is a paginated list of users.
type UserList = ListResponse[User]

// User holds user info.
type User struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    UserAttributes `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

// UserAttributes holds user attributes.
type UserAttributes struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Version  int      `json:"version"`
}

// GetLoggedUser returns the logged user.
func (s *UserService) GetLoggedUser() (*UserResponse, error) {
	return s.GetLoggedUserContext(context.Background())
}

// GetLoggedUserContext returns logged user with context.
func (s *UserService) GetLoggedUserContext(ctx context.Context) (*UserResponse, error) {
	u, err := s.client.buildURL(GetLoggedUserPath)
	if err != nil {
		return nil, err
	}
	var r UserResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(id string) (*SingleResponse[User], error) {
	return s.GetUserContext(context.Background(), id)
}

func (s *UserService) GetUserContext(ctx context.Context, id string) (*SingleResponse[User], error) {
	u, err := s.client.buildURL(fmt.Sprintf(UserPath, id))
	if err != nil {
		return nil, err
	}
	var r SingleResponse[User]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// SearchUsers searches users.
func (s *UserService) SearchUsers(params url.Values) (*UserList, error) {
	return s.SearchUsersContext(context.Background(), params)
}

func (s *UserService) SearchUsersContext(ctx context.Context, params url.Values) (*UserList, error) {
	urlStr, err := s.client.buildURLWithParams(UserListPath, params)
	if err != nil {
		return nil, err
	}
	var l UserList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetFollowedUsers returns followed users.
func (s *UserService) GetFollowedUsers(params url.Values) (*UserList, error) {
	return s.GetFollowedUsersContext(context.Background(), params)
}

// GetFollowedUsersContext returns followed users with context.
func (s *UserService) GetFollowedUsersContext(ctx context.Context, params url.Values) (*UserList, error) {
	urlStr, err := s.client.buildURLWithParams(UserFollowsUserPath, params)
	if err != nil {
		return nil, err
	}
	var l UserList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// CheckFollowedUser checks if following a user.
func (s *UserService) CheckFollowedUser(id string) (bool, error) {
	return s.CheckFollowedUserContext(context.Background(), id)
}

// CheckFollowedUserContext checks with context.
func (s *UserService) CheckFollowedUserContext(ctx context.Context, id string) (bool, error) {
	u, err := s.client.buildURL(fmt.Sprintf(UserFollowsUserIdPath, id))
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
