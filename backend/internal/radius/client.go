package radius

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Client handles RADIUS protocol operations.
type Client struct {
	host     string
	secret   string
	authPort int
	acctPort int
}

// NewClient creates a new RADIUS client.
func NewClient(host, secret string, authPort, acctPort int) *Client {
	return &Client{
		host:     host,
		secret:   secret,
		authPort: authPort,
		acctPort: acctPort,
	}
}

// Authenticate checks username/password against FreeRADIUS.
func (c *Client) Authenticate(ctx context.Context, username, password string) (bool, error) {
	// Phase 4 implementation will use a proper RADIUS library.
	// For now, this is a placeholder that validates the request structure.
	if username == "" || password == "" {
		return false, fmt.Errorf("username and password are required")
	}
	return true, nil
}

// AccountingStart records session start.
func (c *Client) AccountingStart(ctx context.Context, username, sessionID, framedIP string) error {
	return nil
}

// AccountingStop records session stop.
func (c *Client) AccountingStop(ctx context.Context, username, sessionID string, bytesIn, bytesOut uint64) error {
	return nil
}

// GeneratePassword returns an MD5 hash of the shared secret + request authenticator.
func GeneratePassword(secret, requestAuthenticator, password string) string {
	hash := md5.Sum([]byte(secret + requestAuthenticator + password))
	return hex.EncodeToString(hash[:])
}
