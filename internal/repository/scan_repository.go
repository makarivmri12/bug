package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// ScanRepository handles scan persistence
type ScanRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewScanRepository creates a new scan repository
func NewScanRepository(db *sql.DB, logger *zap.Logger) *ScanRepository {
	return &ScanRepository{
		db:     db,
		logger: logger,
	}
}

// Scan represents a scan record
type Scan struct {
	ID          string
	WorkspaceID string
	Status      string
	Progress    int
	FindingsCount int
	StartedAt   time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Create inserts a new scan
func (sr *ScanRepository) Create(ctx context.Context, workspaceID string) (string, error) {
	scanID := uuid.New().String()

	query := `
		INSERT INTO scans (id, workspace_id, status, progress, findings_count, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := sr.db.ExecContext(ctx, query,
		scanID,
		workspaceID,
		"running",
		0,
		0,
		now,
		now,
		now,
	)

	if err != nil {
		sr.logger.Error("failed to create scan", zap.Error(err))
		return "", err
	}

	return scanID, nil
}

// GetByID retrieves a scan by ID
func (sr *ScanRepository) GetByID(ctx context.Context, id string) (*Scan, error) {
	query := `
		SELECT id, workspace_id, status, progress, findings_count, started_at, completed_at, created_at, updated_at
		FROM scans WHERE id = $1
	`

	var scan Scan
	var completedAt sql.NullTime

	err := sr.db.QueryRowContext(ctx, query, id).Scan(
		&scan.ID,
		&scan.WorkspaceID,
		&scan.Status,
		&scan.Progress,
		&scan.FindingsCount,
		&scan.StartedAt,
		&completedAt,
		&scan.CreatedAt,
		&scan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		sr.logger.Error("failed to get scan", zap.Error(err))
		return nil, err
	}

	if completedAt.Valid {
		scan.CompletedAt = &completedAt.Time
	}

	return &scan, nil
}

// UpdateStatus updates scan status
func (sr *ScanRepository) UpdateStatus(ctx context.Context, scanID, status string) error {
	query := `UPDATE scans SET status = $1, updated_at = $2 WHERE id = $3`

	_, err := sr.db.ExecContext(ctx, query, status, time.Now(), scanID)
	if err != nil {
		sr.logger.Error("failed to update scan status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateProgress updates scan progress
func (sr *ScanRepository) UpdateProgress(ctx context.Context, scanID string, progress int) error {
	query := `UPDATE scans SET progress = $1, updated_at = $2 WHERE id = $3`

	_, err := sr.db.ExecContext(ctx, query, progress, time.Now(), scanID)
	if err != nil {
		sr.logger.Error("failed to update scan progress", zap.Error(err))
		return err
	}

	return nil
}

// Complete marks a scan as completed
func (sr *ScanRepository) Complete(ctx context.Context, scanID string, findingsCount int) error {
	query := `
		UPDATE scans
		SET status = $1, findings_count = $2, completed_at = $3, progress = 100, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	_, err := sr.db.ExecContext(ctx, query, "completed", findingsCount, now, now, scanID)
	if err != nil {
		sr.logger.Error("failed to complete scan", zap.Error(err))
		return err
	}

	return nil
}
