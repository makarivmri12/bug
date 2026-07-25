package models

import (
	"time"

	"github.com/google/uuid"
)

// Severity levels for findings
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
)

// Category of vulnerability
type Category string

const (
	CategoryIDOR        Category = "IDOR"
	CategorySQLi        Category = "SQLi"
	CategoryXSS         Category = "XSS"
	CategorySSRF        Category = "SSRF"
	CategoryLogicFlaw   Category = "LogicFlaw"
	CategoryMisconfig   Category = "Misconfig"
	CategoryCVE         Category = "CVE"
	CategoryAuthBypass  Category = "AuthBypass"
	CategoryBrokenAccess Category = "BrokenAccess"
)

// Status of finding
type FindingStatus string

const (
	FindingStatusOpen           FindingStatus = "OPEN"
	FindingStatusConfirmed      FindingStatus = "CONFIRMED"
	FindingStatusFalsePositive  FindingStatus = "FALSE_POSITIVE"
	FindingStatusFixed          FindingStatus = "FIXED"
)

// RequestEvidence contains request details
type RequestEvidence struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Cookies map[string]string `json:"cookies"`
}

// ResponseEvidence contains response details
type ResponseEvidence struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	TimeMS  int64             `json:"time_ms"`
}

// DeltaAnalysis shows difference between baseline and mutated response
type DeltaAnalysis struct {
	StatusChanged   bool   `json:"status_changed"`
	BodyDiff        string `json:"body_diff"`
	TimeDiffMS      int64  `json:"time_diff_ms"`
	HeadersChanged  bool   `json:"headers_changed"`
	ContentLengthDiff int  `json:"content_length_diff"`
}

// Evidence contains all request/response evidence
type Evidence struct {
	Request           RequestEvidence  `json:"request"`
	Response          ResponseEvidence `json:"response"`
	BaselineResponse  ResponseEvidence `json:"baseline_response"`
	Delta             DeltaAnalysis    `json:"delta"`
}

// UniversalFinding represents a discovered vulnerability
type UniversalFinding struct {
	ID                 string      `json:"id"`
	ScanID             string      `json:"scan_id"`
	WorkspaceID        string      `json:"workspace_id"`
	Target             string      `json:"target"`
	Severity           Severity    `json:"severity"`
	Category           Category    `json:"category"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	Evidence           Evidence    `json:"evidence"`
	ReproductionSteps  []string    `json:"reproduction_steps"`
	Remediation        string      `json:"remediation"`
	PoCCurl            string      `json:"poc_curl"`
	Status             FindingStatus `json:"status"`
	AssignedTo         *string     `json:"assigned_to"`
	Comments           []string    `json:"comments"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

// NewUniversalFinding creates a new finding
func NewUniversalFinding(scanID, workspaceID, target string) *UniversalFinding {
	return &UniversalFinding{
		ID:          uuid.New().String(),
		ScanID:      scanID,
		WorkspaceID: workspaceID,
		Target:      target,
		Status:      FindingStatusOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
