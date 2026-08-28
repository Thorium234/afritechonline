# Afritech Online Technology Specification

## 1. Overview

Afritech Online is an ISP management and network automation platform designed to manage customers, internet packages, subscriptions, billing, payments, authentication, and MikroTik network infrastructure.

The project will use:

* **Go** for the backend
* **Next.js** for the frontend
* **MySQL** for development
* **FreeRADIUS** for network authentication and accounting
* **MikroTik RouterOS** for network access and infrastructure management
* **Docker** for development and deployment
* **GitHub** for source control and CI/CD

The architecture should prioritize simplicity, reliability, security, and the ability to support multiple routers and ISP locations in the future.

---

# 2. Technology Stack

| Layer                  | Technology                          | Purpose                            |
| ---------------------- | ----------------------------------- | ---------------------------------- |
| Frontend               | Next.js                             | Web application and user interface |
| Frontend Language      | TypeScript                          | Type-safe frontend development     |
| Backend                | Go                                  | REST API and business logic        |
| Backend Framework      | Gin                                 | HTTP routing and middleware        |
| Database               | MySQL                               | Development database               |
| ORM / Database Access  | GORM or SQL                         | Database operations                |
| Authentication         | JWT + refresh tokens                | API authentication                 |
| Network Authentication | FreeRADIUS                          | Customer authentication            |
| Network Device         | MikroTik RouterOS                   | Internet access management         |
| Payments               | M-Pesa                              | Customer payments                  |
| Containers             | Docker                              | Local development and deployment   |
| API Format             | REST + JSON                         | Frontend/backend communication     |
| Documentation          | OpenAPI                             | API documentation                  |
| Testing                | Go testing + frontend testing tools | Automated testing                  |
| Version Control        | Git                                 | Source control                     |
| CI/CD                  | GitHub Actions                      | Automated testing and builds       |

---

# 3. Backend

## 3.1 Language

The backend will be written in **Go**.

Go is selected because the system will communicate with external services and network infrastructure where reliability, concurrency, performance, and simple deployment are important.

The backend will be responsible for:

* Authentication
* Authorization
* Customer management
* Internet packages
* Subscriptions
* Billing
* Payments
* M-Pesa integration
* MikroTik integration
* FreeRADIUS integration
* Session management
* Notifications
* Reporting
* Administrative operations

---

## 3.2 Backend Framework

The backend will use **Gin** as the HTTP framework.

Gin will provide:

* HTTP routing
* Middleware
* Request handling
* JSON responses
* Authentication middleware
* Error handling
* API grouping

The backend should expose a REST API.

Example:

```text
/api/v1/auth
/api/v1/users
/api/v1/customers
/api/v1/packages
/api/v1/subscriptions
/api/v1/billing
/api/v1/payments
/api/v1/mikrotik
/api/v1/radius
/api/v1/sessions
/api/v1/reports
```

---

# 4. Backend Architecture

The backend should use a layered architecture.

```text
HTTP Request
     │
     ▼
┌───────────────┐
│   Handlers    │
└───────┬───────┘
        │
        ▼
┌───────────────┐
│   Services    │
└───────┬───────┘
        │
        ├──────────────┐
        ▼              ▼
┌───────────────┐ ┌───────────────┐
│ Repositories  │ │ External APIs │
└───────┬───────┘ └───────┬───────┘
        │                 │
        ▼                 ├── M-Pesa
     MySQL                ├── MikroTik
                          └── FreeRADIUS
```

The handler layer should not contain business logic.

For example, this is bad:

```text
HTTP Handler
    ↓
Validate payment
    ↓
Calculate subscription
    ↓
Update database
    ↓
Connect MikroTik
    ↓
Activate customer
```

Instead:

```text
HTTP Handler
    ↓
Payment Service
    ↓
Payment Repository
    ↓
Subscription Service
    ↓
Network Provisioning Service
```

This separation makes the system easier to test and maintain.

---

# 5. Recommended Backend Structure

```text
backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   ├── users/
│   ├── customers/
│   ├── packages/
│   ├── subscriptions/
│   ├── billing/
│   ├── payments/
│   ├── mikrotik/
│   ├── radius/
│   ├── sessions/
│   ├── notifications/
│   └── reports/
│
├── config/
│
├── database/
│   ├── migrations/
│   └── seeds/
│
├── middleware/
│
├── routes/
│
├── pkg/
│   ├── logger/
│   ├── response/
│   └── validator/
│
├── tests/
│
├── go.mod
├── go.sum
└── Dockerfile
```

The `internal` directory should contain application-specific business logic that should not be imported by unrelated external projects.

