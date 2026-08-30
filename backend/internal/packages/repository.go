package packages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a package does not exist.
var ErrNotFound = errors.New("package not found")

// Repository provides data access for internet packages.
type Repository struct {
	db *sql.DB
}

// New creates a package repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const cols = `id, name, description, price, currency, duration_days, download_mbps, upload_mbps, data_limit_gb, is_active, created_at, updated_at`

func scanPackage(row interface{ Scan(...any) error }) (*models.InternetPackage, error) {
	var p models.InternetPackage
	var limit sql.NullInt64
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Currency,
		&p.DurationDays, &p.DownloadMbps, &p.UploadMbps, &limit, &p.IsActive,
		&p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	if limit.Valid {
		v := int(limit.Int64)
		p.DataLimitGB = &v
	}
	return &p, nil
}

// Create inserts a new package.
func (r *Repository) Create(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	if p.Currency == "" {
		p.Currency = "KES"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO internet_packages (name, description, price, currency, duration_days, download_mbps, upload_mbps, data_limit_gb, is_active)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.Currency, p.DurationDays, p.DownloadMbps, p.UploadMbps, p.DataLimitGB, p.IsActive)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	p.ID = uint64(id)
	return p, nil
}

// List returns all packages, optionally filtered by active status.
func (r *Repository) List(ctx context.Context, activeOnly bool) ([]*models.InternetPackage, error) {
	query := `SELECT ` + cols + ` FROM internet_packages`
	args := []any{}
	if activeOnly {
		query += ` WHERE is_active = 1`
	}
	query += ` ORDER BY price ASC`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.InternetPackage
	for rows.Next() {
		p, err := scanPackage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetByID returns a package by ID.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.InternetPackage, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM internet_packages WHERE id = ?`, id)
	return scanPackage(row)
}

// Update modifies a package.
func (r *Repository) Update(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE internet_packages SET name=?, description=?, price=?, duration_days=?, download_mbps=?, upload_mbps=?, data_limit_gb=?, is_active=? WHERE id=?`,
		p.Name, p.Description, p.Price, p.DurationDays, p.DownloadMbps, p.UploadMbps, p.DataLimitGB, p.IsActive, p.ID)
	if err != nil {
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, p.ID)
}

// Delete removes a package.
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM internet_packages WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
