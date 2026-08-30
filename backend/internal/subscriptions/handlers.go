package subscriptions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/pkg/response"
)

// Handler exposes subscription HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a subscription handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	CustomerID uint64 `json:"customer_id"`
	PackageID  uint64 `json:"package_id"`
}

// List returns subscriptions.
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
		"subscriptions": items,
		"pagination": gin.H{
			"page": page, "page_size": pageSize, "total": total,
		},
	})
}

// Get returns a single subscription.
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid subscription id")
		return
	}
	sub, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "subscription not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"subscription": sub})
}

// Create creates a new PENDING subscription for a customer + package selection.
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	fieldErrs := map[string]string{}
	if req.CustomerID == 0 {
		fieldErrs["customer_id"] = "customer_id is required"
	}
	if req.PackageID == 0 {
		fieldErrs["package_id"] = "package_id is required"
	}
	if len(fieldErrs) > 0 {
		response.Validation(c, fieldErrs)
		return
	}

	sub, err := h.service.Create(c.Request.Context(), req.CustomerID, req.PackageID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "customer or package not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"subscription": sub})
}
