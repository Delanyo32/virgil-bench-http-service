package worker

import (
	"log"
	"sync"
	"time"

	"github.com/example/ordersvc/pkg/queue"
)

// Job represents a unit of work to be processed.
type Job struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	CreatedAt time.Time
	Retries   int
}

// JobResult holds the result of processing a job.
type JobResult struct {
	JobID   string
	Success bool
	Error   string
}

// DispatcherStats tracks dispatcher metrics.
type DispatcherStats struct {
	Processed int64
	Pending   int64
	Failed    int64
}

// Dispatcher manages a pool of workers to process jobs from a queue.
type Dispatcher struct {
	queue      *queue.MemoryQueue
	maxWorkers int
	results    chan JobResult
	stats      DispatcherStats
	mu         sync.Mutex
	stopCh     chan struct{}
}

// NewDispatcher creates a new Dispatcher with a fixed worker pool size.
func NewDispatcher(q *queue.MemoryQueue, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		queue:      q,
		maxWorkers: maxWorkers,
		results:    make(chan JobResult, 100),
		stopCh:     make(chan struct{}),
	}
}

// Start begins the dispatcher loop, consuming jobs from the queue.
// FLAW: memory-leak-indicators -- spawns goroutines without context
// cancellation. If ProcessJobs is called repeatedly, goroutines accumulate
// because they are never cleaned up on Stop().
func (d *Dispatcher) Start() {
	log.Println("dispatcher started")

	// FLAW: unbounded goroutine spawning -- each poll iteration can spawn
	// goroutines that outlive the polling loop.
	for {
		select {
		case <-d.stopCh:
			log.Println("dispatcher stopping")
			return
		default:
			jobs := d.queue.Dequeue(d.maxWorkers)
			if len(jobs) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			d.ProcessJobs(jobs)
		}
	}
}

// ProcessJobs dispatches jobs to worker goroutines.
// FLAW: resource-exhaustion -- spawns one goroutine per job with no limit.
// A large batch of jobs creates unbounded goroutines that can exhaust memory.
func (d *Dispatcher) ProcessJobs(jobs []interface{}) {
	for _, rawJob := range jobs {
		job, ok := rawJob.(Job)
		if !ok {
			log.Printf("invalid job type: %T", rawJob)
			continue
		}

		// FLAW: goroutine spawned without context or WaitGroup.
		// These goroutines leak if the dispatcher is stopped while jobs
		// are still running.
		go func(j Job) {
			result := d.processOne(j)
			d.results <- result // blocks forever if results channel is full
		}(job)
	}
}

// processOne handles a single job.
func (d *Dispatcher) processOne(job Job) JobResult {
	log.Printf("processing job %s (type: %s)", job.ID, job.Type)

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	d.mu.Lock()
	d.stats.Processed++
	d.mu.Unlock()

	return JobResult{
		JobID:   job.ID,
		Success: true,
	}
}

// Stop signals the dispatcher to stop processing.
func (d *Dispatcher) Stop() {
	close(d.stopCh)
	// FLAW: does not wait for in-flight goroutines to complete.
	// Goroutines spawned by ProcessJobs continue running after Stop returns.
}

// Stats returns the current dispatcher statistics.
func (d *Dispatcher) Stats() DispatcherStats {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.stats
}

// collectResults drains the results channel.
// Dead code -- never called from Start or any external caller.
func (d *Dispatcher) collectResults() {
	for result := range d.results {
		d.mu.Lock()
		if !result.Success {
			d.stats.Failed++
		}
		d.mu.Unlock()
	}
}
