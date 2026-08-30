package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Thorium234/afritechonline/backend/config"
	"github.com/Thorium234/afritechonline/backend/internal/database"
	"github.com/Thorium234/afritechonline/backend/pkg/logger"
	"github.com/Thorium234/afritechonline/backend/routes"
)

func main() {
	log := logger.New("development")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}
	log = logger.New(cfg.Env)

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()
	log.Info().Msg("Database connected")

	if err := database.Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("database migration failed")
	}
	log.Info().Msg("Database migrations applied")

	if err := database.Seed(db); err != nil {
		log.Fatal().Err(err).Msg("database seeding failed")
	}

	router := routes.Setup(db, cfg, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("Server running on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}
	log.Info().Msg("Server exited gracefully")
}
