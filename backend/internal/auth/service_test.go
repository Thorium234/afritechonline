package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/pkg/token"
)

type mockRefreshStore struct {
	createFunc    func(ctx context.Context, userID uint64, raw string, ttl time.Duration) error
	validateFunc  func(ctx context.Context, raw string) (uint64, error)
	revokeFunc    func(ctx context.Context, raw string) error
	revokeAllFunc func(ctx context.Context, userID uint64) error
}

func (m *mockRefreshStore) Create(ctx context.Context, userID uint64, raw string, ttl time.Duration) error {
	return m.createFunc(ctx, userID, raw, ttl)
}

func (m *mockRefreshStore) Validate(ctx context.Context, raw string) (uint64, error) {
	return m.validateFunc(ctx, raw)
}

func (m *mockRefreshStore) Revoke(ctx context.Context, raw string) error {
	return m.revokeFunc(ctx, raw)
}

func (m *mockRefreshStore) RevokeAllForUser(ctx context.Context, userID uint64) error {
	return m.revokeAllFunc(ctx, userID)
}

type mockUserRepo struct {
	createFunc        func(ctx context.Context, u *models.User) (*models.User, error)
	getByUsernameFunc func(ctx context.Context, username string) (*models.User, error)
	getByEmailFunc    func(ctx context.Context, email string) (*models.User, error)
	getByIDFunc       func(ctx context.Context, id uint64) (*models.User, error)
	recordLoginFunc   func(ctx context.Context, id uint64) error
}

func (m *mockUserRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	return m.createFunc(ctx, u)
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.getByUsernameFunc(ctx, username)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.getByEmailFunc(ctx, email)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockUserRepo) RecordLogin(ctx context.Context, id uint64) error {
	return m.recordLoginFunc(ctx, id)
}

func TestRegister(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	repo := &mockUserRepo{
		getByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return nil, ErrNotFound
		},
		createFunc: func(ctx context.Context, u *models.User) (*models.User, error) {
			return &models.User{ID: 1, Username: u.Username, Email: u.Email, Role: u.Role}, nil
		},
	}
	store := &mockRefreshStore{
		createFunc: func(ctx context.Context, userID uint64, raw string, ttl time.Duration) error {
			return nil
		},
	}

	svc := NewService(repo, store, tokens, 15)
	user, pair, err := svc.Register(context.Background(), "alice", "alice@test.com", "password123", models.RoleCustomer)

	require.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, models.RoleCustomer, user.Role)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, "Bearer", pair.TokenType)
	assert.Equal(t, int64(900), pair.ExpiresIn)
}

func TestRegisterDuplicate(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	repo := &mockUserRepo{
		getByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return &models.User{ID: 1, Username: username}, nil
		},
	}
	store := &mockRefreshStore{}

	svc := NewService(repo, store, tokens, 15)
	_, _, err := svc.Register(context.Background(), "alice", "alice@test.com", "password123", models.RoleCustomer)

	assert.ErrorIs(t, err, ErrDuplicate)
}

func TestLoginSuccess(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	repo := &mockUserRepo{
		getByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return &models.User{ID: 1, Username: username, Email: "a@b.com", Role: models.RoleCustomer, IsActive: true}, nil
		},
		recordLoginFunc: func(ctx context.Context, id uint64) error {
			return nil
		},
	}
	store := &mockRefreshStore{
		createFunc: func(ctx context.Context, userID uint64, raw string, ttl time.Duration) error {
			return nil
		},
	}

	svc := NewService(repo, store, tokens, 15)
	user, pair, err := svc.Login(context.Background(), "alice", "password123")

	require.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestLoginInvalidCredentials(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	repo := &mockUserRepo{
		getByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return nil, ErrNotFound
		},
		getByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, ErrNotFound
		},
	}
	store := &mockRefreshStore{}

	svc := NewService(repo, store, tokens, 15)
	_, _, err := svc.Login(context.Background(), "alice", "wrong")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestRefreshSuccess(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	repo := &mockUserRepo{
		getByIDFunc: func(ctx context.Context, id uint64) (*models.User, error) {
			return &models.User{ID: id, Username: "alice", Email: "a@b.com", Role: models.RoleCustomer, IsActive: true}, nil
		},
	}
	store := &mockRefreshStore{
		validateFunc: func(ctx context.Context, raw string) (uint64, error) {
			return 1, nil
		},
		revokeFunc: func(ctx context.Context, raw string) error {
			return nil
		},
		createFunc: func(ctx context.Context, userID uint64, raw string, ttl time.Duration) error {
			return nil
		},
	}

	svc := NewService(repo, store, tokens, 15)
	user, pair, err := svc.Refresh(context.Background(), "refresh-token")

	require.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestLogout(t *testing.T) {
	tokens := token.New("secret", 15, 168)
	store := &mockRefreshStore{
		revokeFunc: func(ctx context.Context, raw string) error {
			return nil
		},
	}

	svc := NewService(nil, store, tokens, 15)
	err := svc.Logout(context.Background(), "refresh-token")

	assert.NoError(t, err)
}
