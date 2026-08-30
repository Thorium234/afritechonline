package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/internal/users"
	"github.com/Thorium234/afritechonline/backend/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned on a failed login.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrDuplicate is returned when a username/email is already registered.
var ErrDuplicate = errors.New("account already exists")

// TokenPair is returned after a successful login/refresh.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Service orchestrates authentication flows.
type Service struct {
	users      *users.Repository
	store      *RefreshTokenStore
	tokens     *token.Manager
	accessTTLMins int
}

// NewService creates an auth service.
func NewService(repo *users.Repository, store *RefreshTokenStore, tokens *token.Manager, accessTTLMins int) *Service {
	return &Service{users: repo, store: store, tokens: tokens, accessTTLMins: accessTTLMins}
}

// Register creates a new user account.
func (s *Service) Register(ctx context.Context, username, email, password string, role models.Role) (*models.User, TokenPair, error) {
	if _, err := s.users.GetByUsername(ctx, username); err == nil {
		return nil, TokenPair{}, ErrDuplicate
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, TokenPair{}, err
	}

	u := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}
	created, err := s.users.Create(ctx, u)
	if err != nil {
		return nil, TokenPair{}, fmt.Errorf("create user: %w", err)
	}

	pair, err := s.issue(ctx, created)
	if err != nil {
		return created, pair, err
	}
	return created, pair, nil
}

// Login authenticates a user by username/email and password.
func (s *Service) Login(ctx context.Context, identifier, password string) (*models.User, TokenPair, error) {
	u, err := s.users.GetByUsername(ctx, identifier)
	if err != nil {
		u, err = s.users.GetByEmail(ctx, identifier)
		if err != nil {
			return nil, TokenPair{}, ErrInvalidCredentials
		}
	}
	if !u.IsActive {
		return nil, TokenPair{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	if err := s.users.RecordLogin(ctx, u.ID); err != nil {
		// non-fatal
	}

	pair, err := s.issue(ctx, u)
	if err != nil {
		return nil, TokenPair{}, err
	}
	return u, pair, nil
}

// Refresh issues a new token pair from a valid refresh token.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.User, TokenPair, error) {
	userID, err := s.store.Validate(ctx, refreshToken)
	if err != nil {
		return nil, TokenPair{}, ErrInvalidCredentials
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, TokenPair{}, ErrInvalidCredentials
	}
	if !u.IsActive {
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	// Rotate: revoke old refresh token, issue a new one.
	if err := s.store.Revoke(ctx, refreshToken); err != nil {
		return nil, TokenPair{}, err
	}
	pair, err := s.issue(ctx, u)
	if err != nil {
		return nil, TokenPair{}, err
	}
	return u, pair, nil
}

// Logout revokes the given refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.store.Revoke(ctx, refreshToken)
}

func (s *Service) issue(ctx context.Context, u *models.User) (TokenPair, error) {
	access, err := s.tokens.IssueAccessToken(u.ID, string(u.Role), u.Username)
	if err != nil {
		return TokenPair{}, err
	}
	refresh, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	if err := s.store.Create(ctx, u.ID, refresh, s.tokens.RefreshTokenTTL()); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTTLMins * 60),
	}, nil
}
