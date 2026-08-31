package radius

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a RADIUS user does not exist.
var ErrNotFound = errors.New("radius user not found")

// Repository provides data access for RADIUS users.
type Repository struct {
	db *sql.DB
}

// New creates a radius repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const radiusUserCols = `id, username, password, profile, speed, expiry_date, created_at, updated_at`

func scanRadiusUser(row interface{ Scan(...any) error }) (*models.RadiusUser, error) {
	var u models.RadiusUser
	var expiry sql.NullTime
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Profile, &u.Speed, &expiry, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	if expiry.Valid {
		t := expiry.Time
		u.ExpiryDate = &t
	}
	return &u, nil
}

// Create inserts a new RADIUS user.
func (r *Repository) Create(ctx context.Context, u *models.RadiusUser) (*models.RadiusUser, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO radius_users (username, password, profile, speed, expiry_date) VALUES (?, ?, ?, ?, ?)`,
		u.Username, u.Password, u.Profile, u.Speed, u.ExpiryDate)
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

// GetByUsername returns a RADIUS user by username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*models.RadiusUser, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+radiusUserCols+` FROM radius_users WHERE username = ?`, username)
	return scanRadiusUser(row)
}

// Update modifies a RADIUS user.
func (r *Repository) Update(ctx context.Context, u *models.RadiusUser) (*models.RadiusUser, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE radius_users SET password=?, profile=?, speed=?, expiry_date=? WHERE username=?`,
		u.Password, u.Profile, u.Speed, u.ExpiryDate, u.Username)
	if err != nil {
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrNotFound
	}
	return r.GetByUsername(ctx, u.Username)
}

// Delete removes a RADIUS user.
func (r *Repository) Delete(ctx context.Context, username string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM radius_users WHERE username = ?`, username)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// ListExpired returns usernames that have expired.
func (r *Repository) ListExpired(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT username FROM radius_users WHERE expiry_date IS NOT NULL AND expiry_date <= ?`, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, rows.Err()
}
