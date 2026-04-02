package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

// CleanupFunc is a function that performs a cleanup operation.
// It returns the number of items cleaned up and any error encountered.
type CleanupFunc func(ctx context.Context, olderThan time.Time) (int, error)

// CleanupWorker periodically removes expired data using a provided function.
type CleanupWorker struct {
	interval   time.Duration
	retention  time.Duration
	cleanupFn  CleanupFunc
	mu         sync.Mutex
	lastRun    time.Time
	totalItems int
}

// NewCleanupWorker creates a worker that runs at the given interval,
// cleaning up data older than the retention period.
func NewCleanupWorker(interval, retention time.Duration, fn CleanupFunc) *CleanupWorker {
	return &CleanupWorker{
		interval:  interval,
		retention: retention,
		cleanupFn: fn,
	}
}

// Run starts the cleanup loop, which respects context cancellation.
func (w *CleanupWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("cleanup worker shutting down")
			return ctx.Err()
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

// runOnce executes a single cleanup pass.
func (w *CleanupWorker) runOnce(ctx context.Context) {
	cutoff := time.Now().Add(-w.retention)
	count, err := w.cleanupFn(ctx, cutoff)
	if err != nil {
		log.Printf("cleanup error: %v", err)
		return
	}

	w.mu.Lock()
	w.lastRun = time.Now()
	w.totalItems += count
	w.mu.Unlock()

	if count > 0 {
		log.Printf("cleanup: removed %d expired items", count)
	}
}

// Stats returns the time of the last run and total items cleaned.
func (w *CleanupWorker) Stats() (time.Time, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastRun, w.totalItems
}
