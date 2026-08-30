package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSON writes a standard success envelope.
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

// Created writes a 201 Created response.
func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, data)
}

// Error writes a standard error envelope.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"status":  status,
			"message": message,
		},
	})
}

// Err writes a standard error envelope from an error value.
func Err(c *gin.Context, status int, err error) {
	Error(c, status, err.Error())
}

// Validation writes a 422 response with field errors.
func Validation(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"error": gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": "validation failed",
			"fields":  errors,
		},
	})
}
