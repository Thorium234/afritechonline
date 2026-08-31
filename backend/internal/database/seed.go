package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Seed inserts initial reference data and a default admin if they don't exist.
func Seed(db *sql.DB) error {
	seedPackages(db)

	// Create a default super admin if no admin exists.
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = ?`, models.RoleSuperAdmin).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		password := getEnv("SEED_ADMIN_PASSWORD", "")
		if password == "" {
			fmt.Println("SEED_ADMIN_PASSWORD not set; skipping default admin seed")
			return nil
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec(`INSERT INTO users (username, email, password_hash, role)
			VALUES (?, ?, ?, ?)`,
			"admin", "admin@afritech.online", string(hash), models.RoleSuperAdmin)
		if err != nil {
			return fmt.Errorf("seed admin: %w", err)
		}
		fmt.Println("seeded default admin: admin / [password from SEED_ADMIN_PASSWORD]")
	}
	return nil
}

func seedPackages(db *sql.DB) {
	defs := []struct {
		name, desc, currency string
		price                float64
		duration             int
		down, up             int
		limit                *int
	}{
		{"Home Basic", "10 Mbps, 30 days", "KES", 1000, 30, 10, 5, nil},
		{"Home Standard", "20 Mbps, 30 days", "KES", 2000, 30, 20, 10, nil},
		{"Home Premium", "40 Mbps, 30 days", "KES", 3500, 30, 40, 20, nil},
		{"Daily Pass", "5 Mbps, 24 hours", "KES", 100, 1, 5, 5, nil},
	}

	for _, d := range defs {
		_, err := db.Exec(`INSERT INTO internet_packages
			(name, description, price, currency, duration_days, download_mbps, upload_mbps, data_limit_gb, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)
			ON DUPLICATE KEY UPDATE name = VALUES(name)`,
			d.name, d.desc, d.price, d.currency, d.duration, d.down, d.up, d.limit)
		if err != nil {
			fmt.Println("seed package error:", err)
		}
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
