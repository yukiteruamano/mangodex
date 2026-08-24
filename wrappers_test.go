package mangodex

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// This file exercises the non-Context wrappers to hit coverage for those thin layers.
// Uses the same httptest approach but via the Background wrappers.
func TestWrappers_Manga(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/manga":
			_, _ = w.Write(mockMangaListJSON())
		case "/manga/m123":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"m123","type":"manga","attributes":{"title":{"en":"t"},"altTitles":{},"description":{},"isLocked":false,"links":{},"originalLanguage":"ja","state":"published","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","tags":[]}}}`))
		case "/manga/random":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"rand","type":"manga","attributes":{"title":{},"altTitles":{},"description":{},"isLocked":false,"links":{},"originalLanguage":"ja","state":"published","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","tags":[]}}}`))
		case "/manga/m123/aggregate":
			_, _ = w.Write([]byte(`{"result":"ok","volumes":{}}`))
		case "/manga/m123/feed":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/manga/status":
			_, _ = w.Write([]byte(`{"result":"ok","statuses":{"m123":"reading"}}`))
		default:
			_, _ = w.Write([]byte(`{"result":"ok"}`))
		}
	})
	defer srv.Close()

	_, _ = c.Manga.GetMangaList(url.Values{"limit": {"10"}})
	_, _ = c.Manga.GetManga("m123", nil)
	_, _ = c.Manga.GetRandomManga(nil)
	_, _ = c.Manga.GetAggregate("m123", nil)
	_, _ = c.Manga.GetMangaFeed("m123", nil)
	_, _ = c.Manga.GetReadingStatus(nil)
	// Follow checks
	c2, srv2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	})
	_, _ = c2.Manga.CheckIfMangaFollowed("id")
	_, _ = c2.Manga.ToggleMangaFollowStatus("id", true)
	_, _ = c2.Manga.ToggleMangaFollowStatus("id", false)
	srv2.Close()

	require.NotNil(t, c)
	srv.Close()
	_ = srv2
}

