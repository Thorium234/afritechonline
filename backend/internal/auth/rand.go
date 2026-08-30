package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRefreshToken returns a cryptographically random opaque token.
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
