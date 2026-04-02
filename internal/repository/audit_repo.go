package repository

import (
	"sync"
	"time"

	"github.com/example/ordersvc/internal/model"
)

// AuditRepository provides append-only storage for audit log entries.
type AuditRepository struct {
	mu      sync.RWMutex
	entries []model.AuditLog
	nextID  int
}

// NewAuditRepository creates an empty AuditRepository.
func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		entries: make([]model.AuditLog, 0),
		nextID:  1,
	}
}

// Append stores a new audit log entry and assigns an ID.
func (r *AuditRepository) Append(entry *model.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry.ID = r.nextID
	r.nextID++
	r.entries = append(r.entries, *entry)
	return nil
}

// ListAll returns all audit log entries.
func (r *AuditRepository) ListAll() ([]model.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.AuditLog, len(r.entries))
	copy(result, r.entries)
	return result, nil
}

// ListByUser returns audit entries for the specified user ID.
func (r *AuditRepository) ListByUser(userID int) ([]model.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.AuditLog
	for _, e := range r.entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

// ListSince returns audit entries created after the given time.
func (r *AuditRepository) ListSince(since time.Time) ([]model.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.AuditLog
	for _, e := range r.entries {
		if e.CreatedAt.After(since) {
			result = append(result, e)
		}
	}
	return result, nil
}

// Count returns the total number of audit entries.
func (r *AuditRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
