package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/pkg/models"
)

// ScannerPlugin defines the interface for scanner implementations
type ScannerPlugin interface {
	Name() string
	Category() string // "web" | "network" | "logic"
	CanHandle(ctx context.Context, target models.Target) bool
	Scan(ctx context.Context, target models.Target) ([]models.UniversalFinding, error)
}

// Task represents a scanning task in the queue
type Task struct {
	ID        string
	ScanID    string
	Target    models.Target
	Scanners  []ScannerPlugin
	Retries   int
	CreatedAt int64
}

// WorkerPool manages concurrent scanning workers
type WorkerPool struct {
	workers    int
	workQueue  chan *Task
	resultChan chan *models.UniversalFinding
	doneChan   chan struct{}
	wg         sync.WaitGroup
	mu         sync.Mutex
	logger     *zap.Logger
	active     int
}

// NewWorkerPool creates a new worker pool with specified concurrency
func NewWorkerPool(workers int, logger *zap.Logger) *WorkerPool {
	if workers <= 0 {
		workers = 50 // default
	}

	return &WorkerPool{
		workers:    workers,
		workQueue:  make(chan *Task, workers*2),
		resultChan: make(chan *models.UniversalFinding, workers*2),
		doneChan:   make(chan struct{}),
		logger:     logger,
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
	wp.logger.Info("worker pool started", zap.Int("workers", wp.workers))
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.workQueue)
	wp.wg.Wait()
	close(wp.resultChan)
	wp.logger.Info("worker pool stopped")
}

// Submit adds a task to the work queue
func (wp *WorkerPool) Submit(task *Task) error {
	wp.mu.Lock()
	wp.active++
	wp.mu.Unlock()

	select {
	case wp.workQueue <- task:
		return nil
	default:
		wp.mu.Lock()
		wp.active--
		wp.mu.Unlock()
		return fmt.Errorf("work queue full")
	}
}

// ResultChannel returns the result channel
func (wp *WorkerPool) ResultChannel() <-chan *models.UniversalFinding {
	return wp.resultChan
}

// ActiveTasks returns number of active tasks
func (wp *WorkerPool) ActiveTasks() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.active
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			wp.logger.Info("worker shutting down", zap.Int("worker_id", id))
			return
		case task, ok := <-wp.workQueue:
			if !ok {
				wp.logger.Info("work queue closed", zap.Int("worker_id", id))
				return
			}

			wp.processTask(ctx, task)

			wp.mu.Lock()
			wp.active--
			wp.mu.Unlock()
		}
	}
}

// processTask executes scanners against a target
func (wp *WorkerPool) processTask(ctx context.Context, task *Task) {
	wp.logger.Info("processing task",
		zap.String("task_id", task.ID),
		zap.String("target", task.Target.URL),
	)

	for _, scanner := range task.Scanners {
		if !scanner.CanHandle(ctx, task.Target) {
			wp.logger.Debug("scanner cannot handle target",
				zap.String("scanner", scanner.Name()),
				zap.String("target", task.Target.URL),
			)
			continue
		}

		wp.logger.Info("running scanner",
			zap.String("scanner", scanner.Name()),
			zap.String("target", task.Target.URL),
		)

		findings, err := scanner.Scan(ctx, task.Target)
		if err != nil {
			wp.logger.Error("scanner error",
				zap.String("scanner", scanner.Name()),
				zap.Error(err),
			)
			continue
		}

		for _, finding := range findings {
			select {
			case wp.resultChan <- &finding:
			case <-ctx.Done():
				return
			}
		}
	}
}
