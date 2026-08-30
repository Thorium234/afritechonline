package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/Thorium234/afritechonline/backend/pkg/token"
)

// RefreshTokenStore persists refresh tokens and their hashes.
type RefreshTokenStore struct {
	db *sql.DB
}

// NewRefreshTokenStore creates a refresh token store.
func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{db: db}
}

// Create stores a new refresh token hash for a user.
func (s *RefreshTokenStore) Create(ctx context.Context, userID uint64, raw string, ttl time.Duration) error {
	hash := token.HashRefreshToken(raw)
	expires := time.Now().Add(ttl)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)`,
		userID, hash, expires)
	return err
}

// Validate checks a refresh token hash exists, is not expired and not revoked.
func (s *RefreshTokenStore) Validate(ctx context.Context, raw string) (uint64, error) {
	hash := token.HashRefreshToken(raw)
	var userID uint64
	var expiresAt time.Time
	var revokedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = ?`, hash).
		Scan(&userID, &expiresAt, &revokedAt)
	if err != nil {
		return 0, err
	}
	if expired(expiresAt) || revokedAt.Valid {
		return 0, sql.ErrNoRows
	}
	return userID, nil
}

// Revoke invalidates a refresh token.
func (s *RefreshTokenStore) Revoke(ctx context.Context, raw string) error {
	hash := token.HashRefreshToken(raw)
	_, err := s.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ? WHERE token_hash = ? AND revoked_at IS NULL`,
		time.Now(), hash)
	return err
}

// RevokeAllForUser invalidates every active refresh token for a user (logout all).
func (s *RefreshTokenStore) RevokeAllForUser(ctx context.Context, userID uint64) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`,
		time.Now(), userID)
	return err
}

func expired(t time.Time) bool {
	return time.Now().After(t)
}
