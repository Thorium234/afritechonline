package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Env         string
	Port        string
	BaseURL     string
	JWTSecret   string
	AccessTTL   int // minutes
	RefreshTTL  int // hours

	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	DBMaxOpen   int
	DBMaxIdle   int
	DBMaxLifetimeMin int
}

// Load reads configuration from a .env file (if present) and environment.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:        getEnv("APP_ENV", "development"),
		Port:       getEnv("APP_PORT", "8080"),
		BaseURL:    getEnv("APP_BASE_URL", "http://localhost:8080"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		AccessTTL:  getEnvInt("JWT_ACCESS_TTL_MINS", 15),
		RefreshTTL: getEnvInt("JWT_REFRESH_TTL_HOURS", 168),

		DBHost:      getEnv("DB_HOST", "127.0.0.1"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBName:      getEnv("DB_NAME", "afritechonline"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBMaxOpen:   getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdle:   getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxLifetimeMin: getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 5),
	}

	// Validate JWT secret: must be set and not the default placeholder
	if err := validateJWTSecret(cfg.JWTSecret); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DSN returns the MySQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// validateJWTSecret ensures the JWT secret meets security requirements.
func validateJWTSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	
	placeholders := []string{
		"change-me-in-production",
		"change-me-in-production-use-strong-random-string",
		"change-me-in-production-use-strong-random-string-at-least-32-chars",
	}
	
	for _, p := range placeholders {
		if secret == p {
			return fmt.Errorf("JWT_SECRET cannot be a default placeholder. Generate a strong random secret using: openssl rand -base64 32")
		}
	}
	
	// Minimum 32 characters for security (256-bit when base64 encoded)
	if len(secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long (got %d). Generate using: openssl rand -base64 32", len(secret))
	}
	
	return nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
