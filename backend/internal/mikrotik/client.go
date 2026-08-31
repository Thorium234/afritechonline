package mikrotik

import (
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client manages a long-lived SSH session to a MikroTik device.
type Client struct {
	host     string
	username string
	password string
	timeout  time.Duration
}

// NewClient creates a new MikroTik SSH client configuration.
func NewClient(host, username, password string, timeout time.Duration) *Client {
	return &Client{
		host:     host,
		username: username,
		password: password,
		timeout:  timeout,
	}
}

// Connect establishes an SSH connection and returns a session.
func (c *Client) Connect(ctx context.Context) (*ssh.Session, error) {
	cfg := &ssh.ClientConfig{
		User:            c.username,
		Auth:            []ssh.AuthMethod{ssh.Password(c.password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}

	addr := net.JoinHostPort(c.host, fmt.Sprintf("%d", 8728))
	conn, err := dialContext(ctx, "tcp", addr, cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("dial mikrotik: %w", err)
	}

	sshClient, err := ssh.Dial(conn, cfg)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ssh handshake: %w", err)
	}

	session, err := sshClient.NewSession()
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("open session: %w", err)
	}

	return session, nil
}

// TestConnection verifies basic connectivity and returns identity/version.
func (c *Client) TestConnection(ctx context.Context) (string, string, error) {
	session, err := c.Connect(ctx)
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput("/system identity print")
	if err != nil {
		return "", "", fmt.Errorf("identity command failed: %w", err)
	}

	identity := parseIdentity(string(output))
	version, _ := c.getVersion(ctx)
	return identity, version, nil
}

func (c *Client) getVersion(ctx context.Context) (string, error) {
	session, err := c.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput("/system resource print")
	if err != nil {
		return "", err
	}

	lines := splitLines(string(output))
	for _, line := range lines {
		if contains(line, "version:") {
			return trimVersion(line), nil
		}
	}
	return "unknown", nil
}

func dialContext(ctx context.Context, network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	return d.DialContext(ctx, network, address)
}

func parseIdentity(output string) string {
	lines := splitLines(output)
	for _, line := range lines {
		if contains(line, "name:") {
			return trimAfter(line, "name:")
		}
	}
	return "unknown"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func trimAfter(s, delim string) string {
	if idx := indexOf(s, delim); idx >= 0 {
		s = s[idx+len(delim):]
	}
	var b []rune
	for _, r := range s {
		if r != ' ' && r != '\t' {
			break
		}
		b = append(b, r)
	}
	return string(b)
}

func trimVersion(s string) string {
	fields := splitLines(s)
	for _, f := range fields {
		if contains(f, "version:") {
			parts := splitAfter(f, "version:")
			if len(parts) > 1 {
				return trimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func splitAfter(s, delim string) []string {
	var parts []string
	if idx := indexOf(s, delim); idx >= 0 {
		parts = append(parts, s[:idx], s[idx+len(delim):])
	} else {
		parts = append(parts, s)
	}
	return parts
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
