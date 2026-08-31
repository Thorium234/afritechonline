package reports

import (
	"context"
	"database/sql"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Repository provides data access for reporting.
type Repository struct {
	db *sql.DB
}

// New creates a reports repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// RevenueSummary returns revenue grouped by day for the last N days.
func (r *Repository) RevenueSummary(ctx context.Context, days int) ([]*models.RevenueSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DATE(paid_at) as date, SUM(amount) as total, COUNT(*) as count
		 FROM payments
		 WHERE status = 'COMPLETED' AND paid_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		 GROUP BY DATE(paid_at)
		 ORDER BY date DESC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*models.RevenueSummary
	for rows.Next() {
		var s models.RevenueSummary
		if err := rows.Scan(&s.Date, &s.Total, &s.Count); err != nil {
			return nil, err
		}
		summaries = append(summaries, &s)
	}
	return summaries, rows.Err()
}

// CustomerStats returns customer statistics.
func (r *Repository) CustomerStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	rows, err := r.db.QueryContext(ctx, `SELECT status, COUNT(*) FROM customers GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, rows.Err()
}

// ActiveRouters returns count of online routers.
func (r *Repository) ActiveRouters(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM routers WHERE status = ?`, models.RouterStatusOnline).Scan(&count)
	return count, err
}
