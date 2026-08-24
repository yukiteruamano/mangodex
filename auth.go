package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	LoginPath        = "auth/login"
	LogoutPath       = "auth/logout"
	RefreshTokenPath = "auth/refresh"
	CheckAuthPath    = "auth/check"
)

// AuthService provides authentication services.
type AuthService service

// AuthResponse is the typical auth response.
type AuthResponse struct {
	Result  string  `json:"result"`
	Token   token   `json:"token"`
	Message *string `json:"message,omitempty"`
}

func (ar AuthResponse) GetResult() string { return ar.Result }

type token struct {
	Session string `json:"session"`
	Refresh string `json:"refresh"`
}

// Login authenticates with username and password.
func (s *AuthService) Login(user, pwd string) error {
	return s.LoginContext(context.Background(), user, pwd)
}

// LoginContext is Login with custom context.
func (s *AuthService) LoginContext(ctx context.Context, user, pwd string) error {
	u, err := s.client.buildURL(LoginPath)
	if err != nil {
		return err
	}
	req := map[string]string{"username": user, "password": pwd}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("mangodex: marshal login: %w", err)
	}
	var ar AuthResponse
	if err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewReader(rBytes), &ar); err != nil {
		return err
	}
	s.client.mu.Lock()
	s.client.refreshToken = ar.Token.Refresh
	s.client.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	s.client.mu.Unlock()
	return nil
}

// Logout invalidates tokens.
func (s *AuthService) Logout() error { return s.LogoutContext(context.Background()) }

// LogoutContext is Logout with custom context.
func (s *AuthService) LogoutContext(ctx context.Context) error {
	u, err := s.client.buildURL(LogoutPath)
	if err != nil {
		return err
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), nil, &r); err != nil {
		return err
	}
	s.client.mu.Lock()
	s.client.refreshToken = ""
	s.client.header.Del("Authorization")
	s.client.mu.Unlock()
	return nil
}

// RefreshSessionToken refreshes the session using the stored refresh token.
func (s *AuthService) RefreshSessionToken() error {
	return s.RefreshSessionTokenContext(context.Background())
}

// RefreshSessionTokenContext refreshes with custom context.
func (s *AuthService) RefreshSessionTokenContext(ctx context.Context) error {
	s.client.mu.RLock()
	rt := s.client.refreshToken
	s.client.mu.RUnlock()

	u, err := s.client.buildURL(RefreshTokenPath)
	if err != nil {
		return err
	}
	req := map[string]string{"token": rt}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("mangodex: marshal refresh: %w", err)
	}
	var ar AuthResponse
	if err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewReader(rBytes), &ar); err != nil {
		return err
	}
	s.client.mu.Lock()
	s.client.refreshToken = ar.Token.Refresh
	s.client.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	s.client.mu.Unlock()
	return nil
}

// Check verifies if the current session is valid.
func (s *AuthService) Check() (*Response, error) {
	return s.CheckContext(context.Background())
}

// CheckContext verifies session with custom context.
func (s *AuthService) CheckContext(ctx context.Context) (*Response, error) {
	u, err := s.client.buildURL(CheckAuthPath)
	if err != nil {
		return nil, err
	}
	var r Response
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// IsLoggedIn returns true when Authorization header is set.
func (s *AuthService) IsLoggedIn() bool {
	s.client.mu.RLock()
	defer s.client.mu.RUnlock()
	return s.client.header.Get("Authorization") != ""
}

// GetRefreshToken returns the stored refresh token (thread-safe).
func (s *AuthService) GetRefreshToken() string {
	s.client.mu.RLock()
	defer s.client.mu.RUnlock()
	return s.client.refreshToken
}

// SetRefreshToken sets the refresh token (thread-safe).
func (s *AuthService) SetRefreshToken(refreshToken string) {
	s.client.mu.Lock()
	s.client.refreshToken = refreshToken
	s.client.mu.Unlock()
}
