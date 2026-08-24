package mangodex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadExtra_Infrastructure(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ping", r.URL.Path)
		_, _ = w.Write([]byte("pong"))
	})
	defer srv.Close()
	s, err := c.Infrastructure.GetPingContext(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "pong", s)
	_, _ = c.Infrastructure.GetPing()
}

func TestReadExtra_Rating(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rating", r.URL.Path)
		_ = json.NewEncoder(w).Encode(RatingList{Result: "ok", Ratings: map[string]Rating{"m1": {Rating: 8}}})
	})
	defer srv.Close()
	_, err := c.Rating.GetRatingsContext(t.Context(), url.Values{"manga[]": {"m1"}})
	require.NoError(t, err)
	_, _ = c.Rating.GetRatings(url.Values{})
}

func TestReadExtra_Report(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/report/reasons/manga" {
			_ = json.NewEncoder(w).Encode(ReportReasonsList{Result: "ok"})
			return
		}
		_ = json.NewEncoder(w).Encode(ReportListResponse{Result: "ok", Data: []Report{}})
	})
	defer srv.Close()
	_, err := c.Report.GetReportReasonsContext(t.Context(), "manga")
	require.NoError(t, err)
	_, _ = c.Report.GetReportReasons("manga")
	_, err = c.Report.GetReportsContext(t.Context(), url.Values{})
	require.NoError(t, err)
	_, _ = c.Report.GetReports(nil)
}

func TestReadExtra_Relation(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/relation")
		_ = json.NewEncoder(w).Encode(MangaRelationList{Result: "ok"})
	})
	defer srv.Close()
	_, err := c.Relation.GetMangaRelationsContext(t.Context(), "mid")
	require.NoError(t, err)
	_, _ = c.Relation.GetMangaRelations("mid")
	_, err = c.Relation.GetMangaRelationsWithParamsContext(t.Context(), "mid", url.Values{})
	require.NoError(t, err)
	_, _ = c.Relation.GetMangaRelationsWithParams("mid", nil)
}

func TestReadExtra_Settings(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/settings/template/1" {
			_ = json.NewEncoder(w).Encode(SettingsTemplateResponse{Result: "ok"})
			return
		}
		if r.URL.Path == "/settings/template" {
			_ = json.NewEncoder(w).Encode(SettingsTemplateResponse{Result: "ok"})
			return
		}
		_ = json.NewEncoder(w).Encode(SettingsResponse{Result: "ok"})
	})
	defer srv.Close()
	_, err := c.Settings.GetSettingsContext(t.Context())
	require.NoError(t, err)
	_, _ = c.Settings.GetSettings()
	_, err = c.Settings.GetSettingsTemplateContext(t.Context())
	require.NoError(t, err)
	_, _ = c.Settings.GetSettingsTemplate()
	_, err = c.Settings.GetSettingsTemplateVersionContext(t.Context(), "1")
	require.NoError(t, err)
	_, _ = c.Settings.GetSettingsTemplateVersion("1")
}

func TestReadExtra_Upload(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(UploadSessionResponse{Result: "ok"})
	})
	defer srv.Close()
	_, err := c.Upload.GetUploadSessionContext(t.Context())
	require.NoError(t, err)
	_, _ = c.Upload.GetUploadSession()
}

func TestReadExtra_Manga(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/manga/mid/recommendation" {
			_ = json.NewEncoder(w).Encode(MangaRecommendationList{Result: "ok", Data: []Manga{}})
			return
		}
		if r.URL.Path == "/manga/draft" {
			_ = json.NewEncoder(w).Encode(MangaDraftList{Result: "ok"})
			return
		}
		if r.URL.Path == "/manga/draft/d1" {
			_ = json.NewEncoder(w).Encode(SingleResponse[Manga]{Result: "ok"})
			return
		}
		if r.URL.Path == "/manga/mid/status" {
			_ = json.NewEncoder(w).Encode(mangaIDStatusResponse{Result: "ok", Status: "reading"})
			return
		}
		_ = json.NewEncoder(w).Encode(MangaRecommendationList{Result: "ok"})
	})
	defer srv.Close()
	_, err := c.Manga.GetMangaRecommendationsContext(t.Context(), "mid", nil)
	require.NoError(t, err)
	_, _ = c.Manga.GetMangaRecommendations("mid", nil)
	_, err = c.Manga.GetMangaDraftsContext(t.Context(), nil)
	require.NoError(t, err)
	_, _ = c.Manga.GetMangaDrafts(nil)
	_, err = c.Manga.GetMangaDraftContext(t.Context(), "d1", nil)
	require.NoError(t, err)
	_, _ = c.Manga.GetMangaDraft("d1", nil)
	_, err = c.Manga.GetMangaIDStatusContext(t.Context(), "mid")
	require.NoError(t, err)
	_, _ = c.Manga.GetMangaIDStatus("mid")
}

