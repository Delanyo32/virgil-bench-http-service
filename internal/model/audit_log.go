package model

import (
	"fmt"
	"time"
)

// AuditAction enumerates the types of auditable actions.
type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
)

// AuditLog records a single auditable event.
type AuditLog struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Action     AuditAction `json:"action"`
	EntityType string      `json:"entity_type"`
	EntityID   int         `json:"entity_id"`
	Detail     string      `json:"detail,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

// Summary returns a human-readable summary of the audit event.
func (a *AuditLog) Summary() string {
	return fmt.Sprintf("user %d %s %s %d", a.UserID, a.Action, a.EntityType, a.EntityID)
}

// IsWrite returns true if the action modifies data.
func (a *AuditLog) IsWrite() bool {
	return a.Action == AuditActionCreate ||
		a.Action == AuditActionUpdate ||
		a.Action == AuditActionDelete
}

// NewAuditLog creates an AuditLog entry with the current timestamp.
func NewAuditLog(userID int, action AuditAction, entityType string, entityID int) AuditLog {
	return AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		CreatedAt:  time.Now(),
	}
}
