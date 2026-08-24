package mangodex

import (
	"context"
	"io"
	"net/http"
)

const PingPath = "ping"

// InfrastructureService provides infrastructure endpoints.
type InfrastructureService service

// GetPing returns pong if API is healthy.
func (s *InfrastructureService) GetPing() (string, error) {
	return s.GetPingContext(context.Background())
}

// GetPingContext returns ping with context. Response is text/plain.
func (s *InfrastructureService) GetPingContext(ctx context.Context) (string, error) {
	u, err := s.client.buildURL(PingPath)
	if err != nil {
		return "", err
	}
	resp, err := s.client.Request(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