func TestReadExtra_CustomList(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/list/l1/feed":
			_ = json.NewEncoder(w).Encode(ChapterList{Result: "ok"})
		case "/user/list":
			_ = json.NewEncoder(w).Encode(CustomListList{Result: "ok"})
		case "/user/uid/list":
			_ = json.NewEncoder(w).Encode(CustomListList{Result: "ok"})
		case "/user/follows/list":
			_ = json.NewEncoder(w).Encode(CustomListList{Result: "ok"})
		case "/user/follows/list/l1":
			_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
		default:
			_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
		}
	})
	defer srv.Close()
	_, _ = c.CustomList.GetCustomListFeedContext(t.Context(), "l1", nil)
	_, _ = c.CustomList.GetCustomListFeed("l1", nil)
	_, _ = c.CustomList.GetLoggedUserCustomListsContext(t.Context(), nil)
	_, _ = c.CustomList.GetLoggedUserCustomLists(nil)
	_, _ = c.CustomList.GetUserCustomListsContext(t.Context(), "uid", nil)
	_, _ = c.CustomList.GetUserCustomLists("uid", nil)
	_, _ = c.CustomList.GetFollowedCustomListsContext(t.Context(), nil)
	_, _ = c.CustomList.GetFollowedCustomLists(nil)
	_, _ = c.CustomList.CheckFollowedCustomListContext(t.Context(), "l1")
	_, _ = c.CustomList.CheckFollowedCustomList("l1")
	// 404 case
	c2, srv2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Not Found", Detail: "404"}}})
	})
	defer srv2.Close()
	ok, _ := c2.CustomList.CheckFollowedCustomListContext(t.Context(), "l1")
	assert.False(t, ok)
}

func TestReadExtra_ScanlationGroupFollows(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/follows/group" {
			_ = json.NewEncoder(w).Encode(GroupList{Result: "ok"})
			return
		}
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	_, _ = c.ScanlationGroup.GetFollowedGroupsContext(t.Context(), nil)
	_, _ = c.ScanlationGroup.GetFollowedGroups(nil)
	_, _ = c.ScanlationGroup.CheckFollowedGroupContext(t.Context(), "g1")
	_, _ = c.ScanlationGroup.CheckFollowedGroup("g1")
	c2, srv2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Not Found", Detail: "404 not found"}}})
	})
	defer srv2.Close()
	ok, _ := c2.ScanlationGroup.CheckFollowedGroupContext(t.Context(), "g1")
	assert.False(t, ok)
}

func TestReadExtra_UserFollows(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/follows/user" {
			_ = json.NewEncoder(w).Encode(UserList{Result: "ok"})
			return
		}
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	_, _ = c.User.GetFollowedUsersContext(t.Context(), nil)
	_, _ = c.User.GetFollowedUsers(nil)
	_, _ = c.User.CheckFollowedUserContext(t.Context(), "u1")
	_, _ = c.User.CheckFollowedUser("u1")
	c2, srv2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "404", Detail: "not found"}}})
	})
	defer srv2.Close()
	ok, _ := c2.User.CheckFollowedUserContext(t.Context(), "u1")
	assert.False(t, ok)
}

func TestReadExtra_Statistics(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/statistics/chapter":
			_ = json.NewEncoder(w).Encode(ChapterStatistics{Result: "ok"})
		case "/statistics/group/g1":
			_ = json.NewEncoder(w).Encode(GroupStatistics{Result: "ok"})
		case "/statistics/group":
			_ = json.NewEncoder(w).Encode(GroupStatistics{Result: "ok"})
		default:
			_ = json.NewEncoder(w).Encode(GroupStatistics{Result: "ok"})
		}
	})
	defer srv.Close()
	_, _ = c.Statistics.GetChapterStatisticsBatchContext(t.Context(), url.Values{"chapter[]": {"c1"}})
	_, _ = c.Statistics.GetChapterStatisticsBatch(nil)
	_, _ = c.Statistics.GetGroupStatisticsContext(t.Context(), "g1")
	_, _ = c.Statistics.GetGroupStatistics("g1")
	_, _ = c.Statistics.GetGroupStatisticsBatchContext(t.Context(), url.Values{"group[]": {"g1"}})
	_, _ = c.Statistics.GetGroupStatisticsBatch(nil)
}

func TestReadExtra_Feed(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user/history" {
			_ = json.NewEncoder(w).Encode(ChapterList{Result: "ok"})
			return
		}
		_ = json.NewEncoder(w).Encode(MangaReadResponse{Result: "ok", Data: map[string][]string{"m1": {"c1"}}})
	})
	defer srv.Close()
	_, _ = c.Feed.GetReadingHistoryContext(t.Context(), nil)
	_, _ = c.Feed.GetReadingHistory(nil)
	_, _ = c.Feed.GetMangaReadHistoryContext(t.Context(), url.Values{"ids[]": {"m1"}})
	_, _ = c.Feed.GetMangaReadHistory(nil)
}

func TestReadExtra_ApiClient(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/client":
			_ = json.NewEncoder(w).Encode(ApiClientList{Result: "ok"})
		case "/client/c1":
			_ = json.NewEncoder(w).Encode(SingleResponse[ApiClient]{Result: "ok"})
		case "/client/c1/secret":
			_ = json.NewEncoder(w).Encode(ApiClientSecretResponse{Result: "ok", Data: "secret"})
		}
	})
	defer srv.Close()
	_, _ = c.ApiClient.GetApiClientsContext(t.Context(), nil)
	_, _ = c.ApiClient.GetApiClients(nil)
	_, _ = c.ApiClient.GetApiClientContext(t.Context(), "c1", nil)
	_, _ = c.ApiClient.GetApiClient("c1", nil)
	_, _ = c.ApiClient.GetApiClientSecretContext(t.Context(), "c1")
	_, _ = c.ApiClient.GetApiClientSecret("c1")
}
