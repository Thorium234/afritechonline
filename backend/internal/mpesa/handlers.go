package mpesa

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler exposes M-Pesa HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an M-Pesa handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// STKPush initiates an STK Push request.
func (h *Handler) STKPush(c *gin.Context) {
	var req struct {
		InvoiceID  uint64  `json:"invoice_id"`
		PhoneNumber string `json:"phone_number"`
		Amount     float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid request body", "status": http.StatusBadRequest}})
		return
	}

	resp, err := h.service.STKPush(c.Request.Context(), req.InvoiceID, req.PhoneNumber, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Callback handles M-Pesa payment callbacks.
func (h *Handler) Callback(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "invalid body", "status": http.StatusBadRequest}})
		return
	}

	signature := c.GetHeader("x-mpesa-signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "missing signature", "status": http.StatusBadRequest}})
		return
	}

	if err := h.service.ProcessCallback(c.Request.Context(), body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "status": http.StatusInternalServerError}})
		return
	}

	c.Status(http.StatusOK)
}
