-- Afritech Online — Initial schema
-- MySQL / MariaDB

CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL UNIQUE,
    email         VARCHAR(254) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          ENUM('SUPER_ADMIN','ADMIN','STAFF','CUSTOMER') NOT NULL DEFAULT 'CUSTOMER',
    is_active     TINYINT(1)   NOT NULL DEFAULT 1,
    last_login_at DATETIME     NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS customers (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NULL,
    full_name  VARCHAR(128) NOT NULL,
    phone      VARCHAR(20)  NOT NULL,
    email      VARCHAR(254) NOT NULL DEFAULT '',
    username   VARCHAR(64)  NOT NULL,
    status     ENUM('ACTIVE','INACTIVE','SUSPENDED') NOT NULL DEFAULT 'INACTIVE',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_customers_phone (phone),
    UNIQUE KEY uq_customers_username (username),
    CONSTRAINT fk_customers_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internet_packages (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name          VARCHAR(128) NOT NULL,
    description   VARCHAR(255) NOT NULL DEFAULT '',
    price         DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency      VARCHAR(3)   NOT NULL DEFAULT 'KES',
    duration_days INT          NOT NULL DEFAULT 30,
    download_mbps INT          NOT NULL DEFAULT 0,
    upload_mbps   INT          NOT NULL DEFAULT 0,
    data_limit_gb INT          NULL,
    is_active     TINYINT(1)   NOT NULL DEFAULT 1,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_packages_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS subscriptions (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT UNSIGNED NOT NULL,
    package_id  BIGINT UNSIGNED NOT NULL,
    start_date  DATETIME NOT NULL,
    expiry_date DATETIME NOT NULL,
    status      ENUM('PENDING','ACTIVE','EXPIRED','SUSPENDED','CANCELLED') NOT NULL DEFAULT 'PENDING',
    amount      DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency    VARCHAR(3) NOT NULL DEFAULT 'KES',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sub_status (status),
    INDEX idx_sub_expiry (expiry_date),
    INDEX idx_sub_customer (customer_id),
    CONSTRAINT fk_sub_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT fk_sub_package FOREIGN KEY (package_id) REFERENCES internet_packages(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS invoices (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_no      VARCHAR(32) NOT NULL UNIQUE,
    subscription_id BIGINT UNSIGNED NOT NULL,
    customer_id     BIGINT UNSIGNED NOT NULL,
    amount          DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency        VARCHAR(3) NOT NULL DEFAULT 'KES',
    status          ENUM('PENDING','PAID','OVERDUE','CANCELLED') NOT NULL DEFAULT 'PENDING',
    due_date        DATETIME NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_inv_customer (customer_id),
    CONSTRAINT fk_inv_sub FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE,
    CONSTRAINT fk_inv_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS payments (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_id  BIGINT UNSIGNED NOT NULL,
    customer_id BIGINT UNSIGNED NOT NULL,
    amount      DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency    VARCHAR(3) NOT NULL DEFAULT 'KES',
    method      ENUM('MANUAL','MPESA','CARD','OTHER') NOT NULL DEFAULT 'MANUAL',
    reference   VARCHAR(64) NOT NULL DEFAULT '',
    status      ENUM('PENDING','COMPLETED','FAILED','CANCELLED') NOT NULL DEFAULT 'PENDING',
    paid_at     DATETIME NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_payment_reference (reference, method),
    INDEX idx_pay_invoice (invoice_id),
    INDEX idx_pay_status (status),
    CONSTRAINT fk_pay_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    CONSTRAINT fk_pay_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS routers (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(128) NOT NULL,
    host         VARCHAR(128) NOT NULL,
    api_port     INT NOT NULL DEFAULT 8728,
    username     VARCHAR(64) NOT NULL DEFAULT 'admin',
    password_enc VARBINARY(255) DEFAULT NULL,
    location     VARCHAR(128) NOT NULL DEFAULT '',
    status       ENUM('OFFLINE','ONLINE','UNKNOWN') NOT NULL DEFAULT 'UNKNOWN',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  DATETIME NOT NULL,
    revoked_at  DATETIME NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_refresh_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS audit_logs (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    actor_id   BIGINT UNSIGNED NULL,
    actor_type VARCHAR(16) NOT NULL DEFAULT 'USER',
    action     VARCHAR(64) NOT NULL,
    entity     VARCHAR(64) NOT NULL,
    entity_id  VARCHAR(64) NOT NULL DEFAULT '',
    metadata   TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audit_entity (entity, entity_id),
    INDEX idx_audit_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
