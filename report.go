package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ReportReasonsPath = "report/reasons/%s"
	ReportPath        = "report"
)

// ReportService provides report endpoints.
type ReportService service

// ReportReasonsList holds reasons for a category.
type ReportReasonsList struct {
	Result  string         `json:"result"`
	Reasons []ReportReason `json:"reasons"`
}

func (r *ReportReasonsList) GetResult() string { return r.Result }

// ReportReason holds a single reason.
type ReportReason struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Reason   string `json:"reason"`
	Details  string `json:"details"`
}

// ReportListResponse holds list of reports.
type ReportListResponse struct {
	Result   string   `json:"result"`
	Response string   `json:"response"`
	Data     []Report `json:"data"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Total    int      `json:"total"`
}

func (r *ReportListResponse) GetResult() string { return r.Result }

// Report holds report info.
type Report struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes ReportAttributes `json:"attributes"`
}

type ReportAttributes struct {
	Category  string `json:"category"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// GetReportReasons returns reasons for category.
func (s *ReportService) GetReportReasons(category string) (*ReportReasonsList, error) {
	return s.GetReportReasonsContext(context.Background(), category)
}

// GetReportReasonsContext returns reasons with context.
func (s *ReportService) GetReportReasonsContext(ctx context.Context, category string) (*ReportReasonsList, error) {
	u, err := s.client.buildURL(fmt.Sprintf(ReportReasonsPath, category))
	if err != nil {
		return nil, err
	}
	var r ReportReasonsList
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetReports returns list of reports.
func (s *ReportService) GetReports(params url.Values) (*ReportListResponse, error) {
	return s.GetReportsContext(context.Background(), params)
}

// GetReportsContext returns reports with context.
func (s *ReportService) GetReportsContext(ctx context.Context, params url.Values) (*ReportListResponse, error) {
	urlStr, err := s.client.buildURLWithParams(ReportPath, params)
	if err != nil {
		return nil, err
	}
	var r ReportListResponse
	if err := s.client.RequestAndDecode(ctx, http.MethodGet, urlStr, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
