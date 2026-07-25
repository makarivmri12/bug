package heuristic

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// DifferentialAnalyzer compares baseline and mutated responses
type DifferentialAnalyzer struct {
	logger *zap.Logger
}

// NewDifferentialAnalyzer creates a new differential analyzer
func NewDifferentialAnalyzer(logger *zap.Logger) *DifferentialAnalyzer {
	return &DifferentialAnalyzer{
		logger: logger,
	}
}

// Analyze compares baseline and mutated responses
func (da *DifferentialAnalyzer) Analyze(ctx context.Context, baseline, mutated *models.ResponseEvidence) *models.DeltaAnalysis {
	delta := &models.DeltaAnalysis{
		StatusChanged: baseline.Status != mutated.Status,
		TimeDiffMS:    mutated.TimeMS - baseline.TimeMS,
	}

	// Check if headers changed significantly
	delta.HeadersChanged = len(baseline.Headers) != len(mutated.Headers)

	// Calculate body diff
	if baseline.Body != mutated.Body {
		delta.BodyDiff = da.calculateDiff(baseline.Body, mutated.Body)
		delta.ContentLengthDiff = len(mutated.Body) - len(baseline.Body)
	}

	da.logger.Debug("differential analysis complete",
		zap.Bool("status_changed", delta.StatusChanged),
		zap.Int64("time_diff_ms", delta.TimeDiffMS),
		zap.Int("content_length_diff", delta.ContentLengthDiff),
	)

	return delta
}

// calculateDiff produces a simple diff representation
func (da *DifferentialAnalyzer) calculateDiff(before, after string) string {
	// Simple line-based diff
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")

	var diff strings.Builder

	for i := 0; i < len(beforeLines) && i < len(afterLines); i++ {
		if beforeLines[i] != afterLines[i] {
			diff.WriteString(fmt.Sprintf("Line %d:\n", i+1))
			diff.WriteString(fmt.Sprintf("- %s\n", beforeLines[i][:min(50, len(beforeLines[i]))]))
			diff.WriteString(fmt.Sprintf("+ %s\n", afterLines[i][:min(50, len(afterLines[i]))]))
		}
	}

	return diff.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
