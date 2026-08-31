package contextutil

import (
	"context"
	"time"
)

// DefaultTimeout is used when no specific timeout is provided.
const DefaultTimeout = 10 * time.Second

// WithTimeout returns a child context with the given timeout.
// If the parent context is already cancelled, it returns the cancelled context.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// DBTimeout returns a context suitable for database operations.
func DBTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(parent, 5*time.Second)
}
