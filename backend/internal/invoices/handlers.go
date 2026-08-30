package invoices

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/pkg/response"
)

// Handler exposes invoice HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an invoice handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	SubscriptionID uint64 `json:"subscription_id"`
}

// List returns invoices.
func (h *Handler) List(c *gin.Context) {
	customerID, _ := strconv.ParseUint(c.DefaultQuery("customer_id", "0"), 10, 64)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := h.service.List(c.Request.Context(), customerID, status, page, pageSize)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{
		"invoices": items,
		"pagination": gin.H{
			"page": page, "page_size": pageSize, "total": total,
		},
	})
}

// Get returns a single invoice.
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid invoice id")
		return
	}
	inv, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "invoice not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"invoice": inv})
}

// Create generates an invoice for a subscription.
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.SubscriptionID == 0 {
		response.Validation(c, map[string]string{"subscription_id": "subscription_id is required"})
		return
	}
	inv, err := h.service.CreateForSubscription(c.Request.Context(), req.SubscriptionID)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"invoice": inv})
}
