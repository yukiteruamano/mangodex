package mangodex

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"
)

var client = NewDexClient()

func TestLogin(t *testing.T) {
	if os.Getenv("USERNAME") == "" || os.Getenv("PASSWORD") == "" {
		t.Skip("USERNAME/PASSWORD not set, skipping integration test")
	}
	err := client.Auth.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("%v\n", client)
}

func TestGetLoggedUser(t *testing.T) {
	if os.Getenv("USERNAME") == "" {
		t.Skip("skip without auth")
	}
	user, err := client.User.GetLoggedUser()
	if err != nil {
		t.Fatalf("Getting user failed: %v", err)
	}
	t.Log(user)
}

func TestGetMangaList(t *testing.T) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(5))
	params.Set("offset", strconv.Itoa(0))
	params.Set("includes[]", AuthorRel)
	_, err := client.Manga.GetMangaList(params)
	if err != nil {
		t.Fatalf("Getting manga failed: %v", err)
	}
}
