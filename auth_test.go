package mangodex

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_Login_Success(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/auth/login", r.URL.Path)
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "user", body["username"])
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AuthResponse{Result: "ok", Token: token{Session: "sess123", Refresh: "ref456"}})
	})
	defer srv.Close()

	err := c.Auth.LoginContext(t.Context(), "user", "pass")
	require.NoError(t, err)
	assert.Equal(t, "ref456", c.Auth.GetRefreshToken())
	assert.True(t, c.Auth.IsLoggedIn())
	assert.Contains(t, c.header.Get("Authorization"), "sess123")
}

func TestAuth_Login_Failure(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Result: "error", Errors: []Error{{Title: "Unauthorized", Detail: "bad creds"}}})
	})
	defer srv.Close()
	err := c.Auth.LoginContext(t.Context(), "bad", "bad")
	require.Error(t, err)
	assert.False(t, c.Auth.IsLoggedIn())
}

func TestAuth_Logout(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/auth/logout", r.URL.Path)
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	// preset token
	c.Auth.SetRefreshToken("ref")
	c.header.Set("Authorization", "Bearer sess")
	err := c.Auth.LogoutContext(t.Context())
	require.NoError(t, err)
	assert.Empty(t, c.Auth.GetRefreshToken())
	assert.False(t, c.Auth.IsLoggedIn())
}

func TestAuth_Refresh(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/auth/refresh", r.URL.Path)
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "old-ref", body["token"])
		_ = json.NewEncoder(w).Encode(AuthResponse{Result: "ok", Token: token{Session: "new-sess", Refresh: "new-ref"}})
	})
	defer srv.Close()
	c.Auth.SetRefreshToken("old-ref")
	err := c.Auth.RefreshSessionTokenContext(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "new-ref", c.Auth.GetRefreshToken())
	assert.Contains(t, c.header.Get("Authorization"), "new-sess")
}

func TestAuth_Check(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/auth/check", r.URL.Path)
		_ = json.NewEncoder(w).Encode(Response{Result: "ok"})
	})
	defer srv.Close()
	c.header.Set("Authorization", "Bearer tok")
	resp, err := c.Auth.CheckContext(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Result)
}

func TestAuth_SetGetRefreshToken_Concurrent(t *testing.T) {
	c := NewDexClient()
	c.Auth.SetRefreshToken("abc")
	assert.Equal(t, "abc", c.Auth.GetRefreshToken())
}
