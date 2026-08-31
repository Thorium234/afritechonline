package mikrotik

import (
	"context"
	"fmt"
)

// ProfileService manages bandwidth profiles on MikroTik devices.
type ProfileService struct {
	client *Client
}

// NewProfileService creates a MikroTik profile service.
func NewProfileService(client *Client) *ProfileService {
	return &ProfileService{client: client}
}

// CreateProfile creates a hotspot user profile with rate limits.
func (s *ProfileService) CreateProfile(ctx context.Context, name, rateLimit string) error {
	cmd := fmt.Sprintf("/ip hotspot profile add name=%s rate-limit=%s", name, rateLimit)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}

// DeleteProfile removes a hotspot user profile.
func (s *ProfileService) DeleteProfile(ctx context.Context, name string) error {
	cmd := fmt.Sprintf("/ip hotspot profile remove [find name=%s]", name)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}
