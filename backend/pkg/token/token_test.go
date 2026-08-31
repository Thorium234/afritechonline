package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAndParseAccessToken(t *testing.T) {
	mgr := New("test-secret", 15, 168)

	token, err := mgr.IssueAccessToken(42, "ADMIN", "jdoe")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := mgr.ParseAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), claims.UserID)
	assert.Equal(t, "ADMIN", claims.Role)
	assert.Equal(t, "jdoe", claims.Username)
	assert.Equal(t, "afritech-online", claims.Issuer)
}

func TestParseAccessTokenInvalid(t *testing.T) {
	mgr := New("test-secret", 15, 168)

	_, err := mgr.ParseAccessToken("not-a-token")
	assert.Error(t, err)
}

func TestParseAccessTokenWrongSecret(t *testing.T) {
	mgr1 := New("secret-one", 15, 168)
	mgr2 := New("secret-two", 15, 168)

	token, err := mgr1.IssueAccessToken(1, "CUSTOMER", "alice")
	require.NoError(t, err)

	_, err = mgr2.ParseAccessToken(token)
	assert.Error(t, err)
}

func TestRefreshTokenTTL(t *testing.T) {
	mgr := New("test-secret", 15, 168)
	expected := 168 * time.Hour
	assert.Equal(t, expected, mgr.RefreshTokenTTL())
}

func TestHashRefreshToken(t *testing.T) {
	h1 := HashRefreshToken("my-token")
	h2 := HashRefreshToken("my-token")
	assert.Equal(t, h1, h2, "same input must produce same hash")
	assert.Len(t, h1, 44, "SHA-256 base64 raw URL is 44 chars")
}

func TestClaimsExpiration(t *testing.T) {
	mgr := New("test-secret", 15, 168)

	token, err := mgr.IssueAccessToken(1, "CUSTOMER", "bob")
	require.NoError(t, err)

	claims, err := mgr.ParseAccessToken(token)
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), claims.ExpiresAt.Time, time.Minute)
}

func TestTokenSignedWithHMAC(t *testing.T) {
	mgr := New("test-secret", 15, 168)
	token, err := mgr.IssueAccessToken(1, "STAFF", "carol")
	require.NoError(t, err)

	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		assert.True(t, ok, "expected HMAC signing method")
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
}
