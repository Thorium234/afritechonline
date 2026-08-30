package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Thorium234/afritechonline/backend/config"
	"github.com/Thorium234/afritechonline/backend/internal/auth"
	"github.com/Thorium234/afritechonline/backend/internal/customers"
	"github.com/Thorium234/afritechonline/backend/internal/invoices"
	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/internal/packages"
	"github.com/Thorium234/afritechonline/backend/internal/payments"
	"github.com/Thorium234/afritechonline/backend/internal/subscriptions"
	"github.com/Thorium234/afritechonline/backend/internal/users"
	"github.com/Thorium234/afritechonline/backend/middleware"
	"github.com/Thorium234/afritechonline/backend/pkg/token"
)

// Setup configures the Gin engine and all routes.
func Setup(db *sql.DB, cfg *config.Config, log zerolog.Logger) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(log))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ---- Repositories ----
	userRepo := users.New(db)
	customerRepo := customers.New(db)
	packageRepo := packages.New(db)
	subscriptionRepo := subscriptions.New(db)
	invoiceRepo := invoices.New(db)
	paymentRepo := payments.New(db)

	// ---- Services ----
	tokens := token.New(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	refreshStore := auth.NewRefreshTokenStore(db)
	authService := auth.NewService(userRepo, refreshStore, tokens, cfg.AccessTTL)

	customerService := customers.NewService(customerRepo)
	packageService := packages.NewService(packageRepo)
	subscriptionService := subscriptions.NewService(subscriptionRepo, packageService, customerService)
	invoiceService := invoices.NewService(invoiceRepo, subscriptionService)
	paymentService := payments.NewService(paymentRepo, invoiceService, subscriptionService)

	// ---- Handlers ----
	authHandler := auth.NewHandler(authService)
	customerHandler := customers.NewHandler(customerService)
	packageHandler := packages.NewHandler(packageService)
	subscriptionHandler := subscriptions.NewHandler(subscriptionService)
	invoiceHandler := invoices.NewHandler(invoiceService)
	paymentHandler := payments.NewHandler(paymentService)

	// ---- Middleware ----
	authMW := middleware.NewAuthMiddleware(tokens, userRepo)
	staffOnly := middleware.RequireRole(models.RoleAdmin, models.RoleSuperAdmin, models.RoleStaff)
	adminOnly := middleware.RequireRole(models.RoleAdmin, models.RoleSuperAdmin)

	v1 := r.Group("/api/v1")

	// Auth (public)
	a := v1.Group("/auth")
	{
		a.POST("/register", authHandler.Register)
		a.POST("/login", authHandler.Login)
		a.POST("/refresh", authHandler.Refresh)
		a.POST("/logout", authHandler.Logout)
		a.GET("/me", authMW.RequireAuth(), authHandler.Me)
	}

	// Protected
	protected := v1.Group("")
	protected.Use(authMW.RequireAuth())

	// Customers
	customers := protected.Group("/customers")
	{
		customers.GET("", staffOnly, customerHandler.List)
		customers.POST("", staffOnly, customerHandler.Create)
		customers.GET("/:id", staffOnly, customerHandler.Get)
		customers.PUT("/:id", staffOnly, customerHandler.Update)
		customers.DELETE("/:id", adminOnly, customerHandler.Delete)
	}

	// Packages
	packages := protected.Group("/packages")
	{
		packages.GET("", packageHandler.List)
		packages.POST("", staffOnly, packageHandler.Create)
		packages.GET("/:id", packageHandler.Get)
		packages.PUT("/:id", staffOnly, packageHandler.Update)
		packages.DELETE("/:id", adminOnly, packageHandler.Delete)
	}

	// Subscriptions
	subscriptions := protected.Group("/subscriptions")
	{
		subscriptions.GET("", subscriptionHandler.List)
		subscriptions.POST("", staffOnly, subscriptionHandler.Create)
		subscriptions.GET("/:id", subscriptionHandler.Get)
	}

	// Invoices
	invoices := protected.Group("/invoices")
	{
		invoices.GET("", invoiceHandler.List)
		invoices.POST("", staffOnly, invoiceHandler.Create)
		invoices.GET("/:id", invoiceHandler.Get)
	}

	// Payments
	payments := protected.Group("/payments")
	{
		payments.GET("", paymentHandler.List)
		payments.POST("", staffOnly, paymentHandler.Create)
		payments.GET("/:id", paymentHandler.Get)
		payments.POST("/:id/complete", staffOnly, paymentHandler.Complete)
		payments.POST(