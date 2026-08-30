package invoices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when an invoice does not exist.
var ErrNotFound = errors.New("invoice not found")

// Repository provides data access for invoices.
type Repository struct {
	db *sql.DB
}

// New creates an invoice repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const cols = `id, invoice_no, subscription_id, customer_id, amount, currency, status, due_date, created_at, updated_at`

func scanInvoice(row interface{ Scan(...any) error }) (*models.Invoice, error) {
	var inv models.Invoice
	if err := row.Scan(&inv.ID, &inv.InvoiceNo, &inv.SubscriptionID, &inv.CustomerID,
		&inv.Amount, &inv.Currency, &inv.Status, &inv.DueDate, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
		return nil, err
	}
	return &inv, nil
}

// NextNumber generates a unique sequential invoice number.
func (r *Repository) NextNumber(ctx context.Context) (string, error) {
	var seq int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(id), 0) + 1 FROM invoices`).Scan(&seq); err != nil {
		return "", err
	}
	t := time.Now()
	return fmt.Sprintf("INV-%s-%06d", t.Format("20060102"), seq), nil
}

// Create inserts a new invoice.
func (r *Repository) Create(ctx context.Context, inv *models.Invoice) (*models.Invoice, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO invoices (invoice_no, subscription_id, customer_id, amount, currency, status, due_date)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		inv.InvoiceNo, inv.SubscriptionID, inv.CustomerID, inv.Amount, inv.Currency, inv.Status, inv.DueDate)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	inv.ID = uint64(id)
	return inv, nil
}

// GetByID returns an invoice.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.Invoice, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM invoices WHERE id = ?`, id)
	return scanInvoice(row)
}

// GetBySubscription returns the latest invoice for a subscription.
func (r *Repository) GetBySubscription(ctx context.Context, subID uint64) (*models.Invoice, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+cols+` FROM invoices WHERE subscription_id = ? ORDER BY created_at DESC LIMIT 1`, subID)
	return scanInvoice(row)
}

// List returns invoices with pagination.
func (r *Repository) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Invoice, int64, error) {
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
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + cols + ` FROM invoices` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, inv)
	}
	return out, total, rows.Err()
}

// MarkPaid sets an invoice status to PAID.
func (r *Repository) MarkPaid(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE invoices SET status = 'PAID' WHERE id = ?`, id)
	return err
}
