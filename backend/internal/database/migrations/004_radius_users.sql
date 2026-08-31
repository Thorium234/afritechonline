-- FreeRADIUS users table
CREATE TABLE IF NOT EXISTS radius_users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username    VARCHAR(64) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    profile     VARCHAR(128) NOT NULL DEFAULT 'default',
    speed       VARCHAR(32) NOT NULL DEFAULT '10M/5M',
    expiry_date DATETIME NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_radius_username (username),
    INDEX idx_radius_expiry (expiry_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
