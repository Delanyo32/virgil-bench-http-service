package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/example/ordersvc/internal/model"
)

// UserRepository handles user persistence.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID retrieves a user by their ID.
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(user *model.User) error {
	err := r.db.QueryRow(
		`INSERT INTO users (email, name, password_hash, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		user.Email, user.Name, user.PasswordHash, time.Now(), time.Now(),
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update modifies a user's name and email.
func (r *UserRepository) Update(user *model.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4`,
		user.Name, user.Email, time.Now(), user.ID,
	)
	return err
}

// ListAll returns all users.
func (r *UserRepository) ListAll() ([]model.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, name, password_hash, created_at, updated_at FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
