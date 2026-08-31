package mikrotik

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a router does not exist.
var ErrNotFound = errors.New("router not found")

// Repository provides data access for routers.
type Repository struct {
	db *sql.DB
}

// New creates a router repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const routerCols = `id, name, host, api_port, username, password_enc, location, status, created_at, updated_at`

func scanRouter(row interface{ Scan(...any) error }) (*models.Router, error) {
	var r models.Router
	if err := row.Scan(&r.ID, &r.Name, &r.Host, &r.APIPort, &r.Username,
		&r.PasswordEnc, &r.Location, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}

// Create inserts a new router.
func (r *Repository) Create(ctx context.Context, router *models.Router) (*models.Router, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO routers (name, host, api_port, username, password_enc, location, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		router.Name, router.Host, router.APIPort, router.Username, router.PasswordEnc, router.Location, router.Status)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	router.ID = uint64(id)
	return router, nil
}

// GetByID returns a router by ID.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.Router, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+routerCols+` FROM routers WHERE id = ?`, id)
	return scanRouter(row)
}

// List returns all routers.
func (r *Repository) List(ctx context.Context) ([]*models.Router, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+routerCols+` FROM routers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Router
	for rows.Next() {
		router, err := scanRouter(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, router)
	}
	return out, rows.Err()
}

// Update modifies a router.
func (r *Repository) Update(ctx context.Context, router *models.Router) (*models.Router, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE routers SET name=?, host=?, api_port=?, username=?, password_enc=?, location=?, status=? WHERE id=?`,
		router.Name, router.Host, router.APIPort, router.Username, router.PasswordEnc, router.Location, router.Status, router.ID)
	if err != nil {
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, router.ID)
}

// Delete removes a router.
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM routers WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateStatus changes a router's online status.
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE routers SET status = ? WHERE id = ?`, status, id)
	return err
}
