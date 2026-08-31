# Database Schema

## Tables

### users
Platform user accounts (admins, staff, customers).

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| username | VARCHAR(64) | Unique username |
| email | VARCHAR(254) | Unique email |
| password_hash | VARCHAR(255) | Bcrypt hash |
| role | ENUM | SUPER_ADMIN, ADMIN, STAFF, CUSTOMER |
| is_active | TINYINT | Account status |
| last_login_at | DATETIME | Last login timestamp |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### customers
ISP subscribers.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| user_id | BIGINT UNSIGNED | Optional linked user account |
| full_name | VARCHAR(128) | Customer full name |
| phone | VARCHAR(20) | Unique phone number |
| email | VARCHAR(254) | Email address |
| username | VARCHAR(64) | Unique username |
| status | ENUM | ACTIVE, INACTIVE, SUSPENDED |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### internet_packages
Internet service plans.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| name | VARCHAR(128) | Package name |
| description | VARCHAR(255) | Package description |
| price | DECIMAL(12,2) | Package price |
| currency | VARCHAR(3) | Currency code (default: KES) |
| duration_days | INT | Subscription duration |
| download_mbps | INT | Download speed |
| upload_mbps | INT | Upload speed |
| data_limit_gb | INT | Optional data limit |
| is_active | TINYINT | Active status |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### subscriptions
Customer package subscriptions.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| customer_id | BIGINT UNSIGNED | Customer reference |
| package_id | BIGINT UNSIGNED | Package reference |
| start_date | DATETIME | Subscription start |
| expiry_date | DATETIME | Subscription expiry |
| status | ENUM | PENDING, ACTIVE, EXPIRED, SUSPENDED, CANCELLED |
| amount | DECIMAL(12,2) | Subscription amount |
| currency | VARCHAR(3) | Currency code |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### invoices
Billing documents.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| invoice_no | VARCHAR(32) | Unique invoice number |
| subscription_id | BIGINT UNSIGNED | Subscription reference |
| customer_id | BIGINT UNSIGNED | Customer reference |
| amount | DECIMAL(12,2) | Invoice amount |
| currency | VARCHAR(3) | Currency code |
| status | ENUM | PENDING, PAID, OVERDUE, CANCELLED |
| due_date | DATETIME | Payment due date |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### payments
Payment records.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| invoice_id | BIGINT UNSIGNED | Invoice reference |
| customer_id | BIGINT UNSIGNED | Customer reference |
| amount | DECIMAL(12,2) | Payment amount |
| currency | VARCHAR(3) | Currency code |
| method | ENUM | MANUAL, MPESA, CARD, OTHER |
| reference | VARCHAR(64) | External payment reference |
| status | ENUM | PENDING, COMPLETED, FAILED, CANCELLED |
| paid_at | DATETIME | Payment timestamp |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### routers
MikroTik network devices.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| name | VARCHAR(128) | Router name |
| host | VARCHAR(128) | IP address or hostname |
| api_port | INT | API port (default: 8728) |
| username | VARCHAR(64) | Router username |
| password_enc | VARBINARY(255) | Encrypted password |
| location | VARCHAR(128) | Physical location |
| status | ENUM | OFFLINE, ONLINE, UNKNOWN |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### refresh_tokens
JWT refresh tokens.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| user_id | BIGINT UNSIGNED | User reference |
| token_hash | VARCHAR(255) | SHA-256 hash of token |
| expires_at | DATETIME | Token expiration |
| revoked_at | DATETIME | Revocation timestamp |
| created_at | DATETIME | Creation timestamp |

### audit_logs
Administrative action logs.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| actor_id | BIGINT UNSIGNED | User who performed action |
| actor_type | VARCHAR(16) | Type of actor (USER, SYSTEM) |
| action | VARCHAR(64) | Action performed |
| entity | VARCHAR(64) | Entity type affected |
| entity_id | VARCHAR(64) | Entity identifier |
| metadata | TEXT | Additional data |
| created_at | DATETIME | Timestamp |

### radius_users
FreeRADIUS user accounts.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT UNSIGNED | Primary key |
| username | VARCHAR(64) | RADIUS username |
| password | VARCHAR(255) | RADIUS password |
| profile | VARCHAR(128) | Speed profile |
| speed | VARCHAR(32) | Speed limit |
| expiry_date | DATETIME | Account expiry |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

## Migrations

Migrations are applied automatically on startup. They are tracked in the `schema_migrations` table.

Current migrations:
1. `001_init.sql` - Initial schema
2. `002_invoice_unique_subscription.sql` - Enforce one invoice per subscription
3. `003_invoice_counters.sql` - Invoice sequence counter
4. `004_radius_users.sql` - FreeRADIUS users table
