package scanner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// WebScanner performs web-based vulnerability scanning
type WebScanner struct {
	name   string
	client *http.Client
	logger *zap.Logger
}

// NewWebScanner creates a new web scanner instance
func NewWebScanner(logger *zap.Logger) *WebScanner {
	return &WebScanner{
		name: "web-scanner",
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		logger: logger,
	}
}

// Name returns the scanner name
func (ws *WebScanner) Name() string {
	return ws.name
}

// Category returns the scanner category
func (ws *WebScanner) Category() string {
	return "web"
}

// CanHandle determines if this scanner can handle the target
func (ws *WebScanner) CanHandle(ctx context.Context, target models.Target) bool {
	return target.Type == models.TargetTypeWeb || target.Type == models.TargetTypeAPI
}

// Scan performs the actual scanning
func (ws *WebScanner) Scan(ctx context.Context, target models.Target) ([]models.UniversalFinding, error) {
	var findings []models.UniversalFinding

	// Probe the target
	baselineResp, err := ws.probeTarget(ctx, target)
	if err != nil {
		ws.logger.Error("failed to probe target", zap.Error(err))
		return nil, err
	}

	// Test for common vulnerabilities
	finding := ws.testForCommonVulns(ctx, target, baselineResp)
	if finding != nil {
		findings = append(findings, *finding)
	}

	return findings, nil
}

// probeTarget makes a baseline request to the target
func (ws *WebScanner) probeTarget(ctx context.Context, target models.Target) (*models.ResponseEvidence, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
	if err != nil {
		return nil, err
	}

	// Add custom headers
	for k, v := range target.Headers {
		req.Header.Add(k, v)
	}

	// Add cookies
	for k, v := range target.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	start := time.Now()
	resp, err := ws.client.Do(req)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &models.ResponseEvidence{
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(body),
		TimeMS:  duration.Milliseconds(),
	}, nil
}

// testForCommonVulns tests for common vulnerabilities
func (ws *WebScanner) testForCommonVulns(ctx context.Context, target models.Target, baseline *models.ResponseEvidence) *models.UniversalFinding {
	// Example: Test for missing security headers
	if baseline.Headers["X-Frame-Options"] == "" ||
		baseline.Headers["Content-Security-Policy"] == "" ||
		baseline.Headers["X-Content-Type-Options"] == "" {

		finding := models.NewUniversalFinding("", target.WorkspaceID, target.URL)
		finding.Title = "Missing Security Headers"
		finding.Description = "The target is missing important HTTP security headers"
		finding.Category = models.CategoryMisconfig
		finding.Severity = models.SeverityMedium
		finding.Evidence = models.Evidence{
			BaselineResponse: *baseline,
		}
		finding.Remediation = "Add proper security headers to HTTP responses"

		return finding
	}

	return nil
}
