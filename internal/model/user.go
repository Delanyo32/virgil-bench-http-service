package model

import (
	"time"
)

// User represents a registered customer.
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DisplayName returns the user's name or email if name is empty.
func (u *User) DisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}
