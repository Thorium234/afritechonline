# Implementation Guide - Using New Security Features

This guide shows how to integrate the newly implemented security features into your existing code.

---

## 1. Integrating Rate Limiting

### In `backend/routes/routes.go`:

```go
import (
    "time"
    "github.com/Thorium234/afritechonline/backend/middleware"
)

func Setup(db *sql.DB, cfg *config.Config, log zerolog.Logger) *gin.Engine {
    // ... existing code ...
    
    // Create rate limiters for different endpoints
    authLimiter := middleware.NewRateLimiter(10, 15*time.Minute)    // 10 requests per 15 min
    apiLimiter := middleware.NewRateLimiter(100, time.Minute)      // 100 requests per minute
    
    // Apply to authentication endpoints
    authRoutes := r.Group("/api/v1/auth")
    authRoutes.Use(middleware.RateLimitMiddleware(authLimiter))
    {
        authRoutes.POST("/login", authHandler.Login)
        authRoutes.POST("/register", authHandler.Register)
        authRoutes.POST("/refresh", authHandler.Refresh)
    }
    
    // Apply to general API endpoints
    api := r.Group("/api/v1")
    api.Use(middleware.RateLimitMiddleware(apiLimiter))
    {
        // ... your routes ...
    }
    
    return r
}
```

---

## 2. Using Audit Logging

### In service files (e.g., `backend/internal/payments/service.go`):

```go
import (
    "github.com/Thorium234/afritechonline/backend/pkg/audit"
)

type Service struct {
    repo     *Repository
    auditLog *audit.Logger  // Add this field
}

func NewService(repo *Repository, auditLog *audit.Logger) *Service {
    return &Service{
        repo:     repo,
        auditLog: auditLog,
    }
}

// Example: Complete a payment
func (s *Service) CompletePayment(ctx context.Context, paymentID uint64, userID uint64) error {
    // ... your logic ...
    
    if err := s.repo.Update(ctx, payment); err != nil {
        // Log failed payment
        s.auditLog.LogPaymentEvent(ctx, 
            audit.EventPaymentFailed, 
            paymentID, userID, payment.Amount, "FAILURE",
            map[string]interface{}{"error": err.Error()})
        return err
    }
    
    // Log successful payment
    s.auditLog.LogPaymentEvent(ctx,
        audit.EventPaymentCompleted,
        paymentID, userID, payment.Amount, "SUCCESS",
        map[string]interface{}{
            "method": payment.Method,
            "reference": payment.Reference,
        })
    
    return nil
}
```

### Inject audit logger during initialization:

```go
func Setup(db *sql.DB, cfg *config.Config, log zerolog.Logger) *gin.Engine {
    // ... existing code ...
    
    auditLog := audit.New(log)
    
    // Pass to services
    paymentService := payments.NewService(paymentRepo, auditLog)
    subscriptionService := subscriptions.NewService(subscriptionRepo, auditLog)
    
    return r
}
```

---

## 3. Using Transaction Wrapper

### In `backend/internal/payments/service.go`:

```go
import (
    "github.com/Thorium234/afritechonline/backend/internal/database"
)

// CompletePayment atomically updates payment and subscription
func (s *Service) CompletePayment(ctx context.Context, paymentID uint64) error {
    return database.WithTx(s.db, ctx, func(tx *database.Tx) error {
        // Step 1: Update payment status
        payment := &models.Payment{ID: paymentID, Status: "COMPLETED"}
        if err := tx.QueryRow(
            "UPDATE payments SET status = $1 WHERE id = $2 RETURNING id",
            payment.Status, paymentID).Scan(&payment.ID); err != nil {
            return fmt.Errorf("update payment: %w", err)
        }
        
        // Step 2: Update subscription status
        if err := tx.Exec(
            "UPDATE subscriptions SET status = $1 WHERE id = $2",
            "ACTIVE", payment.SubscriptionID); err != nil {
            return fmt.Errorf("update subscription: %w", err)
        }
        
        // Step 3: Update invoice status
        if err := tx.Exec(
            "UPDATE invoices SET status = $1 WHERE id = $2",
            "PAID", payment.InvoiceID); err != nil {
            return fmt.Errorf("update invoice: %w", err)
        }
        
        return nil // Transaction auto-commits
    })
}
```