---

# 6. Frontend

## 6.1 Framework

The frontend will use **Next.js with TypeScript**.

Next.js will provide:

* Customer portal
* Administrator dashboard
* Authentication interface
* Package selection
* Subscription management
* Payment interface
* Customer management
* Billing dashboard
* Network monitoring interface
* Reports

---

# 7. Frontend Architecture

The frontend communicates with the Go backend through REST APIs.

```text
┌─────────────────────────┐
│        Next.js           │
│                          │
│ Admin Dashboard          │
│ Customer Portal          │
│ Authentication           │
│ Billing                  │
│ Packages                 │
└────────────┬────────────┘
             │
             │ HTTPS / JSON
             ▼
┌─────────────────────────┐
│       Go Backend         │
│                          │
│ REST API                 │
│ Business Logic           │
│ Authentication           │
└────────────┬────────────┘
             │
             ▼
          MySQL
```

The frontend must **never communicate directly with MySQL, MikroTik, FreeRADIUS, or M-Pesa**.

All sensitive operations must go through the backend.

---

# 8. Recommended Frontend Structure

```text
frontend/
│
├── app/
│   ├── login/
│   ├── dashboard/
│   ├── customers/
│   ├── packages/
│   ├── subscriptions/
│   ├── billing/
│   ├── payments/
│   └── settings/
│
├── components/
│   ├── ui/
│   ├── forms/
│   ├── tables/
│   └── dashboard/
│
├── lib/
│   ├── api/
│   ├── auth/
│   └── utils/
│
├── hooks/
│
├── types/
│
├── services/
│
├── public/
│
├── package.json
├── tsconfig.json
└── Dockerfile
```

---

# 9. Database

## 9.1 Development Database

**MySQL** will be used as the primary development database.

The database should be designed using relational principles.

Initial entities are expected to include:

```text
users
customers
internet_packages
subscriptions
invoices
payments
routers
radius_accounts
network_sessions
notifications
audit_logs
```

The final schema should be created after the business requirements are finalized.

---

# 10. Database Principles

The database must follow these rules:

1. Use primary keys for all major entities.
2. Use foreign keys where relationships exist.
3. Avoid duplicated data.
4. Use appropriate indexes.
5. Use transactions for financial operations.
6. Never store plaintext passwords.
7. Never store M-Pesa PINs.
8. Store payment transaction identifiers for reconciliation.
9. Record timestamps for important operations.
10. Maintain audit records for sensitive administrative actions.

---

# 11. Payment Architecture

M-Pesa integration will be handled exclusively by the Go backend.

```text
Customer
   │
   ▼
Next.js
   │
   ▼
Go API
   │
   ▼
M-Pesa
   │
   ▼
STK Push
   │
   ▼
Customer Payment
   │
   ▼
M-Pesa Callback
   │
   ▼
Go API
   │
   ▼
Verify Transaction
   │
   ▼
Create Payment Record
   │
   ▼
Activate Subscription
```

The frontend must never determine whether a payment succeeded.

The backend is the authoritative source for payment status.

---

# 12. Subscription Activation

A subscription should only become active after successful payment verification.

```text
PAYMENT_PENDING
       │
       ▼
Payment Received
       │
       ▼
Payment Verified
       │
       ▼
PAYMENT_COMPLETED
       │
       ▼
Subscription Activated
       │
       ▼
Network Access Enabled
```

Failed payments must not activate internet access.

---

# 13. MikroTik Integration

MikroTik communication will be handled by the Go backend.

The frontend must never contain:

```text
MikroTik IP
MikroTik username
MikroTik password
Router API credentials
```

The backend will provide a dedicated network abstraction.

```text
Network Service
      │
      ├── MikroTik Adapter
      │
      ├── FreeRADIUS Adapter
      │
      └── Future Network Adapter
```

This allows additional network equipment to be supported later without rewriting the billing system.

---

# 14. FreeRADIUS

FreeRADIUS will handle customer network authentication and accounting.

The intended flow is:

```text
Customer Device
       │
       ▼
MikroTik
       │
       ▼
FreeRADIUS
       │
       ▼
Authentication
       │
       ▼
Access Granted
```

Afritech Online will manage the customer and subscription information while FreeRADIUS handles network authentication.

---

# 15. Authentication

The backend will use token-based authentication.

The system should support:

* Login
* Logout
* Access tokens
* Refresh tokens
* Password hashing
* Role-based authorization
* Session invalidation

Initial roles may include:

```text
SUPER_ADMIN
ADMIN
STAFF
CUSTOMER
```

