package payments

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Status constants for payments.
const (
	StatusPending   = "PENDING"
	StatusCompleted = "COMPLETED"
	StatusFailed    = "FAILED"
	StatusCancelled = "CANCELLED"
)

// InvoiceService marks invoices paid and fetches them.
type InvoiceService interface {
	Get(ctx context.Context, id uint64) (*models.Invoice, error)
	MarkPaid(ctx context.Context, id uint64) error
}

// SubscriptionService activates subscriptions.
type SubscriptionService interface {
	Get(ctx context.Context, id uint64) (*models.Subscription, error)
	Activate(ctx context.Context, subID uint64) (*models.Subscription, error)
}

// ErrInvalidState indicates a payment cannot be completed in its current state.
var ErrInvalidState = errors.New("invalid payment state")

// Service orchestrates payment recording and completion.
// Completion is idempotent and activates the associated subscription exactly once.
type Service struct {
	repo          *Repository
	invoices      InvoiceService
	subscriptions SubscriptionService
}

// NewService creates a payment service.
func NewService(repo *Repository, invoices InvoiceService, subscriptions SubscriptionService) *Service {
	return &Service{repo: repo, invoices: invoices, subscriptions: subscriptions}
}

// Create records a new payment against an invoice (PENDING by default).
func (s *Service) Create(ctx context.Context, p *models.Payment) (*models.Payment, error) {
	return s.repo.Create(ctx, p)
}

// Get returns a payment.
func (s *Service) Get(ctx context.Context, id uint64) (*models.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

// List returns payments with filters.
func (s *Service) List(ctx context.Context, customerID uint64, status string, page, pageSize int) ([]*models.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, customerID, status, pageSize, offset)
}

// Complete marks a payment COMPLETED and, if it was not already completed,
// marks the invoice PAID and activates the subscription. Idempotent: a second
// call for an already COMPLETED payment is a no-op.
func (s *Service) Complete(ctx context.Context, paymentID uint64) (*models.Payment, error) {
	p, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status == StatusCompleted {
		return p, nil // idempotent
	}
	if p.Status == StatusFailed || p.Status == StatusCancelled {
		return nil, ErrInvalidState
	}

	if err := s.repo.MarkCompleted(ctx, p.ID); err != nil {
		return nil, err
	}

	// Mark the invoice paid (idempotent-friendly).
	inv, err := s.invoices.Get(ctx, p.InvoiceID)
	if err != nil {
		return nil, err
	}
	if inv.Status != "PAID" {
		if err := s.invoices.MarkPaid(ctx, inv.ID); err != nil {
			return nil, err
		}
	}

	// Activate the subscription (idempotent inside the subscription service).
	if _, err := s.subscriptions.Activate(ctx, inv.SubscriptionID); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, p.ID)
}
