package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GroupListPath          = "group"
	GroupPath              = "group/%s"
	UserFollowsGroupPath   = "user/follows/group"
	UserFollowsGroupIdPath = "user/follows/group/%s"
)

// ScanlationGroupService provides scanlation group services.
type ScanlationGroupService service

// GroupList is a paginated list of groups.
type GroupList = ListResponse[ScanlationGroup]

// ScanlationGroup holds group information.
type ScanlationGroup struct {
	ID            string                    `json:"id"`
	Type          string                    `json:"type"`
	Attributes    ScanlationGroupAttributes `json:"attributes"`
	Relationships []Relationship            `json:"relationships"`
}

// ScanlationGroupAttributes holds group attributes.
type ScanlationGroupAttributes struct {
	Name            string             `json:"name"`
	AltNames        []LocalisedStrings `json:"altNames"`
	Website         *string            `json:"website"`
	IRCServer       *string            `json:"ircServer"`
	Discord         *string            `json:"discord"`
	ContactEmail    *string            `json:"contactEmail"`
	Description     *string            `json:"description"`
	Twitter         *string            `json:"twitter"`
	FocusedLanguage []string           `json:"focusedLanguage"`
	Locked          bool               `json:"locked"`
	Official        bool               `json:"official"`
	Inactive        bool               `json:"inactive"`
	PublishDelay    string             `json:"publishDelay"`
	Version         int                `json:"version"`
	CreatedAt       string             `json:"createdAt"`
	UpdatedAt       string             `json:"updatedAt"`
}

// SearchGroups searches scanlation groups.
func (s *ScanlationGroupService) SearchGroups(params url.Values) (*GroupList, error) {
	return s.SearchGroupsContext(context.Background(), params)
}

// SearchGroupsContext searches groups with context.
func (s *ScanlationGroupService) SearchGroupsContext(ctx context.Context, params url.Values) (*GroupList, error) {
	urlStr, err := s.client.buildURLWithParams(GroupListPath, params)
	if err != nil {
		return nil, err
	}
	var l GroupList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetGroup returns a single group by ID.
func (s *ScanlationGroupService) GetGroup(id string, params url.Values) (*SingleResponse[ScanlationGroup], error) {
	return s.GetGroupContext(context.Background(), id, params)
}

// GetGroupContext returns group with context.
func (s *ScanlationGroupService) GetGroupContext(ctx context.Context, id string, params url.Values) (*SingleResponse[ScanlationGroup], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(GroupPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[ScanlationGroup]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetFollowedGroups returns followed groups.
func (s *ScanlationGroupService) GetFollowedGroups(params url.Values) (*GroupList, error) {
	return s.GetFollowedGroupsContext(context.Background(), params)
}

// GetFollowedGroupsContext returns followed groups with context.
func (s *ScanlationGroupService) GetFollowedGroupsContext(ctx context.Context, params url.Values) (*GroupList, error) {
	urlStr, err := s.client.buildURLWithParams(UserFollowsGroupPath, params)
	if err != nil {
		return nil, err
	}
	var l GroupList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// CheckFollowedGroup checks if following a group.
func (s *ScanlationGroupService) CheckFollowedGroup(id string) (bool, error) {
	return s.CheckFollowedGroupContext(context.Background(), id)
}

// CheckFollowedGroupContext checks with context.
func (s *ScanlationGroupService) CheckFollowedGroupContext(ctx context.Context, id string) (bool, error) {
	u, err := s.client.buildURL(fmt.Sprintf(UserFollowsGroupIdPath, id))
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
