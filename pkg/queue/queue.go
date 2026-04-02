package queue

// Queue defines the interface for job queue operations.
type Queue interface {
	Enqueue(item interface{}) error
	Dequeue(count int) []interface{}
	Size() int
}
