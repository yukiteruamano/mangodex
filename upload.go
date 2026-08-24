package mangodex

import (
	"context"
	"net/http"
)

const (
	UploadSessionPath = "upload"
)

// UploadService provides upload session endpoints (read only).
type UploadService service

// UploadSession holds session info.
type UploadSession struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes UploadSessionAttributes `json:"attributes"`
}

type UploadSessionAttributes struct {
	IsCommitted bool   `json:"isCommitted"`
	IsProcessed bool   `json:"isProcessed"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type UploadSessionResponse struct {
	Result   string        `json:"result"`
	Response string        `json:"response"`
	Data     UploadSession `json:"data"`
}

func (r *UploadSessionResponse) GetResult() string { return r.Result }

// GetUploadSession returns current upload session.
func (s *UploadService) GetUploadSession() (*UploadSessionResponse, error) {
	return s.GetUploadSessionContext(context.Background())
}

// GetUploadSessionContext returns session with context.
func (s *UploadService) GetUploadSessionContext(ctx context.Context) (*UploadSessionResponse, error) {
	u, err := s.client.buildURL(UploadSessionPath)
	if err != nil {
		return nil, err
	}
	var r UploadSessionResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
