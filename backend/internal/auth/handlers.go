package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/pkg/response"
	"github.com/Thorium234/afritechonline/backend/pkg/validator"
)

// Handler exposes authentication HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Register godoc placeholder.
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fieldErrs := map[string]string{}
	if req.Username == "" {
		fieldErrs["username"] = "username is required"
	}
	if req.Password == "" || len(req.Password) < 8 {
		fieldErrs["password"] = "password must be at least 8 characters"
	}
	if req.Email == "" || !validator.IsValidEmail(req.Email) {
		fieldErrs["email"] = "a valid email is required"
	}
	role := models.Role(req.Role)
	if role == "" {
		role = models.RoleCustomer
	}
	if !validRole(role) {
		fieldErrs["role"] = "role must be one of SUPER_ADMIN, ADMIN, STAFF, CUSTOMER"
	}
	if len(fieldErrs) > 0 {
		response.Validation(c, fieldErrs)
		return
	}

	user, pair, err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Password, role)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			response.Error(c, http.StatusConflict, "an account with that username or email already exists")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"user": user, "tokens": pair})
}

// Login authenticates and issues tokens.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Identifier == "" || req.Password == "" {
		response.Error(c, http.StatusBadRequest, "identifier and password are required")
		return
	}

	user, pair, err := h.service.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"user": user, "tokens": pair})
}

// Refresh rotates a refresh token.
func (h *Handler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, "refresh_token is required")
		return
	}
	user, pair, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"user": user, "tokens": pair})
}

// Logout revokes a refresh token.
func (h *Handler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, "refresh_token is required")
		return
	}
	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "logged out"})
}

// Me returns the currently authenticated user.
func (h *Handler) Me(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"user": user})
}

func validRole(r models.Role) bool {
	switch r {
	case models.RoleSuperAdmin, models.RoleAdmin, models.RoleStaff, models.RoleCustomer:
		return true
	}
	return false
}
