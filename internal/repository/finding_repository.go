package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// FindingRepository handles finding persistence
type FindingRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewFindingRepository creates a new finding repository
func NewFindingRepository(db *sql.DB, logger *zap.Logger) *FindingRepository {
	return &FindingRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new finding into the database
func (fr *FindingRepository) Create(ctx context.Context, finding *models.UniversalFinding) error {
	query := `
		INSERT INTO findings (
			id, scan_id, workspace_id, target, severity, category, title,
			description, evidence, reproduction_steps, remediation, poc_curl,
			status, assigned_to, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := fr.db.ExecContext(ctx, query,
		finding.ID,
		finding.ScanID,
		finding.WorkspaceID,
		finding.Target,
		finding.Severity,
		finding.Category,
		finding.Title,
		finding.Description,
		"{}" /* TODO: serialize evidence to JSON */,
		pq.Array(finding.ReproductionSteps),
		finding.Remediation,
		finding.PoCCurl,
		finding.Status,
		finding.AssignedTo,
		finding.CreatedAt,
		finding.UpdatedAt,
	)

	if err != nil {
		fr.logger.Error("failed to create finding", zap.Error(err))
		return err
	}

	return nil
}

// GetByID retrieves a finding by ID
func (fr *FindingRepository) GetByID(ctx context.Context, id string) (*models.UniversalFinding, error) {
	query := `
		SELECT id, scan_id, workspace_id, target, severity, category, title,
		       description, reproduction_steps, remediation, poc_curl, status,
		       assigned_to, created_at, updated_at
		FROM findings WHERE id = $1
	`

	var finding models.UniversalFinding
	var reproSteps pq.StringArray
	var assignedTo sql.NullString

	err := fr.db.QueryRowContext(ctx, query, id).Scan(
		&finding.ID,
		&finding.ScanID,
		&finding.WorkspaceID,
		&finding.Target,
		&finding.Severity,
		&finding.Category,
		&finding.Title,
		&finding.Description,
		&reproSteps,
		&finding.Remediation,
		&finding.PoCCurl,
		&finding.Status,
		&assignedTo,
		&finding.CreatedAt,
		&finding.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fr.logger.Error("failed to get finding", zap.Error(err))
		return nil, err
	}

	finding.ReproductionSteps = []string(reproSteps)
	if assignedTo.Valid {
		finding.AssignedTo = &assignedTo.String
	}

	return &finding, nil
}

// GetByScanID retrieves all findings for a scan
func (fr *FindingRepository) GetByScanID(ctx context.Context, scanID string) ([]models.UniversalFinding, error) {
	query := `
		SELECT id, scan_id, workspace_id, target, severity, category, title,
		       description, reproduction_steps, remediation, poc_curl, status,
		       assigned_to, created_at, updated_at
		FROM findings WHERE scan_id = $1 ORDER BY created_at DESC
	`

	rows, err := fr.db.QueryContext(ctx, query, scanID)
	if err != nil {
		fr.logger.Error("failed to query findings", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var findings []models.UniversalFinding
	for rows.Next() {
		var finding models.UniversalFinding
		var reproSteps pq.StringArray
		var assignedTo sql.NullString

		err := rows.Scan(
			&finding.ID,
			&finding.ScanID,
			&finding.WorkspaceID,
			&finding.Target,
			&finding.Severity,
			&finding.Category,
			&finding.Title,
			&finding.Description,
			&reproSteps,
			&finding.Remediation,
			&finding.PoCCurl,
			&finding.Status,
			&assignedTo,
			&finding.CreatedAt,
			&finding.UpdatedAt,
		)

		if err != nil {
			fr.logger.Error("failed to scan finding", zap.Error(err))
			return nil, err
		}

		finding.ReproductionSteps = []string(reproSteps)
		if assignedTo.Valid {
			finding.AssignedTo = &assignedTo.String
		}

		findings = append(findings, finding)
	}

	return findings, rows.Err()
}

// Update updates an existing finding
func (fr *FindingRepository) Update(ctx context.Context, finding *models.UniversalFinding) error {
	query := `
		UPDATE findings SET
			status = $1, assigned_to = $2, updated_at = $3
		WHERE id = $4
	`

	finding.UpdatedAt = time.Now()
	_, err := fr.db.ExecContext(ctx, query,
		finding.Status,
		finding.AssignedTo,
		finding.UpdatedAt,
		finding.ID,
	)

	if err != nil {
		fr.logger.Error("failed to update finding", zap.Error(err))
		return err
	}

	return nil
}

// Delete deletes a finding
func (fr *FindingRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM findings WHERE id = $1"

	_, err := fr.db.ExecContext(ctx, query, id)
	if err != nil {
		fr.logger.Error("failed to delete finding", zap.Error(err))
		return err
	}

	return nil
}
