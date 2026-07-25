package models

import "time"

// TargetType defines what kind of target to scan
type TargetType string

const (
	TargetTypeWeb     TargetType = "WEB"
	TargetTypeAPI     TargetType = "API"
	TargetTypeNetwork TargetType = "NETWORK"
)

// Target represents a scan target
type Target struct {
	ID           string      `json:"id"`
	WorkspaceID  string      `json:"workspace_id"`
	URL          string      `json:"url"`
	Type         TargetType  `json:"type"`
	Description  string      `json:"description"`
	AuthTokens   map[string]string `json:"auth_tokens"`
	Headers      map[string]string `json:"headers"`
	Cookies      map[string]string `json:"cookies"`
	ProxyURL     string      `json:"proxy_url"`
	Timeout      int         `json:"timeout"`
	RateLimit    int         `json:"rate_limit"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ScanRequest represents a scan job
type ScanRequest struct {
	ID           string        `json:"id"`
	WorkspaceID  string        `json:"workspace_id"`
	Targets      []Target      `json:"targets"`
	ScannerTypes []string      `json:"scanner_types"`
	Depth        int           `json:"depth"`
	Concurrency  int           `json:"concurrency"`
	Timeout      int           `json:"timeout"`
	CreatedAt    time.Time     `json:"created_at"`
}