Permissions should be based on roles and specific capabilities rather than simply hiding frontend pages.

The backend must enforce authorization independently of the frontend.

---

# 16. API Security

Production APIs must use HTTPS.

Sensitive endpoints should require authentication.

Example:

```text
POST /api/v1/auth/login
POST /api/v1/auth/refresh

GET /api/v1/customers
POST /api/v1/customers

GET /api/v1/packages
POST /api/v1/packages

POST /api/v1/payments/initiate

POST /api/v1/payments/mpesa/callback

GET /api/v1/subscriptions
POST /api/v1/subscriptions
```

The M-Pesa callback endpoint requires special handling because it is called by an external payment provider.

Payment callbacks must be validated and processed idempotently.

---

# 17. API Versioning

All public API endpoints should be versioned.

Use:

```text
/api/v1/
```

Example:

```text
/api/v1/customers
/api/v1/packages
/api/v1/payments
```

Future breaking changes can then use:

```text
/api/v2/
```

---

# 18. Environment Configuration

Secrets must never be hardcoded.

Development configuration should use `.env`.

Example:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_NAME=afritechonline
DB_USER=root
DB_PASSWORD=

JWT_SECRET=

MPESA_CONSUMER_KEY=
MPESA_CONSUMER_SECRET=
MPESA_SHORTCODE=
MPESA_PASSKEY=

MIKROTIK_HOST=
MIKROTIK_USERNAME=
MIKROTIK_PASSWORD=

RADIUS_HOST=
RADIUS_SECRET=
```

A `.env.example` file should be committed to Git.

The real `.env` file must be included in `.gitignore`.

---

# 19. Docker

Docker will be used to make development environments reproducible.

The development environment should eventually contain:

```text
┌──────────────────────┐
│      Next.js         │
│      Frontend        │
└──────────┬───────────┘
           │
┌──────────▼───────────┐
│       Go API         │
│       Backend        │
└──────────┬───────────┘
           │
┌──────────▼───────────┐
│        MySQL         │
└──────────────────────┘

