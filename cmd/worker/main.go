package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/example/ordersvc/internal/config"
	"github.com/example/ordersvc/internal/worker"
	"github.com/example/ordersvc/pkg/queue"
)

func main() {
	cfg := config.Load()

	memQueue := queue.NewMemoryQueue(100)
	dispatcher := worker.NewDispatcher(memQueue, cfg.WorkerCount)
	scheduler := worker.NewScheduler(memQueue)

	// FLAW: goroutine leak -- starts goroutines without context cancellation.
	// If main exits abnormally, these goroutines are never cleaned up.
	go dispatcher.Start()
	go scheduler.RunForever()

	// FLAW: no graceful shutdown -- runs forever with no signal handling.
	// goroutines from dispatcher and scheduler leak if process is killed.
	log.Printf("worker started with %d goroutines", cfg.WorkerCount)

	// Busy-wait health check with no context
	for {
		stats := dispatcher.Stats()
		log.Printf("worker stats: processed=%d, pending=%d, failed=%d",
			stats.Processed, stats.Pending, stats.Failed)
		time.Sleep(30 * time.Second) // FLAW: magic number for check interval
	}

	// FLAW: dead code -- unreachable cleanup
	fmt.Println("shutting down worker")
	dispatcher.Stop()
	_ = os
}
