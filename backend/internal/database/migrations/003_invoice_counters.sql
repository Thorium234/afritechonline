-- Invoice sequence counter
CREATE TABLE IF NOT EXISTS invoice_counters (
    prefix VARCHAR(16) NOT NULL PRIMARY KEY,
    next_val BIGINT UNSIGNED NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO invoice_counters (prefix, next_val) VALUES ('daily', 1)
	ON DUPLICATE KEY UPDATE prefix = prefix;