Optional:
┌──────────────────────┐
│     FreeRADIUS       │
└──────────────────────┘
```

MikroTik will normally remain an external network device rather than a Docker container.

---

# 20. Development Workflow

Clone the repository:

```bash
git clone https://github.com/Thorium234/afritechonline.git
cd afritechonline
```

Development branches should follow:

```text
main
develop
feature/*
fix/*
hotfix/*
```

Example:

```bash
git checkout -b feature/customer-management
```

Commit messages should follow a consistent convention.

Examples:

```text
feat: add customer management
feat: implement mpesa payment initiation
fix: prevent duplicate payment processing
refactor: separate mikrotik service
docs: update API documentation
test: add subscription service tests
```

---

# 21. Testing

Testing is mandatory for core business logic.

Backend tests should cover:

* Authentication
* Authorization
* Customer creation
* Package creation
* Subscription creation
* Billing calculations
* Payment processing
* Payment idempotency
* Subscription activation
* Subscription expiration
* MikroTik provisioning
* Error handling

The most important business rules must have automated tests.

For example:

```text
Successful Payment
        ↓
Subscription Activated

Failed Payment
        ↓
Subscription NOT Activated

Duplicate Callback
        ↓
Payment NOT Duplicated
```

---

# 22. Logging

The backend should implement structured logging.

Important events include:

* Login attempts
* Authentication failures
* Payment requests
* Payment callbacks
* Payment verification
* Subscription activation
* Subscription expiration
* MikroTik operations
* Administrative actions
* System errors

Passwords, tokens, payment secrets, and router credentials must never be written to logs.

---

# 23. Monitoring

The production system should eventually support:

* Application health checks
* Database health checks
* Router connectivity monitoring
* API response monitoring
* Error tracking
* Payment failure monitoring
* Authentication failure monitoring

The backend should expose a health endpoint such as:

```text
GET /health
```

Example response:

```json
{
  "status": "ok"
}
```

---

# 24. Deployment

The application should be deployable using Docker.

Production architecture:

```text
                   INTERNET
                       │
                       ▼
                ┌──────────────┐
                │ Reverse Proxy│
                │ HTTPS / TLS  │
                └──────┬───────┘
                       │
             ┌─────────┴─────────┐
             ▼                   ▼
       Next.js Frontend       Go API
                                 │
                     ┌───────────┼───────────┐
                     ▼           ▼           ▼
                   MySQL     FreeRADIUS   MikroTik
```

Production credentials must be supplied through secure environment configuration or a dedicated secret management system.

---

# 25. Development Priority

Development should follow this order:

### Phase 1

```text
Project Setup
    ↓
Go Backend
    ↓
MySQL
    ↓
Database Migrations
    ↓
Authentication
```

### Phase 2

```text
Customers
    ↓
Internet Packages
    ↓
Subscriptions
    ↓
Billing
```

### Phase 3

```text
M-Pesa
    ↓
Payment Verification
    ↓
Subscription Activation
```

### Phase 4

```text
FreeRADIUS
    ↓
MikroTik
    ↓
Network Provisioning
```

### Phase 5

```text
Next.js Dashboard
    ↓
Customer Portal
    ↓
Reports
    ↓
Monitoring
```

---

# 26. Non-Negotiable Architecture Rules

The following rules apply throughout the project.

### Rule 1: Frontend never accesses the database

```text
Next.js → Go API → MySQL
```

Never:

```text
Next.js → MySQL
```

### Rule 2: Frontend never controls MikroTik directly

```text
Next.js → Go API → Network Service → MikroTik
```

### Rule 3: Payment status comes from the backend

The frontend cannot declare:

```text
payment = successful
```

The backend must verify the transaction.

### Rule 4: Business logic belongs in services

Do not place billing, subscription, payment, or network logic directly inside HTTP handlers.

### Rule 5: Secrets never enter Git

Never commit:

```text
.env
passwords
API keys
M-Pesa credentials
MikroTik credentials
RADIUS secrets
JWT secrets
```

### Rule 6: Financial operations must be transactional

Payment records and subscription activation must be handled safely so that partial operations cannot create inconsistent account states.

### Rule 7: External operations must be resilient

M-Pesa, MikroTik, and FreeRADIUS can fail or become unavailable.

The backend must handle:

```text
Timeouts
Retries
Duplicate callbacks
Connection failures
Partial failures
Invalid responses
```

### Rule 8: Design for multiple routers

Do not build the application around one hardcoded MikroTik router.

The data model should eventually support:

```text
ISP
 │
 ├── Router 1
 ├── Router 2
 ├── Router 3
 └── Router N
```

---

# 27. MVP Definition

The first production-quality MVP is considered successful when the following flow works reliably:

```text
Admin
  │
  ▼
Create Internet Package
  │
  ▼
Create Customer
  │
  ▼
Customer Selects Package
  │
  ▼
M-Pesa Payment
  │
  ▼
Payment Verified
  │
  ▼
Subscription Activated
  │
  ▼
FreeRADIUS Account Active
  │
  ▼
MikroTik Authenticates Customer
  │
  ▼
Internet Access Granted
  │
  ▼
Subscription Expires
  │
  ▼
Access Disabled
```

This is the core product.

Everything else is secondary until this workflow is reliable.

---

# 28. Technology Decision Summary

```text
Frontend
    Next.js
    TypeScript

Backend
    Go
    Gin

Database
    MySQL

Authentication
    JWT / Refresh Tokens

Network Authentication
    FreeRADIUS

Network Infrastructure
    MikroTik RouterOS

Payments
    M-Pesa

Development
    Docker
    Git
    GitHub

API
    REST
    JSON
    OpenAPI

Testing
    Go testing
    Frontend automated tests

Deployment
    Docker
    HTTPS
```

---

## Final Architecture

The intended Afritech Online architecture is:

```text
                         ┌───────────────────────┐
                         │       CUSTOMER        │
                         └───────────┬───────────┘
                                     │
                                     ▼
                         ┌───────────────────────┐
                         │       Next.js         │
                         │       Frontend        │
                         └───────────┬───────────┘
                                     │
                                  HTTPS
                                     │
                                     ▼
                         ┌───────────────────────┐
                         │       Go + Gin        │
                         │       Backend         │
                         └───────────┬───────────┘
                                     │
              ┌──────────────────────┼──────────────────────┐
              │                      │                      │
              ▼                      ▼                      ▼
       ┌─────────────┐       ┌─────────────┐       ┌─────────────┐
       │    MySQL    │       │   M-Pesa    │       │   Network   │
       │             │       │             │       │   Services  │
       └─────────────┘       └─────────────┘       └──────┬──────┘
                                                           │
                                             ┌─────────────┴─────────────┐
                                             │                           │
                                             ▼                           ▼
                                      ┌─────────────┐             ┌─────────────┐
                                      │ FreeRADIUS  │             │   MikroTik  │
                                      └─────────────┘             └─────────────┘
```

**This document defines the default technology direction for Afritech Online. Any change to the core stack should be justified technically rather than introduced casually.**
