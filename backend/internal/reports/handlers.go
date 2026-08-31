package reports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler exposes reports HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a reports handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RevenueSummary returns revenue summary.
func (h *Handler) RevenueSummary(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	summaries, err := h.service.RevenueSummary(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"summaries": summaries}})
}

// CustomerStats returns customer statistics.
func (h *Handler) CustomerStats(c *gin.Context) {
	stats, err := h.service.CustomerStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"stats": stats}})
}

// ActiveRouters returns count of online routers.
func (h *Handler) ActiveRouters(c *gin.Context) {
	count, err := h.service.ActiveRouters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"active_routers": count}})
}
