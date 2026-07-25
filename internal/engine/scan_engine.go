package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// ScanEngine orchestrates the scanning workflow
type ScanEngine struct {
	workerPool *WorkerPool
	redisClient *redis.Client
	logger *zap.Logger
	plugins map[string]ScannerPlugin
}

// NewScanEngine creates a new scan engine
func NewScanEngine(workers int, redisClient *redis.Client, logger *zap.Logger) *ScanEngine {
	return &ScanEngine{
		workerPool: NewWorkerPool(workers, logger),
		redisClient: redisClient,
		logger: logger,
		plugins: make(map[string]ScannerPlugin),
	}
}

// RegisterScanner registers a scanner plugin
func (se *ScanEngine) RegisterScanner(scanner ScannerPlugin) {
	se.plugins[scanner.Name()] = scanner
	se.logger.Info("scanner registered",
		zap.String("name", scanner.Name()),
		zap.String("category", scanner.Category()),
	)
}

// Start starts the scan engine and worker pool
func (se *ScanEngine) Start(ctx context.Context) error {
	se.workerPool.Start(ctx)
	se.logger.Info("scan engine started")
	return nil
}

// Stop gracefully stops the scan engine
func (se *ScanEngine) Stop() error {
	se.workerPool.Stop()
	se.logger.Info("scan engine stopped")
	return nil
}

// ExecuteScan performs a complete scan of targets
func (se *ScanEngine) ExecuteScan(ctx context.Context, req models.ScanRequest) (string, error) {
	scanID := uuid.New().String()

	se.logger.Info("starting scan",
		zap.String("scan_id", scanID),
		zap.Int("targets", len(req.Targets)),
	)

	// Filter scanners based on requested types
	selectedScanners := se.selectScanners(req.ScannerTypes)
	if len(selectedScanners) == 0 {
		return "", fmt.Errorf("no scanners available for types: %v", req.ScannerTypes)
	}

	// Submit tasks for each target
	for _, target := range req.Targets {
		task := &Task{
			ID:       uuid.New().String(),
			ScanID:   scanID,
			Target:   target,
			Scanners: selectedScanners,
			Retries:  0,
			CreatedAt: time.Now().Unix(),
		}

		if err := se.workerPool.Submit(task); err != nil {
			se.logger.Error("failed to submit task", zap.Error(err))
			return "", err
		}
	}

	// Store scan metadata in Redis
	metadata := map[string]interface{}{
		"scan_id": scanID,
		"workspace_id": req.WorkspaceID,
		"target_count": len(req.Targets),
		"status": "running",
		"started_at": time.Now().Unix(),
	}

	if err := se.storeScanMetadata(ctx, scanID, metadata); err != nil {
		se.logger.Error("failed to store scan metadata", zap.Error(err))
	}

	return scanID, nil
}

// selectScanners returns scanners for requested types
func (se *ScanEngine) selectScanners(scannerTypes []string) []ScannerPlugin {
	var selected []ScannerPlugin

	if len(scannerTypes) == 0 {
		// Use all scanners if none specified
		for _, scanner := range se.plugins {
			selected = append(selected, scanner)
		}
		return selected
	}

	for _, scanType := range scannerTypes {
		if scanner, ok := se.plugins[scanType]; ok {
			selected = append(selected, scanner)
		}
	}

	return selected
}

// storeScanMetadata stores scan metadata in Redis
func (se *ScanEngine) storeScanMetadata(ctx context.Context, scanID string, metadata map[string]interface{}) error {
	key := fmt.Sprintf("scan:%s", scanID)
	// Simple JSON encoding would be used in production
	return se.redisClient.Set(ctx, key, fmt.Sprintf("%v", metadata), 24*time.Hour).Err()
}

// GetScanStatus returns the status of a scan
func (se *ScanEngine) GetScanStatus(ctx context.Context, scanID string) (string, error) {
	key := fmt.Sprintf("scan:%s", scanID)
	val, err := se.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
