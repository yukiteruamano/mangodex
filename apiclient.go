package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ApiClientListPath   = "client"
	ApiClientPath       = "client/%s"
	ApiClientSecretPath = "client/%s/secret"
)

// ApiClientService provides ApiClient endpoints (read only).
type ApiClientService service

// ApiClientList is list of clients.
type ApiClientList = ListResponse[ApiClient]

// ApiClient holds client info.
type ApiClient struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes ApiClientAttributes `json:"attributes"`
}

type ApiClientAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ApiClientSecretResponse struct {
	Result string `json:"result"`
	Data   string `json:"data"`
}

func (r *ApiClientSecretResponse) GetResult() string { return r.Result }

// GetApiClients returns list.
func (s *ApiClientService) GetApiClients(params url.Values) (*ApiClientList, error) {
	return s.GetApiClientsContext(context.Background(), params)
}

// GetApiClientsContext returns with context.
func (s *ApiClientService) GetApiClientsContext(ctx context.Context, params url.Values) (*ApiClientList, error) {
	urlStr, err := s.client.buildURLWithParams(ApiClientListPath, params)
	if err != nil {
		return nil, err
	}
	var l ApiClientList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetApiClient returns single.
func (s *ApiClientService) GetApiClient(id string, params url.Values) (*SingleResponse[ApiClient], error) {
	return s.GetApiClientContext(context.Background(), id, params)
}

// GetApiClientContext returns with context.
func (s *ApiClientService) GetApiClientContext(ctx context.Context, id string, params url.Values) (*SingleResponse[ApiClient], error) {
	urlStr, err := s.client.buildURLWithParams(fmt.Sprintf(ApiClientPath, id), params)
	if err != nil {
		return nil, err
	}
	var r SingleResponse[ApiClient]
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetApiClientSecret returns secret.
func (s *ApiClientService) GetApiClientSecret(id string) (*ApiClientSecretResponse, error) {
	return s.GetApiClientSecretContext(context.Background(), id)
}

// GetApiClientSecretContext returns secret with context.
func (s *ApiClientService) GetApiClientSecretContext(ctx context.Context, id string) (*ApiClientSecretResponse, error) {
	u, err := s.client.buildURL(fmt.Sprintf(ApiClientSecretPath, id))
	if err != nil {
		return nil, err
	}
	var r ApiClientSecretResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
