package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a user does not exist.
var ErrNotFound = errors.New("user not found")

// Repository provides data access for users.
type Repository struct {
	db *sql.DB
}

// New creates a users repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const userColumns = `id, username, email, password_hash, role, is_active, last_login_at, created_at, updated_at`

func scanUser(row interface{ Scan(...any) error }) (*models.User, error) {
	var u models.User
	var lastLogin sql.NullTime
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.IsActive, &lastLogin, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		t := lastLogin.Time
		u.LastLoginAt = &t
	}
	return &u, nil
}

// Create inserts a new user and returns it with its ID.
func (r *Repository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (username, email, password_hash, role, is_active) VALUES (?, ?, ?, ?, ?)`,
		u.Username, u.Email, u.PasswordHash, u.Role, u.IsActive)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = uint64(id)
	return u, nil
}

// GetByUsername returns a user by username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE username = ?`, username)
	return scanUser(row)
}

// GetByEmail returns a user by email.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE email = ?`, email)
	return scanUser(row)
}

// GetByID returns a user by primary key.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE id = ?`, id)
	return scanUser(row)
}

// RecordLogin updates the last login timestamp.
func (r *Repository) RecordLogin(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET last_login_at = ? WHERE id = ?`, time.Now(), id)
	return err
}
