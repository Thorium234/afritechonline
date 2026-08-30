package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrInvalidState is returned when a subscription is not in a transitionable state.
var ErrInvalidState = errors.New("invalid subscription state")

// Status constants for subscriptions.
const (
	StatusPending   = "PENDING"
	StatusActive    = "ACTIVE"
	StatusExpired   = "EXPIRED"
	StatusSuspended = "SUSPENDED"
	StatusCancelled = "CANCELLED"
)

// PackageProvider resolves package details used when creating subscriptions.
type PackageProvider interface {
	GetActive(ctx context.Context, id uint64) (*models.InternetPackage, error)
}

// CustomerProvider validates that a customer exists.
type CustomerProvider interface {
	Get(ctx context.Context, id uint64) (*models.Customer, error)
	SetActive(ctx context.Context, id uint64, active bool) error
}

// Service implements subscription business rules.
type Service struct {
	repo     *Repository
	packages PackageProvider
	customers CustomerProvider
	now      func() time.Time
}

// NewService creates a subscription service.
func NewService(repo *Repository, packages PackageProvider, customers CustomerProvider) *Service {
	return &Service{repo: repo, packages: packages, customers: customers, now: time.Now}
}

// Create builds and persists a PENDING subscription from a package selection.
func (s *Service) Create(ctx context.Context, customerID, packageID uint64) (*models.Subscription, error) {
	// Verify customer exists.
	if _, err := s.customers.Get(ctx, customerID); err != nil {
		return nil, err
	}
	// Verify package is active and resolve price/duration.
	pkg, err := s.packages.GetActive(ctx, packageID)
	if err != nil {
		return nil, err
	}

	now := s.now()
	sub := &models.Subscription{
		CustomerID: customerID,
		PackageID:  packageID,
		StartDate:  now,
		ExpiryDate: now.AddDate(0, 0, pkg.DurationDays),
		Status:     StatusPending,
		Amount:     pkg.Price,
		Currency:   pkg.Currency,
	}
	return s.repo.Create(ctx, sub)
}

// Get returns a subscription.
func (s *Service) Get(ctx context.Context, id uint64) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return sub, err
}

// List returns subscriptions with filters.
func (s *Service) List(ctx context.Context, customerID uint64, status string, page, pageSize int) ([]*models.Subscription, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, customerID, status, pageSize, offset)
}

// Activate transitions a PENDING subscription to ACTIVE on successful payment.
// It also marks the customer ACTIVE.
func (s *Service) Activate(ctx context.Context, subID uint64) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, subID)
	if err != nil {
		return nil, err
	}
	// Idempotency: if already ACTIVE, do not create duplicate side effects.
	if sub.Status == StatusActive {
		return sub, nil
	}
	if sub.Status != StatusPending {
		return nil, ErrInvalidState
	}

	if err := s.repo.UpdateStatus(ctx, sub.ID, StatusActive); err != nil {
		return nil, err
	}
	if err := s.customers.SetActive(ctx, sub.CustomerID, true); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, subID)
}

// Expire transitions an ACTIVE subscription to EXPIRED.
func (s *Service) Expire(ctx context.Context, subID uint64) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, subID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateStatus(ctx, sub.ID, StatusExpired); err != nil {
		return nil, err
	}
	if err := s.customers.SetActive(ctx, sub.CustomerID, false); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, subID)
}

// RefreshExpiredStatus marks any ACTIVE subscriptions that have passed their
// expiry date as EXPIRED. Returns the number changed.
func (s *Service) RefreshExpiredStatus(ctx context.Context) (int64, error) {
	expired, err := s.repo.ListExpired(ctx, s.now())
	if err != nil {
		return 0, err
	}
	for _, id := range expired {
		if err := s.repo.UpdateStatus(ctx, id, StatusExpired); err != nil {
			return 0, err
		}
	}
	return int64(len(expired)), nil
}
