package mangodex

import (
	"context"
	"fmt"
	"net/http"
)

const (
	SettingsPath                = "settings"
	SettingsTemplatePath        = "settings/template"
	SettingsTemplateVersionPath = "settings/template/%s"
)

// SettingsService provides settings endpoints.
type SettingsService service

// SettingsResponse holds user settings.
type SettingsResponse struct {
	Result   string   `json:"result"`
	Settings Settings `json:"settings"`
}

func (r *SettingsResponse) GetResult() string { return r.Result }

type Settings struct {
	Version int `json:"version"`
}

// SettingsTemplateResponse holds template.
type SettingsTemplateResponse struct {
	Result   string           `json:"result"`
	Template SettingsTemplate `json:"template"`
}

func (r *SettingsTemplateResponse) GetResult() string { return r.Result }

type SettingsTemplate struct {
	Version int `json:"version"`
}

// GetSettings returns user settings.
func (s *SettingsService) GetSettings() (*SettingsResponse, error) {
	return s.GetSettingsContext(context.Background())
}

// GetSettingsContext returns settings with context.
func (s *SettingsService) GetSettingsContext(ctx context.Context) (*SettingsResponse, error) {
	u, err := s.client.buildURL(SettingsPath)
	if err != nil {
		return nil, err
	}
	var r SettingsResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetSettingsTemplate returns template.
func (s *SettingsService) GetSettingsTemplate() (*SettingsTemplateResponse, error) {
	return s.GetSettingsTemplateContext(context.Background())
}

// GetSettingsTemplateContext returns template with context.
func (s *SettingsService) GetSettingsTemplateContext(ctx context.Context) (*SettingsTemplateResponse, error) {
	u, err := s.client.buildURL(SettingsTemplatePath)
	if err != nil {
		return nil, err
	}
	var r SettingsTemplateResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetSettingsTemplateVersion returns template version.
func (s *SettingsService) GetSettingsTemplateVersion(version string) (*SettingsTemplateResponse, error) {
	return s.GetSettingsTemplateVersionContext(context.Background(), version)
}

// GetSettingsTemplateVersionContext returns template version with context.
func (s *SettingsService) GetSettingsTemplateVersionContext(ctx context.Context, version string) (*SettingsTemplateResponse, error) {
	u, err := s.client.buildURL(fmt.Sprintf(SettingsTemplateVersionPath, version))
	if err != nil {
		return nil, err
	}
	var r SettingsTemplateResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
