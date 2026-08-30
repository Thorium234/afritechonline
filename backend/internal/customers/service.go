package customers

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// ErrDuplicate indicates a phone or username already exists.
var ErrDuplicate = errors.New("duplicate customer")

// Service wraps the repository with business rules and error mapping.
type Service struct {
	repo *Repository
}

// NewService creates a customer service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListWithPagination returns customers and total count.
func (s *Service) ListWithPagination(ctx context.Context, search string, page, pageSize int) ([]*models.Customer, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, search, pageSize, offset)
}

// Get returns a single customer.
func (s *Service) Get(ctx context.Context, id uint64) (*models.Customer, error) {
	c, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

// Create registers a new customer.
func (s *Service) Create(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	if c.Status == "" {
		c.Status = "INACTIVE"
	}
	created, err := s.repo.Create(ctx, c)
	if err != nil {
		if isDuplicate(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return created, nil
}

// Update modifies a customer.
func (s *Service) Update(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	updated, err := s.repo.Update(ctx, c)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil && isDuplicate(err) {
		return nil, ErrDuplicate
	}
	return updated, err
}

// Delete removes a customer.
func (s *Service) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// SetActive flips the customer's status between ACTIVE and INACTIVE.
func (s *Service) SetActive(ctx context.Context, id uint64, active bool) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if active && c.Status == "SUSPENDED" {
		return nil
	}
	status := "INACTIVE"
	if active {
		status = "ACTIVE"
	}
	c.Status = status
	_, err = s.repo.Update(ctx, c)
	return err
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
