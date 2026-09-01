package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket rate limiting per IP address.
type RateLimiter struct {
	requests map[string]*tokenBucket
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type tokenBucket struct {
	tokens    float64
	lastSeen  time.Time
	refillRate float64
}

// NewRateLimiter creates a new rate limiter with the given requests per window.
func NewRateLimiter(requestsPerWindow int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*tokenBucket),
		limit:    requestsPerWindow,
		window:   window,
	}
	
	// Cleanup stale entries every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			now := time.Now()
			for ip, bucket := range rl.requests {
				if now.Sub(bucket.lastSeen) > 24*time.Hour {
					delete(rl.requests, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	
	return rl
}

// Allow checks if the request from the given IP should be allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	bucket, exists := rl.requests[ip]
	now := time.Now()
	
	if !exists {
		// New IP: initialize with full bucket
		rl.requests[ip] = &tokenBucket{
			tokens:     float64(rl.limit),
			lastSeen:   now,
			refillRate: float64(rl.limit) / rl.window.Seconds(),
		}
		rl.requests[ip].tokens-- // Use one token for this request
		return true
	}
	
	// Refill tokens based on time elapsed
	elapsed := now.Sub(bucket.lastSeen).Seconds()
	bucket.tokens = min(float64(rl.limit), bucket.tokens+elapsed*bucket.refillRate)
	bucket.lastSeen = now
	
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}
	
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware returns a Gin middleware that applies rate limiting.
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c.Request)
		
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"status":  http.StatusTooManyRequests,
					"message": "rate limit exceeded",
				},
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RateLimitEndpoint returns middleware for rate limiting specific endpoints.
func RateLimitEndpoint(rl *RateLimiter) gin.HandlerFunc {
	return RateLimitMiddleware(rl)
}

// getClientIP extracts the client IP from the request, respecting X-Forwarded-For header.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		// X-Forwarded-For might not have port
		return xff
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to remote address
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	
	return r.RemoteAddr
}
