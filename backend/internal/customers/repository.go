package customers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrNotFound is returned when a customer does not exist.
var ErrNotFound = errors.New("customer not found")

// Repository provides data access for customers.
type Repository struct {
	db *sql.DB
}

// New creates a customer repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const cols = `id, user_id, full_name, phone, email, username, status, created_at, updated_at`

func scanCustomer(row interface{ Scan(...any) error }) (*models.Customer, error) {
	var c models.Customer
	var userID sql.NullInt64
	if err := row.Scan(&c.ID, &userID, &c.FullName, &c.Phone, &c.Email, &c.Username,
		&c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	if userID.Valid {
		id := uint64(userID.Int64)
		c.UserID = &id
	}
	return &c, nil
}

func (r *Repository) dbq(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// Create inserts a new customer.
func (r *Repository) Create(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	if c.Status == "" {
		c.Status = "INACTIVE"
	}
	var userID sql.NullInt64
	if c.UserID != nil {
		userID = sql.NullInt64{Int64: int64(*c.UserID), Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO customers (user_id, full_name, phone, email, username, status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID, c.FullName, c.Phone, c.Email, c.Username, c.Status)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.ID = uint64(id)
	return c, nil
}

// List returns customers with optional search.
func (r *Repository) List(ctx context.Context, search string, limit, offset int) ([]*models.Customer, int64, error) {
	where := ""
	args := []any{}
	if search != "" {
		where = " WHERE full_name LIKE ? OR phone LIKE ? OR username LIKE ? OR email LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT ` + cols + ` FROM customers` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.dbq(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, c)
	}
	return out, total, rows.Err()
}

// GetByID returns a customer by ID.
func (r *Repository) GetByID(ctx context.Context, id uint64) (*models.Customer, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM customers WHERE id = ?`, id)
	return scanCustomer(row)
}

// GetByPhone returns a customer by phone.
func (r *Repository) GetByPhone(ctx context.Context, phone string) (*models.Customer, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM customers WHERE phone = ?`, phone)
	return scanCustomer(row)
}

// Update updates customer mutable fields.
func (r *Repository) Update(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE customers SET full_name = ?, phone = ?, email = ?, username = ?, status = ? WHERE id = ?`,
		c.FullName, c.Phone, c.Email, c.Username, c.Status, c.ID)
	if err != nil {
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, c.ID)
}

// Delete removes a customer.
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
