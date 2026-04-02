package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/example/ordersvc/pkg/queue"
)

// ScheduledTask represents a recurring job configuration.
type ScheduledTask struct {
	Name     string
	Interval time.Duration
	JobType  string
	Payload  map[string]interface{}
}

// Scheduler manages periodic task execution.
type Scheduler struct {
	queue *queue.MemoryQueue
	tasks []ScheduledTask
}

// NewScheduler creates a new Scheduler.
func NewScheduler(q *queue.MemoryQueue) *Scheduler {
	s := &Scheduler{
		queue: q,
	}

	// Register default scheduled tasks
	s.tasks = []ScheduledTask{
		{
			Name:     "cleanup_expired_orders",
			Interval: 5 * time.Minute,
			JobType:  "cleanup",
			Payload:  map[string]interface{}{"max_age_hours": 24},
		},
		{
			Name:     "sync_inventory",
			Interval: 10 * time.Minute,
			JobType:  "inventory_sync",
			Payload:  map[string]interface{}{"source": "warehouse_api"},
		},
		{
			Name:     "generate_daily_report",
			Interval: 24 * time.Hour,
			JobType:  "generate_report",
			Payload:  map[string]interface{}{"type": "daily_summary"},
		},
	}

	return s
}

// RunForever starts all scheduled tasks.
// FLAW: memory-leak-indicators -- spawns goroutines per task with no
// WaitGroup or context to track or cancel them. On shutdown, these
// goroutines leak indefinitely.
func (s *Scheduler) RunForever() {
	for _, task := range s.tasks {
		// FLAW: goroutine launched with no tracking mechanism
		go s.runTask(task)
	}

	log.Println("scheduler started with", len(s.tasks), "tasks")

	// Block forever -- no graceful shutdown support
	select {}
}

// runTask executes a single scheduled task on its interval.
func (s *Scheduler) runTask(task ScheduledTask) {
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	log.Printf("scheduled task '%s' every %v", task.Name, task.Interval)

	for range ticker.C {
		job := Job{
			ID:        fmt.Sprintf("%s-%d", task.Name, time.Now().UnixNano()),
			Type:      task.JobType,
			Payload:   task.Payload,
			CreatedAt: time.Now(),
		}

		if err := s.queue.Enqueue(job); err != nil {
			log.Printf("failed to enqueue task '%s': %v", task.Name, err)
		}
	}
}

// AddTask adds a new scheduled task. Not goroutine-safe --
// modifying the tasks slice while RunForever iterates it.
func (s *Scheduler) AddTask(task ScheduledTask) {
	s.tasks = append(s.tasks, task)
}

// RemoveTask removes a scheduled task by name.
// Dead code -- never called from any module.
func (s *Scheduler) RemoveTask(name string) {
	for i, t := range s.tasks {
		if t.Name == name {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return
		}
	}
}

// ListTasks returns the names of all scheduled tasks.
func (s *Scheduler) ListTasks() []string {
	names := make([]string, len(s.tasks))
	for i, t := range s.tasks {
		names[i] = t.Name
	}
	return names
}
