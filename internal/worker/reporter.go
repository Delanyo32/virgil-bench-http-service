package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

// MetricSnapshot holds a point-in-time set of metrics.
type MetricSnapshot struct {
	Timestamp     time.Time
	RequestCount  int64
	ErrorCount    int64
	AvgLatencyMs  float64
}

// Reporter periodically collects and flushes application metrics.
type Reporter struct {
	interval     time.Duration
	mu           sync.Mutex
	requestCount int64
	errorCount   int64
	latencySum   float64
	latencyN     int64
	snapshots    []MetricSnapshot
}

// NewReporter creates a Reporter that flushes metrics at the given interval.
func NewReporter(interval time.Duration) *Reporter {
	return &Reporter{
		interval:  interval,
		snapshots: make([]MetricSnapshot, 0),
	}
}

// RecordRequest records a completed request with its latency.
func (r *Reporter) RecordRequest(latencyMs float64, isError bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requestCount++
	r.latencySum += latencyMs
	r.latencyN++
	if isError {
		r.errorCount++
	}
}

// Run starts the reporting loop, which respects context cancellation.
func (r *Reporter) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.flush()
			log.Println("reporter shutting down")
			return ctx.Err()
		case <-ticker.C:
			r.flush()
		}
	}
}

// flush takes a snapshot of current metrics and resets counters.
func (r *Reporter) flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var avgLatency float64
	if r.latencyN > 0 {
		avgLatency = r.latencySum / float64(r.latencyN)
	}

	snap := MetricSnapshot{
		Timestamp:    time.Now(),
		RequestCount: r.requestCount,
		ErrorCount:   r.errorCount,
		AvgLatencyMs: avgLatency,
	}

	r.snapshots = append(r.snapshots, snap)
	log.Printf("metrics: requests=%d errors=%d avg_latency=%.2fms",
		r.requestCount, r.errorCount, avgLatency)

	r.requestCount = 0
	r.errorCount = 0
	r.latencySum = 0
	r.latencyN = 0
}

// Snapshots returns all recorded metric snapshots.
func (r *Reporter) Snapshots() []MetricSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]MetricSnapshot, len(r.snapshots))
	copy(result, r.snapshots)
	return result
}
