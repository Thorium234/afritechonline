package mikrotik

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// Handler exposes router HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a router handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List returns all routers.
func (h *Handler) List(c *gin.Context) {
	routers, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"routers": routers}})
}

// Get returns a single router.
func (h *Handler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid router id", "status": http.StatusBadRequest}})
		return
	}
	router, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "router not found", "status": http.StatusNotFound}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"router": router}})
}

// Create registers a new router.
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		APIPort  int    `json:"api_port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Location string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid request body", "status": http.StatusBadRequest}})
		return
	}

	router := &models.Router{
		Name:     req.Name,
		Host:     req.Host,
		APIPort:  req.APIPort,
		Username: req.Username,
		PasswordEnc: req.Password,
		Location: req.Location,
		Status:   models.RouterStatusUnknown,
	}
	created, err := h.service.Register(c.Request.Context(), router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"router": created}})
}

// Update modifies a router.
func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid router id", "status": http.StatusBadRequest}})
		return
	}
	router, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "router not found", "status": http.StatusNotFound}})
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Host     *string `json:"host"`
		APIPort  *int    `json:"api_port"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Location *string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid request body", "status": http.StatusBadRequest}})
		return
	}
	if req.Name != nil {
		router.Name = *req.Name
	}
	if req.Host != nil {
		router.Host = *req.Host
	}
	if req.APIPort != nil {
		router.APIPort = *req.APIPort
	}
	if req.Username != nil {
		router.Username = *req.Username
	}
	if req.Password != nil {
		router.PasswordEnc = *req.Password
	}
	if req.Location != nil {
		router.Location = *req.Location
	}

	updated, err := h.service.Update(c.Request.Context(), router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"router": updated}})
}

// Delete removes a router.
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid router id", "status": http.StatusBadRequest}})
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "router deleted"}})
}

// TestConnection checks router reachability.
func (h *Handler) TestConnection(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid router id", "status": http.StatusBadRequest}})
		return
	}
	router, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "router not found", "status": http.StatusNotFound}})
		return
	}

	identity, version, err := h.service.TestConnection(c.Request.Context(), router)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"connected": false, "error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"connected": true, "identity": identity, "version": version}})
}

// Status returns router status.
func (h *Handler) Status(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid router id", "status": http.StatusBadRequest}})
		return
	}
	router, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "router not found", "status": http.StatusNotFound}})
		return
	}

	identity, version, err := h.service.TestConnection(c.Request.Context(), router)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": router.Status, "connected": false}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": models.RouterStatusOnline, "connected": true, "identity": identity, "version": version}})
}

func parseID(c *gin.Context) (uint64, error) {
	s := c.Param("id")
	var id uint64
	if _, err := fmt.Sscanf(s, "%d", &id); err != nil {
		return 0, err
	}
	return id, nil
}
