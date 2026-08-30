package packages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/pkg/response"
)

// Handler exposes internet package HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a package handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	DurationDays int     `json:"duration_days"`
	DownloadMbps int     `json:"download_mbps"`
	UploadMbps   int     `json:"upload_mbps"`
	DataLimitGB  *int    `json:"data_limit_gb"`
	IsActive     *bool   `json:"is_active"`
}

// List returns packages.
func (h *Handler) List(c *gin.Context) {
	activeOnlyParam := c.DefaultQuery("active", "")
	activeOnly := false
	if activeOnlyParam == "true" || activeOnlyParam == "1" {
		activeOnly = true
	}
	items, err := h.service.List(c.Request.Context(), activeOnly)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"packages": items})
}

// Get returns a single package.
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid package id")
		return
	}
	p, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "package not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"package": p})
}

// Create adds a new package.
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	fieldErrs := map[string]string{}
	if req.Name == "" {
		fieldErrs["name"] = "name is required"
	}
	if req.Price <= 0 {
		fieldErrs["price"] = "price must be greater than zero"
	}
	if req.DurationDays <= 0 {
		fieldErrs["duration_days"] = "duration_days must be greater than zero"
	}
	if len(fieldErrs) > 0 {
		response.Validation(c, fieldErrs)
		return
	}

	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	p := &models.InternetPackage{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Currency:     orDefault(req.Currency, "KES"),
		DurationDays: req.DurationDays,
		DownloadMbps: req.DownloadMbps,
		UploadMbps:   req.UploadMbps,
		DataLimitGB:  req.DataLimitGB,
		IsActive:     active,
	}
	created, err := h.service.Create(c.Request.Context(), p)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"package": created})
}

// Update modifies a package.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid package id")
		return
	}
	p, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "package not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.Currency != "" {
		p.Currency = req.Currency
	}
	if req.DurationDays > 0 {
		p.DurationDays = req.DurationDays
	}
	p.DownloadMbps = req.DownloadMbps
	p.UploadMbps = req.UploadMbps
	p.DataLimitGB = req.DataLimitGB
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	updated, err := h.service.Update(c.Request.Context(), p)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "package not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"package": updated})
}

// Delete removes a package.
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid package id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(c, http.StatusNotFound, "package not found")
			return
		}
		response.Err(c, http.StatusInternalServerError, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "package deleted"})
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
