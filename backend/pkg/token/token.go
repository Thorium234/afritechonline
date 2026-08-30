package token

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalid indicates a token could not be validated.
var ErrInvalid = errors.New("invalid token")

// Claims is the JWT payload for access tokens.
type Claims struct {
	UserID   uint64 `json:"uid"`
	Role     string `json:"role"`
	Username string `json:"uname"`
	jwt.RegisteredClaims
}

// Manager issues and validates access tokens and refresh-token hashes.
type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// New creates a token manager.
func New(secret string, accessTTLMins, refreshTTLHours int) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  time.Duration(accessTTLMins) * time.Minute,
		refreshTTL: time.Duration(refreshTTLHours) * time.Hour,
	}
}

// IssueAccessToken creates a signed JWT for the given user.
func (m *Manager) IssueAccessToken(userID uint64, role, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Role:     role,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "afritech-online",
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

// ParseAccessToken validates a JWT and returns its claims.
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalid
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalid
	}
	if !t.Valid {
		return nil, ErrInvalid
	}
	return claims, nil
}

// RefreshTokenTTL returns the refresh token lifetime duration.
func (m *Manager) RefreshTokenTTL() time.Duration {
	return m.refreshTTL
}

// HashRefreshToken returns a URL-safe SHA-256 hash for a refresh token.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
