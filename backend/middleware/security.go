package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds common security headers to responses.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Control referrer information
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Restrict browser features and APIs
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
		
		// Content Security Policy - restrict resource loading
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'")
		
		// HSTS: Force HTTPS for 1 year (production only)
		// Uncomment for production deployment
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing with strict validation.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		isAllowed := false
		
		// Validate origin against whitelist
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}
		
		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		// Only allow specific methods
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		
		// Only allow necessary headers
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With, Accept")
		
		// Allow preflight cache for 12 hours
		c.Header("Access-Control-Max-Age", "43200")
		
		// Expose only necessary headers
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-JSON-Response-Code")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RequestTimeout ensures requests don't run indefinitely.
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ContentTypeValidation ensures requests have appropriate Content-Type headers.
func ContentTypeValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate for requests with bodies
		if c.Request.ContentLength > 0 && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") {
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "" || (contentType != "application/json" && contentType != "application/x-www-form-urlencoded") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": gin.H{
						"status":  http.StatusUnsupportedMediaType,
						"message": "Content-Type must be application/json or application/x-www-form-urlencoded",
					},
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
