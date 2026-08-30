package packages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Service wraps the package repository with business rules.
type Service struct {
	repo *Repository
}

// NewService creates a package service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns packages.
func (s *Service) List(ctx context.Context, activeOnly bool) ([]*models.InternetPackage, error) {
	return s.repo.List(ctx, activeOnly)
}

// Get returns a single package.
func (s *Service) Get(ctx context.Context, id uint64) (*models.InternetPackage, error) {
	p, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

// GetActive returns an active package by ID.
func (s *Service) GetActive(ctx context.Context, id uint64) (*models.InternetPackage, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !p.IsActive {
		return nil, ErrNotFound
	}
	return p, nil
}

// Create adds a package.
func (s *Service) Create(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	return s.repo.Create(ctx, p)
}

// Update modifies a package.
func (s *Service) Update(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	return s.repo.Update(ctx, p)
}

// Delete removes a package.
func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
