package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

// Limit returns a middleware that rate limits by client IP.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
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

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	var recent []time.Time
	for _, t := range rl.clients[key] {
		if now.Sub(t) < rl.window {
			recent = append(recent, t)
		}
	}
	if len(recent) >= rl.limit {
		return false
	}
	rl.clients[key] = append(recent, now)
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.clients {
			var recent []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					recent = append(recent, t)
				}
			}
			if len(recent) == 0 {
				delete(rl.clients, key)
			} else {
				rl.clients[key] = recent
			}
		}
		rl.mu.Unlock()
	}
}
