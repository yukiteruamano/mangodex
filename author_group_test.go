package mangodex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthor_SearchAndGet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/author" {
			_ = json.NewEncoder(w).Encode(ListResponse[Author]{Result: "ok", Data: []Author{{ID: "a1"}}})
			return
		}
		assert.Equal(t, "/author/a1", r.URL.Path)
		_ = json.NewEncoder(w).Encode(SingleResponse[Author]{Result: "ok", Data: Author{ID: "a1"}})
	})
	defer srv.Close()
	list, err := c.Author.SearchAuthorsContext(t.Context(), url.Values{"name": {"test"}})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
	single, err := c.Author.GetAuthorContext(t.Context(), "a1", nil)
	require.NoError(t, err)
	assert.Equal(t, "a1", single.Data.ID)
}

func TestGroup_SearchAndGet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/group" {
			_ = json.NewEncoder(w).Encode(ListResponse[ScanlationGroup]{Result: "ok", Data: []ScanlationGroup{{ID: "g1"}}})
			return
		}
		assert.Equal(t, "/group/g1", r.URL.Path)
		_ = json.NewEncoder(w).Encode(SingleResponse[ScanlationGroup]{Result: "ok", Data: ScanlationGroup{ID: "g1"}})
	})
	defer srv.Close()
	list, err := c.ScanlationGroup.SearchGroupsContext(t.Context(), url.Values{"name": {"test"}})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
	single, err := c.ScanlationGroup.GetGroupContext(t.Context(), "g1", nil)
	require.NoError(t, err)
	assert.Equal(t, "g1", single.Data.ID)
}

func TestCustomList_SearchAndGet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/list" {
			_ = json.NewEncoder(w).Encode(ListResponse[CustomList]{Result: "ok", Data: []CustomList{{ID: "l1"}}})
			return
		}
		_ = json.NewEncoder(w).Encode(SingleResponse[CustomList]{Result: "ok", Data: CustomList{ID: "l1"}})
	})
	defer srv.Close()
	list, err := c.CustomList.SearchCustomListsContext(t.Context(), url.Values{})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
	single, err := c.CustomList.GetCustomListContext(t.Context(), "l1")
	require.NoError(t, err)
	assert.Equal(t, "l1", single.Data.ID)
}

func TestTag_GetTags(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/manga/tag", r.URL.Path)
		_ = json.NewEncoder(w).Encode(TagList{Result: "ok", Data: []Tag{{ID: "t1"}}})
	})
	defer srv.Close()
	list, err := c.Tag.GetTagsContext(t.Context())
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
	_, err = c.Tag.SearchTagsContext(t.Context(), url.Values{})
	require.NoError(t, err)
}

func TestFeed_GetFollowedFeed(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user/follows/manga/feed", r.URL.Path)
		_ = json.NewEncoder(w).Encode(ListResponse[Chapter]{Result: "ok", Data: []Chapter{{ID: "ch1"}}})
	})
	defer srv.Close()
	list, err := c.Feed.GetFollowedMangaFeedContext(t.Context(), url.Values{"limit": {"10"}})
	require.NoError(t, err)
	assert.Len(t, list.Data, 1)
}

func TestStatistics_MangaAndChapter(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/statistics/manga/mid" {
			_ = json.NewEncoder(w).Encode(MangaStatistics{Result: "ok", Statistics: map[string]MangaStatistic{"mid": {Follows: intPtr(100)}}}) //nolint:modernize
			return
		}
		if r.URL.Path == "/statistics/manga" {
			_ = json.NewEncoder(w).Encode(MangaStatistics{Result: "ok", Statistics: map[string]MangaStatistic{}})
			return
		}
		_ = json.NewEncoder(w).Encode(ChapterStatistics{Result: "ok", Statistics: map[string]ChapterStatistic{"cid": {}}})
	})
	defer srv.Close()
	stats, err := c.Statistics.GetMangaStatisticsContext(t.Context(), "mid")
	require.NoError(t, err)
	assert.Equal(t, "ok", stats.Result)
	_, err = c.Statistics.GetMangaStatisticsBatchContext(t.Context(), url.Values{"manga[]": {"mid"}})
	require.NoError(t, err)
	cs, err := c.Statistics.GetChapterStatisticsContext(t.Context(), "cid")
	require.NoError(t, err)
	assert.Equal(t, "ok", cs.Result)
}

func TestUser_GetUserAndSearch(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/me" {
			_ = json.NewEncoder(w).Encode(UserResponse{Result: "ok", Data: User{ID: "me"}})
			return
		}
		if r.URL.Path == "/user" {
			_ = json.NewEncoder(w).Encode(ListResponse[User]{Result: "ok", Data: []User{{ID: "u1"}}})
			return
		}
		if r.URL.Path == "/user/follows/manga" {
			_ = json.NewEncoder(w).Encode(ListResponse[Manga]{Result: "ok", Response: "collection", Data: []Manga{{ID: "m1"}}})
			return
		}
		_ = json.NewEncoder(w).Encode(SingleResponse[User]{Result: "ok", Data: User{ID: "u1"}})
	})
	defer srv.Close()
	_, err := c.User.GetLoggedUserContext(t.Context())
	require.NoError(t, err)
	_, err = c.User.GetUserContext(t.Context(), "u1")
	require.NoError(t, err)
	_, err = c.User.SearchUsersContext(t.Context(), url.Values{"username": {"test"}})
	require.NoError(t, err)
	_, err = c.User.GetUserFollowedMangaListContext(t.Context(), 10, 0, []string{"manga"})
	require.NoError(t, err)
	_, err = c.User.SearchFollowedMangaContext(t.Context(), url.Values{"limit": {"10"}})
	require.NoError(t, err)
}
