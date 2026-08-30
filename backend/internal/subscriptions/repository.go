package subscriptions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a subscription does not exist.
var ErrNotFound = errors.New("subscription not found")

// Repository provides data access for subscriptions.
type Repository struct {
	db *sql.DB
}

// New creates a subscription repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const cols = `id, customer_id, package_id, start_date, expiry_date, status, amount, currency, created_at, updated_at`

func scanSubscription(row interface{ Scan(...any) error }) (*models.Subscription, error) {
	var s models.Subscription
	if err := row.Scan(&s.ID, &s.CustomerID, &s.PackageID, &s.StartDate, &s.ExpiryDate,
		&s.Status, &s.Amount, &s.Currency, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

// Create inserts a new subscription.
func (r *Repository) Create(ctx context.Context, s *models.Subscription) (*models.Subscription, error) {
	if s.Currency == "" {
		s.Currency = "KES"
	}
	if s.Status == "" {
		s.Status = "PENDING"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO subscriptions (customer_id, package_id, start_date, expiry_date, status, amount, currency)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		s.CustomerID, s.PackageID, s.StartDate, s.ExpiryDate, s.Status, s.Amount, s.Currency)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	s.ID = uint64(id)
	return s, nil
}

// GetByID returns a subscription.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.Subscription, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM subscriptions WHERE id = ?`, id)
	return scanSubscription(row)
}

// List returns subscriptions, optionally filtered by customer.
func (r *Repository) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Subscription, int64, error) {
	where := " WHERE 1=1"
	args := []any{}
	if customerID > 0 {
		where += " AND customer_id = ?"
		args = append(args, customerID)
	}
	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscriptions`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + cols + ` FROM subscriptions` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Subscription
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

// UpdateStatus changes a subscription status.
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE subscriptions SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// GetActiveForCustomer returns the currently active subscription for a customer, if any.
func (r *Repository) GetActiveForCustomer(ctx context.Context, customerID uint64) (*models.Subscription, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+cols+` FROM subscriptions WHERE customer_id = ? AND status = 'ACTIVE' ORDER BY created_at DESC LIMIT 1`,
		customerID)
	return scanSubscription(row)
}

// ListExpired returns the IDs of ACTIVE subscriptions that have passed a cut-off time.
func (r *Repository) ListExpired(ctx context.Context, cutoff interface{}) ([]uint64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id FROM subscriptions WHERE status = 'ACTIVE' AND expiry_date <= ?`, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
