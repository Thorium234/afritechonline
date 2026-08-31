package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Service orchestrates router operations.
type Service struct {
	repo      *Repository
	sshClient *Client
}

// NewService creates a MikroTik service.
func NewService(repo *Repository, sshClient *Client) *Service {
	return &Service{repo: repo, sshClient: sshClient}
}

// Register adds a new router to the database.
func (s *Service) Register(ctx context.Context, router *models.Router) (*models.Router, error) {
	if router.Status == "" {
		router.Status = models.RouterStatusUnknown
	}
	router.APIPort = defaultAPIPort(router.APIPort)
	return s.repo.Create(ctx, router)
}

// Get returns a router by ID.
func (s *Service) Get(ctx context.Context, id uint64) (*models.Router, error) {
	router, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return router, nil
}

// List returns all routers.
func (s *Service) List(ctx context.Context) ([]*models.Router, error) {
	return s.repo.List(ctx)
}

// Update modifies a router.
func (s *Service) Update(ctx context.Context, router *models.Router) (*models.Router, error) {
	return s.repo.Update(ctx, router)
}

// Delete removes a router.
func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// TestConnection verifies connectivity to the router.
func (s *Service) TestConnection(ctx context.Context, router *models.Router) (string, string, error) {
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	identity, version, err := client.TestConnection(ctx)
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, router.ID, models.RouterStatusOffline)
		return "", "", err
	}
	_ = s.repo.UpdateStatus(ctx, router.ID, models.RouterStatusOnline)
	return identity, version, nil
}

// ProvisionUser creates a network user on the router.
func (s *Service) ProvisionUser(ctx context.Context, routerID uint64, username, password, profile string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	userService := NewUserService(client)
	return userService.CreateUser(ctx, username, password, profile)
}

// RemoveUser deletes a network user from the router.
func (s *Service) RemoveUser(ctx context.Context, routerID uint64, username string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	userService := NewUserService(client)
	return userService.DeleteUser(ctx, username)
}

// EnableUser enables a network user on the router.
func (s *Service) EnableUser(ctx context.Context, routerID uint64, username string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	userService := NewUserService(client)
	return userService.EnableUser(ctx, username)
}

// DisableUser disables a network user on the router.
func (s *Service) DisableUser(ctx context.Context, routerID uint64, username string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	userService := NewUserService(client)
	return userService.DisableUser(ctx, username)
}

// CreateProfile creates a bandwidth profile on the router.
func (s *Service) CreateProfile(ctx context.Context, routerID uint64, name, rateLimit string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	profileService := NewProfileService(client)
	return profileService.CreateProfile(ctx, name, rateLimit)
}

// ActiveSessions returns active sessions on the router.
func (s *Service) ActiveSessions(ctx context.Context, routerID uint64) ([]string, error) {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	sessionService := NewSessionService(client)
	return sessionService.ActiveSessions(ctx)
}

// DisconnectSession terminates an active session on the router.
func (s *Service) DisconnectSession(ctx context.Context, routerID uint64, username string) error {
	router, err := s.repo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}
	client := NewClient(router.Host, router.Username, router.PasswordEnc, 5*time.Second)
	sessionService := NewSessionService(client)
	return sessionService.DisconnectSession(ctx, username)
}

func defaultAPIPort(port int) int {
	if port <= 0 {
		return 8728
	}
	return port
}
