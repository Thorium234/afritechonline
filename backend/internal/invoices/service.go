package invoices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Status constants for invoices.
const (
	StatusPending  = "PENDING"
	StatusPaid     = "PAID"
	StatusOverdue  = "OVERDUE"
	StatusCancelled = "CANCELLED"
)

// SubscriptionProvider resolves subscription details for invoice generation.
type SubscriptionProvider interface {
	Get(ctx context.Context, id uint64) (*models.Subscription, error)
}

// Service implements invoice business rules.
type Service struct {
	repo          *Repository
	subscriptions SubscriptionProvider
}

// NewService creates an invoice service.
func NewService(repo *Repository, subscriptions SubscriptionProvider) *Service {
	return &Service{repo: repo, subscriptions: subscriptions}
}

// CreateForSubscription generates an invoice for a subscription.
func (s *Service) CreateForSubscription(ctx context.Context, subID uint64) (*models.Invoice, error) {
	sub, err := s.subscriptions.Get(ctx, subID)
	if err != nil {
		return nil, err
	}

	number, err := s.repo.NextNumber(ctx)
	if err != nil {
		return nil, err
	}
	inv := &models.Invoice{
		InvoiceNo:      number,
		SubscriptionID: sub.ID,
		CustomerID:     sub.CustomerID,
		Amount:         sub.Amount,
		Currency:       sub.Currency,
		Status:         StatusPending,
		DueDate:        sub.ExpiryDate,
	}
	created, err := s.repo.Create(ctx, inv)
	if err != nil {
		if isDuplicateSubscription(err) {
			return s.repo.GetBySubscription(ctx, subID)
		}
		return nil, err
	}
	return created, nil
}

// Get returns an invoice.
func (s *Service) Get(ctx context.Context, id uint64) (*models.Invoice, error) {
	inv, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return inv, err
}

// List returns invoices with filters.
func (s *Service) List(ctx context.Context, customerID uint64, status string, page, pageSize int) ([]*models.Invoice, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, customerID, status, pageSize, offset)
}

// MarkPaid flags an invoice as paid.
func (s *Service) MarkPaid(ctx context.Context, id uint64) error {
	return s.repo.MarkPaid(ctx, id)
}

func isDuplicateSubscription(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, sql.ErrNoRows) == false && (contains(err.Error(), "Duplicate") || contains(err.Error(), "uq_invoice_subscription"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
