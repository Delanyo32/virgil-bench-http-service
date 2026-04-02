package service

import (
	"fmt"
	"time"

	"github.com/example/ordersvc/internal/model"
	"github.com/example/ordersvc/internal/repository"
)

// AuditService records and queries audit log entries.
type AuditService struct {
	repo *repository.AuditRepository
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// RecordAction creates an audit log entry for the given action.
func (s *AuditService) RecordAction(userID int, action model.AuditAction, entityType string, entityID int) error {
	entry := model.NewAuditLog(userID, action, entityType, entityID)
	if err := s.repo.Append(&entry); err != nil {
		return fmt.Errorf("failed to record audit action: %w", err)
	}
	return nil
}

// GetUserHistory returns all audit entries for a user.
func (s *AuditService) GetUserHistory(userID int) ([]model.AuditLog, error) {
	entries, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user history: %w", err)
	}
	return entries, nil
}

// GetRecentActivity returns audit entries from the last given duration.
func (s *AuditService) GetRecentActivity(d time.Duration) ([]model.AuditLog, error) {
	since := time.Now().Add(-d)
	entries, err := s.repo.ListSince(since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	return entries, nil
}

// Count returns the total number of recorded audit events.
func (s *AuditService) Count() int {
	return s.repo.Count()
}
