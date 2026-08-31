package mikrotik

import (
	"context"
	"fmt"
	"strings"
)

// SessionService manages active sessions on MikroTik devices.
type SessionService struct {
	client *Client
}

// NewSessionService creates a MikroTik session service.
func NewSessionService(client *Client) *SessionService {
	return &SessionService{client: client}
}

// ActiveSessions returns a list of active hotspot sessions.
func (s *SessionService) ActiveSessions(ctx context.Context) ([]string, error) {
	session, err := s.client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output, err := session.CombinedOutput("/ip hotspot active print detail")
	if err != nil {
		return nil, err
	}

	var sessions []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "user=") || strings.Contains(line, "address=") {
			sessions = append(sessions, strings.TrimSpace(line))
		}
	}
	return sessions, nil
}

// DisconnectSession terminates an active session by user name.
func (s *SessionService) DisconnectSession(ctx context.Context, username string) error {
	cmd := fmt.Sprintf("/ip hotspot active remove [find user=%s]", username)
	session, err := s.client.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}
