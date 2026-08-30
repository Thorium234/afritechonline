package models

import "time"

// Role defines a user's authorization level.
type Role string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleAdmin      Role = "ADMIN"
	RoleStaff      Role = "STAFF"
	RoleCustomer   Role = "CUSTOMER"
)

// User represents a user of the platform (admin, staff or customer login).
type User struct {
	ID           uint64     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Customer represents an internet subscriber.
type Customer struct {
	ID        uint64     `json:"id"`
	UserID    *uint64    `json:"user_id,omitempty"`
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Status    string     `json:"status"` // ACTIVE, INACTIVE, SUSPENDED
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// InternetPackage represents a purchasable internet plan.
type InternetPackage struct {
	ID               uint64     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Price            float64    `json:"price"`
	Currency         string     `json:"currency"`
	DurationDays     int        `json:"duration_days"`
	DownloadMbps     int        `json:"download_mbps"`
	UploadMbps       int        `json:"upload_mbps"`
	DataLimitGB      *int       `json:"data_limit_gb,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Subscription represents a customer's purchase of an internet package.
type Subscription struct {
	ID        uint64    `json:"id"`
	CustomerID uint64   `json:"customer_id"`
	PackageID  uint64   `json:"package_id"`
	StartDate  time.Time `json:"start_date"`
	ExpiryDate time.Time `json:"expiry_date"`
	Status     string   `json:"status"` // PENDING, ACTIVE, EXPIRED, SUSPENDED, CANCELLED
	Amount     float64  `json:"amount"`
	Currency   string   `json:"currency"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Invoice represents a billing document for a subscription.
type Invoice struct {
	ID             uint64    `json:"id"`
	InvoiceNo      string    `json:"invoice_no"`
	SubscriptionID uint64    `json:"subscription_id"`
	CustomerID     uint64    `json:"customer_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"` // PENDING, PAID, OVERDUE, CANCELLED
	DueDate        time.Time `json:"due_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Payment represents a payment against an invoice.
type Payment struct {
	ID         uint64    `json:"id"`
	InvoiceID  uint64    `json:"invoice_id"`
	CustomerID uint64    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	Method     string    `json:"method"` // MANUAL, MPESA, CARD, OTHER
	Reference  string    `json:"reference"`
	Status     string    `json:"status"` // PENDING, COMPLETED, FAILED, CANCELLED
	PaidAt     *time.Time `json:"paid_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Router represents a MikroTik network device managed by the platform.
type Router struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	APIPort     int       `json:"api_port"`
	Username    string    `json:"username"`
	PasswordEnc string    `json:"-"`
	Location    string    `json:"location"`
	Status      string    `json:"status"` // OFFLINE, ONLINE, UNKNOWN
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuditLog records sensitive administrative actions.
type AuditLog struct {
	ID        uint64    `json:"id"`
	ActorID   *uint64   `json:"actor_id,omitempty"`
	ActorType string    `json:"actor_type"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
