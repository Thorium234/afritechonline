package payments

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/pkg/response"
)

// Handler exposes payment HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a payment handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	InvoiceID  uint64  `json:"invoice_id"`
	CustomerID uint64  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Method     string  `json:"method"`
	Reference  string  `json:"reference"`
}

// List returns payments.
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
		"payments": items,
		"pagination": gin.H{
			"page": page, "page_size": pageSize, "total": total,
		},
	})
}

// Get returns a single payment.
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payment id")
		return
	}
	p, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "payment not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"payment": p})
}

// Create records a payment.
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	fieldErrs := map[string]string{}
	if req.InvoiceID == 0 {
		fieldErrs["invoice_id"] = "invoice_id is required"
	}
	if req.CustomerID == 0 {
		fieldErrs["customer_id"] = "customer_id is required"
	}
	if req.Amount <= 0 {
		fieldErrs["amount"] = "amount must be greater than zero"
	}
	if len(fieldErrs) > 0 {
		response.Validation(c, fieldErrs)
		return
	}

	method := orDefault(req.Method, "MANUAL")
	if !validMethod(method) {
		response.Validation(c, map[string]string{"method": "method must be MANUAL, MPESA, CARD or OTHER"})
		return
	}

	p := &models.Payment{
		InvoiceID:  req.InvoiceID,
		CustomerID: req.CustomerID,
		Amount:     req.Amount,
		Currency:   orDefault(req.Currency, "KES"),
		Method:     method,
		Reference:  req.Reference,
		Status:     StatusPending,
	}
	created, err := h.service.Create(c.Request.Context(), p)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"payment": created})
}

// Complete confirms a payment and activates the subscription (manual verification).
func (h *Handler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payment id")
		return
	}
	p, err := h.service.Complete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "payment not found")
			return
		}
		if errors.Is(err, ErrInvalidState) {
			response.Error(c, http.StatusConflict, "payment cannot be completed in its current state")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"payment": p})
}

// Fail marks a payment as FAILED (PAYMENT failure never activates internet).
func (h *Handler) Fail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Err