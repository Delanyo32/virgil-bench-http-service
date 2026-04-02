package queue

import (
	"container/heap"
	"sync"
)

// PriorityItem represents an item in the priority queue.
type PriorityItem struct {
	Value    interface{}
	Priority int
	index    int
}

// priorityHeap implements heap.Interface for priority queue operations.
type priorityHeap []*PriorityItem

func (h priorityHeap) Len() int           { return len(h) }
func (h priorityHeap) Less(i, j int) bool { return h[i].Priority > h[j].Priority }
func (h priorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Push adds an item to the heap.
func (h *priorityHeap) Push(x interface{}) {
	item := x.(*PriorityItem)
	item.index = len(*h)
	*h = append(*h, item)
}

// Pop removes the highest-priority item from the heap.
func (h *priorityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[:n-1]
	return item
}

// PriorityQueue is a thread-safe priority queue backed by container/heap.
type PriorityQueue struct {
	mu   sync.Mutex
	heap priorityHeap
}

// NewPriorityQueue creates an empty priority queue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		heap: make(priorityHeap, 0),
	}
	heap.Init(&pq.heap)
	return pq
}

// Enqueue adds an item with the given priority.
func (pq *PriorityQueue) Enqueue(value interface{}, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item := &PriorityItem{
		Value:    value,
		Priority: priority,
	}
	heap.Push(&pq.heap, item)
}

// Dequeue removes and returns the highest-priority item.
// Returns nil if the queue is empty.
func (pq *PriorityQueue) Dequeue() interface{} {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.heap.Len() == 0 {
		return nil
	}

	item := heap.Pop(&pq.heap).(*PriorityItem)
	return item.Value
}

// Len returns the number of items in the queue.
func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.heap.Len()
}

// IsEmpty returns true if the queue has no items.
func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() == 0
}
