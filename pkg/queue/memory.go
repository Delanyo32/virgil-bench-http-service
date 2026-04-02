package queue

import (
	"fmt"
	"log"
	"sync"
)

// MemoryQueue implements Queue using an in-memory channel.
type MemoryQueue struct {
	// FLAW: unbuffered channel causes deadlock when Enqueue is called from
	// the same goroutine that calls Dequeue, or when the channel is full.
	// The buffer size passed to NewMemoryQueue is used, but if the queue
	// fills up, Enqueue blocks forever.
	ch       chan interface{}
	mu       sync.Mutex
	enqueued int64
	dequeued int64
}

// NewMemoryQueue creates a new MemoryQueue with the given capacity.
func NewMemoryQueue(capacity int) *MemoryQueue {
	return &MemoryQueue{
		ch: make(chan interface{}, capacity),
	}
}

// Enqueue adds an item to the queue.
// FLAW: blocks indefinitely if the channel is full. No timeout or
// context cancellation. In production, this can cause goroutine pileup.
func (q *MemoryQueue) Enqueue(item interface{}) error {
	q.ch <- item // blocks if channel is full -- potential deadlock
	q.mu.Lock()
	q.enqueued++
	q.mu.Unlock()
	return nil
}

// Dequeue removes up to count items from the queue.
func (q *MemoryQueue) Dequeue(count int) []interface{} {
	var items []interface{}
	for i := 0; i < count; i++ {
		select {
		case item := <-q.ch:
			items = append(items, item)
			q.mu.Lock()
			q.dequeued++
			q.mu.Unlock()
		default:
			return items
		}
	}
	return items
}

// Size returns the current number of items in the queue.
func (q *MemoryQueue) Size() int {
	return len(q.ch)
}

// Stats returns queue statistics.
func (q *MemoryQueue) Stats() (int64, int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.enqueued, q.dequeued
}

// Drain removes all items from the queue and returns them.
// Dead code -- not called from any module.
func (q *MemoryQueue) Drain() []interface{} {
	var items []interface{}
	for {
		select {
		case item := <-q.ch:
			items = append(items, item)
		default:
			return items
		}
	}
}

// EnqueueBatch adds multiple items to the queue.
func (q *MemoryQueue) EnqueueBatch(items []interface{}) error {
	for i, item := range items {
		if err := q.Enqueue(item); err != nil {
			return fmt.Errorf("failed to enqueue item %d: %w", i, err)
		}
	}
	return nil
}

// unused logging helper
func logQueueEvent(event string, count int) {
	log.Printf("queue event: %s (count: %d)", event, count)
}
