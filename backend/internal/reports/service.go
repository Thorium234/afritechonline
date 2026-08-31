package reports

import (
	"context"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Service implements reporting business rules.
type Service struct {
	repo *Repository
}

// NewService creates a reports service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RevenueSummary returns revenue summary for the last N days.
func (s *Service) RevenueSummary(ctx context.Context, days int) ([]*models.RevenueSummary, error) {
	return s.repo.RevenueSummary(ctx, days)
}

// CustomerStats returns customer statistics.
func (s *Service) CustomerStats(ctx context.Context) (map[string]int64, error) {
	return s.repo.CustomerStats(ctx)
}

// ActiveRouters returns count of online routers.
func (s *Service) ActiveRouters(ctx context.Context) (int64, error) {
	return s.repo.ActiveRouters(ctx)
}