---

## 4. Using Password Validation

### In `backend/internal/auth/handlers.go`:

```go
import (
    "github.com/Thorium234/afritechonline/backend/pkg/validator"
)

func (h *Handler) Register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "invalid request body")
        return
    }
    
    fieldErrs := map[string]string{}
    
    // Validate email
    if req.Email == "" || !validator.HasValidEmail(req.Email) {
        fieldErrs["email"] = "a valid email is required"
    }
    
    // Validate password with NIST requirements
    passwordErrs := validator.ValidatePasswordDefault(req.Password)
    if len(passwordErrs) > 0 {
        fieldErrs["password"] = passwordErrs[0] // Return first error
    }
    
    if len(fieldErrs) > 0 {
        response.Validation(c, fieldErrs)
        return
    }
    
    // ... rest of registration logic ...
}
```

---

## 5. Enhanced Security Middleware

### Already integrated in routes.go (no changes needed):

The following are already applied in the middleware stack:
- `middleware.SecurityHeaders()` - Adds security headers
- `middleware.CORSMiddleware()` - Controls cross-origin access
- `middleware.ContentTypeValidation()` - Validates request types
- `middleware.RequestTimeout()` - Prevents slow clients

New headers now include:
- `Content-Security-Policy: default-src 'self'...`
- `Strict-Transport-Security` (commented for production)
- `Permissions-Policy`
- `X-Content-Type-Options: nosniff`

---

## 6. Testing the Implementations

### Test Rate Limiting

```bash
# Test with Apache Bench (simulate 20 concurrent requests)
ab -n 20 -c 20 http://localhost:8080/api/v1/auth/login

# Should see: 429 Too Many Requests after limit exceeded
```

### Test Audit Logging

Check logs for audit events:
```bash
docker logs afritechonline_backend | grep "event_type"
```

### Test Password Validation

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "weak"
  }'

# Response: 422 with validation error about password length
```

### Test Transactions

Enable transaction logging in database operations to verify atomicity.

---

## 7. Environment Setup

### Generate secure JWT secret:

```bash
# Linux/Mac
openssl rand -base64 32

# Windows PowerShell
[Convert]::ToBase64String([System.Security.Cryptography.RandomNumberGenerator]::GetBytes(32))
```

Copy the output to your `.env`:
```
JWT_SECRET=your-generated-secret-here
```

---

## 8. Deployment Checklist

Before going to production:

- [ ] Update JWT_SECRET in `.env` with generated value
- [ ] Set `APP_ENV=production`
- [ ] Uncomment HSTS header in `middleware/security.go`
- [ ] Enable HTTPS on reverse proxy (Nginx)
- [ ] Test rate limiting doesn't block legitimate users
- [ ] Verify audit logs are collected
- [ ] Monitor error rates during gradual rollout
- [ ] Setup alerts for repeated failed login attempts

---

## 9. Configuration Examples

### Strict Password Requirements (14 characters)

```go
errors := validator.ValidatePassword(password, validator.DefaultPasswordRequirements())
// Requires: 14+ chars, uppercase, lowercase, numbers
```

### Relaxed Password Requirements (backward compatibility)

```go
errors := validator.ValidatePassword(password, validator.LegacyPasswordRequirements())
// Requires: 8+ chars, uppercase, lowercase, numbers
```

### Custom Requirements

```go
custom := validator.PasswordRequirements{
    MinLength:        20,
    RequireUppercase: true,
    RequireLowercase: true,
    RequireNumbers:   true,
    RequireSpecial:   true,
    ProhibitCommon:   true,
}
errors := validator.ValidatePassword(password, custom)
```

---

## Support & References

- **Audit Logger**: Check `backend/pkg/audit/audit.go` for all event types
- **Rate Limiting**: Check `backend/middleware/ratelimit.go` for configuration
- **Password Validation**: Check `backend/pkg/validator/password.go` for rules
- **Transactions**: Check `backend/internal/database/transaction.go` for examples

For issues or questions, refer to `SECURITY_FIXES_SUMMARY.md`.
