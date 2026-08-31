package radius

import (
	"context"
	"fmt"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Service implements RADIUS business rules.
type Service struct {
	repo *Repository
}

// NewService creates a radius service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateUser adds a new RADIUS user.
func (s *Service) CreateUser(ctx context.Context, u *models.RadiusUser) (*models.RadiusUser, error) {
	return s.repo.Create(ctx, u)
}

// GetUser retrieves a RADIUS user.
func (s *Service) GetUser(ctx context.Context, username string) (*models.RadiusUser, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

// UpdateUser modifies a RADIUS user.
func (s *Service) UpdateUser(ctx context.Context, u *models.RadiusUser) (*models.RadiusUser, error) {
	return s.repo.Update(ctx, u)
}

// DeleteUser removes a RADIUS user.
func (s *Service) DeleteUser(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, username)
}

// RefreshExpired disables expired users.
func (s *Service) RefreshExpired(ctx context.Context) (int64, error) {
	usernames, err := s.repo.ListExpired(ctx)
	if err != nil {
		return 0, err
	}
	for _, username := range usernames {
		_, _ = s.repo.Update(ctx, &models.RadiusUser{
			Username: username,
			Profile:  "disabled",
		})
	}
	return int64(len(usernames)), nil
}

// ProvisionFromSubscription creates a RADIUS user from a subscription.
func (s *Service) ProvisionFromSubscription(ctx context.Context, sub *models.Subscription, pkg *models.InternetPackage) (*models.RadiusUser, error) {
	user := &models.RadiusUser{
		Username:  "",
		Password:  "",
		Profile:   pkg.Name,
		Speed:     fmt.Sprintf("%dM/%dM", pkg.DownloadMbps, pkg.UploadMbps),
		ExpiryDate: &sub.ExpiryDate,
	}
	return s.repo.Create(ctx, user)
}