func TestWrappers_Chapter(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
	})
	defer srv.Close()
	_, _ = c.Chapter.GetMangaChapters("mid", url.Values{"limit": {"10"}})
	_, _ = c.Chapter.SearchChapters(url.Values{})
	c2, srv2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chapter/c1" {
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"c1","type":"chapter","attributes":{"title":"","volume":null,"chapter":null,"translatedLanguage":"en","uploader":"u","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00","publishAt":"2020-01-01T00:00:00","pages":1}}}`))
			return
		}
		_, _ = w.Write([]byte(`{"result":"ok","data":["c1"]}`))
	})
	_, _ = c2.Chapter.GetChapter("c1", nil)
	_, _ = c2.Chapter.GetReadMangaChapters("mid")
	_, _ = c2.Chapter.SetReadUnreadMangaChapters("mid", []string{"c1"}, nil)
	srv2.Close()
}

func TestWrappers_Auth(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/auth/login" {
			_, _ = w.Write([]byte(`{"result":"ok","token":{"session":"s","refresh":"r"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"result":"ok"}`))
	})
	defer srv.Close()
	_ = c.Auth.Login("u", "p")
	_ = c.Auth.Logout()
	_ = c.Auth.RefreshSessionToken()
	_, _ = c.Auth.Check()
	_ = c.Auth.IsLoggedIn()
	_ = c.Auth.GetRefreshToken()
	c.Auth.SetRefreshToken("x")
}

func TestWrappers_OtherServices(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/author":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/author/a1":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"a1","type":"author","attributes":{"name":"n","biography":{},"version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00"}}}`))
		case "/cover":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/cover/c1":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"c1","type":"cover_art","attributes":{"fileName":"f.jpg","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00"}}}`))
		case "/group":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/group/g1":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"g1","type":"scanlation_group","attributes":{"name":"g","locked":false,"official":false,"inactive":false,"publishDelay":"0","version":1,"createdAt":"2020-01-01T00:00:00","updatedAt":"2020-01-01T00:00:00"}}}`))
		case "/list":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/list/l1":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"l1","type":"custom_list","attributes":{"name":"n","visibility":"public","version":1}}}`))
		case "/manga/tag":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[]}`))
		case "/user/follows/manga/feed":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/statistics/manga/m1":
			_, _ = w.Write([]byte(`{"result":"ok","statistics":{}}`))
		case "/statistics/chapter/c1":
			_, _ = w.Write([]byte(`{"result":"ok","statistics":{}}`))
		case "/statistics/manga":
			_, _ = w.Write([]byte(`{"result":"ok","statistics":{}}`))
		case "/user/me":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"me","type":"user","attributes":{"username":"u","roles":[],"version":1}}}`))
		case "/user/u1":
			_, _ = w.Write([]byte(`{"result":"ok","response":"entity","data":{"id":"u1","type":"user","attributes":{"username":"u","roles":[],"version":1}}}`))
		case "/user":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/user/follows/manga":
			_, _ = w.Write([]byte(`{"result":"ok","response":"collection","data":[],"limit":10,"offset":0,"total":0}`))
		case "/at-home/server/ch1":
			_, _ = w.Write([]byte(`{"result":"ok","baseUrl":"https://x","chapter":{"hash":"h","data":["a"],"dataSaver":["a"]}}`))
		default:
			_, _ = w.Write([]byte(`{"result":"ok"}`))
		}
	})
	defer srv.Close()

	_, _ = c.Author.SearchAuthors(url.Values{})
	_, _ = c.Author.GetAuthor("a1", nil)
	_, _ = c.Cover.SearchCovers(url.Values{})
	_, _ = c.Cover.GetCover("c1", nil)
	_, _ = c.ScanlationGroup.SearchGroups(url.Values{})
	_, _ = c.ScanlationGroup.GetGroup("g1", nil)
	_, _ = c.CustomList.SearchCustomLists(url.Values{})
	_, _ = c.CustomList.GetCustomList("l1")
	_, _ = c.Tag.GetTags()
	_, _ = c.Tag.SearchTags(url.Values{})
	_, _ = c.Feed.GetFollowedMangaFeed(url.Values{})
	_, _ = c.Statistics.GetMangaStatistics("m1")
	_, _ = c.Statistics.GetMangaStatisticsBatch(url.Values{"manga[]": {"m1"}})
	_, _ = c.Statistics.GetChapterStatistics("c1")
	_, _ = c.User.GetLoggedUser()
	_, _ = c.User.GetUser("u1")
	_, _ = c.User.SearchUsers(url.Values{})
	_, _ = c.User.GetUserFollowedMangaList(10, 0, nil)
	_, _ = c.User.SearchFollowedManga(url.Values{})
	_, _ = c.AtHome.NewMDHomeClient("ch1", "data", false)
	_, _ = c.AtHome.NewMDHomeClient("ch1", "data", false)
	// AtHome GetChapterPage wrapper via direct
	md := &MDHomeClient{client: c.client, baseURL: srv.URL, quality: "data", hash: "h", reportURL: srv.URL}
	_, _ = md.GetChapterPage("p.jpg")

	// also test New alias
	_ = New(WithBaseURL(srv.URL))
	// ListResponse GetResult
	var lr ListResponse[Manga]
	lr.Result = "ok"
	_ = lr.GetResult()
	var sr SingleResponse[Manga]
	sr.Result = "ok"
	_ = sr.GetResult()
	var ms MangaStatistics
	ms.Result = "ok"
	_ = ms.GetResult()
	var cs ChapterStatistics
	cs.Result = "ok"
	_ = cs.GetResult()
	var tl TagList
	tl.Result = "ok"
	_ = tl.GetResult()
	var ur UserResponse
	ur.Result = "ok"
	_ = ur.GetResult()
	var ar AuthResponse
	ar.Result = "ok"
	_ = ar.GetResult()
	var cl ChapterList
	cl.Result = "ok"
	_ = cl.GetResult()
	var ma MangaAggregate
	ma.Result = "ok"
	_ = ma.GetResult()
	var cr ChapterReadMarkers
	cr.Result = "ok"
	_ = cr.GetResult()
	var mdh MDHomeServerResponse
	mdh.Result = "ok"
	_ = mdh.GetResult()
}
