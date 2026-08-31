package mikrotik

import (
	"context"
	"fmt"
)

// UserService manages network users on MikroTik devices.
type UserService struct {
	client *Client
}

// NewUserService creates a MikroTik user service.
func NewUserService(client *Client) *UserService {
	return &UserService{client: client}
}

// CreateUser adds a new hotspot/PPPoE user to the router.
func (s *UserService) CreateUser(ctx context.Context, username, password, profile string) error {
	cmd := fmt.Sprintf("/ip hotspot user add name=%s password=%s profile=%s", username, password, profile)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}

// DeleteUser removes a user from the router.
func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	cmd := fmt.Sprintf("/ip hotspot user remove [find name=%s]", username)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}

// DisableUser disables a user account.
func (s *UserService) DisableUser(ctx context.Context, username string) error {
	cmd := fmt.Sprintf("/ip hotspot user disable [find name=%s]", username)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}

// EnableUser enables a disabled user account.
func (s *UserService) EnableUser(ctx context.Context, username string) error {
	cmd := fmt.Sprintf("/ip hotspot user enable [find name=%s]", username)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}
