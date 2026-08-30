package customers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/pkg/response"
	"github.com/Thorium234/afritechonline/backend/pkg/validator"
)

// Handler exposes customer HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a customer handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	FullName string  `json:"full_name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Status   string  `json:"status"`
	UserID   *uint64 `json:"user_id,omitempty"`
}

type updateRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Status   *string `json:"status"`
}

// List returns a paginated, optionally searchable list of customers.
func (h *Handler) List(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	customers, total, err := h.service.ListWithPagination(c.Request.Context(), search, page, pageSize)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{
		"customers": customers,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// Get returns a single customer.
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	customer, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "customer not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"customer": customer})
}

// Create registers a new customer.
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fieldErrs := map[string]string{}
	if req.FullName == "" {
		fieldErrs["full_name"] = "full_name is required"
	}
	if req.Username == "" {
		fieldErrs["username"] = "username is required"
	}
	if req.Phone == "" || !validator.IsValidPhone(req.Phone) {
		fieldErrs["phone"] = "a valid phone number is required"
	}
	if s := validateStatus(req.Status, true); s != "" {
		fieldErrs["status"] = s
	}
	if req.Email != "" && !validator.IsValidEmail(req.Email) {
		fieldErrs["email"] = "a valid email is required"
	}
	if len(fieldErrs) > 0 {
		response.Validation(c, fieldErrs)
		return
	}

	customer := &models.Customer{
		FullName: req.FullName,
		Phone:    validator.NormalizePhone(req.Phone),
		Email:    req.Email,
		Username: req.Username,
		Status:   req.Status,
		UserID:   req.UserID,
	}
	created, err := h.service.Create(c.Request.Context(), customer)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			response.Error(c, http.StatusConflict, "a customer with that phone or username already exists")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"customer": created})
}

// Update modifies an existing customer.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	existing, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "customer not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.Phone != nil {
		existing.Phone = validator.NormalizePhone(*req.Phone)
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.Username != nil {
		existing.Username = *req.Username
	}
	if req.Status != nil {
		if s := validateStatus(*req.Status, false); s != "" {
			response.Validation(c, map[string]string{"status": s})
			return
		}
		existing.Status = *req.Status
	}

	updated, err := h.service.Update(c.Request.Context(), existing)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			response.Error(c, http.StatusConflict, "a customer with that phone or username already exists")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"customer": updated})
}

// Delete removes a customer.
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "customer not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "customer deleted"})
}

var validStatuses = map[string]bool{
	"ACTIVE": true, "INACTIVE": true, "SUSPENDED": true,
}

// validateStatus returns an error string if the status is invalid, else "".
func validateStatus(s string, emptyAllowed bool) string {
	if s == "" {
		if emptyAllowed {
			return ""
		}
		return "status is required"
	}
	if !validStatuses[s] {
		return "status must be ACTIVE, INACTIVE or SUSPENDED"
	}
	return ""
}
