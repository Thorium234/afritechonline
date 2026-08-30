package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/internal/users"
	"github.com/Thorium234/afritechonline/backend/pkg/response"
	"github.com/Thorium234/afritechonline/backend/pkg/token"
)

// AuthMiddleware authenticates requests via a Bearer access token.
// It loads the full user into the request context.
type AuthMiddleware struct {
	tokens *token.Manager
	users  *users.Repository
}

// NewAuthMiddleware creates the auth middleware.
func NewAuthMiddleware(tokens *token.Manager, users *users.Repository) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens, users: users}
}

// RequireAuth validates the access token and loads the user.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims, err := m.tokens.ParseAccessToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		user, err := m.users.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "user no longer exists")
			c.Abort()
			return
		}
		if !user.IsActive {
			response.Error(c, http.StatusUnauthorized, "account is disabled")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("role", user.Role)
		c.Next()
	}
}

// RequireRole restricts access to one or more roles.
func RequireRole(roles ...models.Role) gin.HandlerFunc {
	allowed := map[models.Role]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			response.Error(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}
		if !allowed[role.(models.Role)] {
			response.Error(c, http.StatusForbidden, "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
