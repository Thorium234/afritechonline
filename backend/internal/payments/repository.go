package payments

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a payment does not exist.
var ErrNotFound = errors.New("payment not found")

// Repository provides data access for payments.
type Repository struct {
	db *sql.DB
}

// New creates a payment repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const cols = `id, invoice_id, customer_id, amount, currency, method, reference, status, paid_at, created_at, updated_at`

func scanPayment(row interface{ Scan(...any) error }) (*models.Payment, error) {
	var p models.Payment
	var paidAt sql.NullTime
	if err := row.Scan(&p.ID, &p.InvoiceID, &p.CustomerID, &p.Amount, &p.Currency,
		&p.Method, &p.Reference, &p.Status, &paidAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	if paidAt.Valid {
		t := paidAt.Time
		p.PaidAt = &t
	}
	return &p, nil
}

// Create inserts a new payment.
func (r *Repository) Create(ctx context.Context, p *models.Payment) (*models.Payment, error) {
	if p.Currency == "" {
		p.Currency = "KES"
	}
	if p.Status == "" {
		p.Status = "PENDING"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (invoice_id, customer_id, amount, currency, method, reference, status, paid_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.InvoiceID, p.CustomerID, p.Amount, p.Currency, p.Method, p.Reference, p.Status, p.PaidAt)
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

// GetByID returns a payment.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.Payment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM payments WHERE id = ?`, id)
	return scanPayment(row)
}

// GetByReference returns a payment by its external reference + method.
func (r *Repository) GetByReference(ctx context.Context, reference, method string) (*models.Payment, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+cols+` FROM payments WHERE reference = ? AND method = ?`, reference, method)
	return scanPayment(row)
}

// List returns payments with pagination.
func (r *Repository) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Payment, int64, error) {
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
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + cols + ` FROM payments` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Payment
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

// MarkCompleted sets a payment as COMPLETED with a paid timestamp.
func (r *Repository) MarkCompleted(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE payments SET status = 'COMPLETED', paid_at = ? WHERE id = ?`, time.Now(), id)
	return err
}

// MarkFailed sets a payment as FAILED.
func (r *Repository) MarkFailed(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE payments SET status = 'FAILED' WHERE id = ?`, id)
	return err
}

// Begin starts a transaction.
func (r *Repository) Begin(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// WithTx returns a repository bound to the provided transaction.
func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}
