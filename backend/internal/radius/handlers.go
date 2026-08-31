package radius

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler exposes radius HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a radius handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List returns RADIUS users.
func (h *Handler) List(c *gin.Context) {
	// Placeholder: list all RADIUS users
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"users": []string{}}})
}

// Get returns a single RADIUS user.
func (h *Handler) Get(c *gin.Context) {
	username := c.Param("username")
	user, err := h.service.GetUser(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "user not found", "status": http.StatusNotFound}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"user": user}})
}

// Create adds a new RADIUS user.
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Profile  string `json:"profile"`
		Speed    string `json:"speed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid request body", "status": http.StatusBadRequest}})
		return
	}
	// Implementation continues...
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"user": nil}})
}

// Update modifies a RADIUS user.
func (h *Handler) Update(c *gin.Context) {
	username := c.Param("username")
	var req struct {
		Password string `json:"password"`
		Profile  string `json:"profile"`
		Speed    string `json:"speed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid request body", "status": http.StatusBadRequest}})
		return
	}
	// Implementation continues...
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"user": nil}})
}

// Delete removes a RADIUS user.
func (h *Handler) Delete(c *gin.Context) {
	username := c.Param("username")
	if err := h.service.DeleteUser(c.Request.Context(), username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "user deleted"}})
}
