# Afritech Online

**Afritech Online** is an ISP management and internet billing platform designed to automate customer management, internet packages, subscriptions, billing, payments, network authentication, and MikroTik network management.

The platform is intended for small and medium-sized Internet Service Providers (ISPs), Wi-Fi hotspot operators, and community networks that need a centralized system for managing customers and internet access.

> **Project Status:** 🚧 Under active development
> **Current Phase:** Phase 1, Go + MySQL Backend

---

## Overview

Managing an ISP manually becomes increasingly difficult as the number of customers grows. Customer registration, package assignment, payments, internet activation, session tracking, and subscription expiration can become disconnected processes.

Afritech Online aims to bring these operations into one platform.

The long-term system will connect:

* Customer management
* Internet packages
* Subscriptions
* Billing
* M-Pesa payments
* MikroTik RouterOS
* FreeRADIUS
* Internet authentication
* Network sessions
* Automatic account expiration
* Notifications
* Reporting
* Network monitoring

The core objective is to automate the customer lifecycle:

```text
Customer
    │
    ▼
Select Internet Package
    │
    ▼
Create Subscription
    │
    ▼
Make Payment
    │
    ▼
Payment Verification
    │
    ▼
Activate Subscription
    │
    ▼
Network Provisioning
    │
    ▼
Internet Access
    │
    ▼
Subscription Expires
    │
    ▼
Network Access Disabled
```

The final system is intended to remove the need for administrators to manually activate and deactivate customers whenever payments are received or subscriptions expire.

---

# Technology Stack

| Component              | Technology        |
| ---------------------- | ----------------- |
| Backend                | Go                |
| Backend Framework      | Gin               |
| Frontend               | Next.js           |
| Frontend Language      | TypeScript        |
| Development Database   | MySQL             |
| Network Authentication | FreeRADIUS        |
| Network Infrastructure | MikroTik RouterOS |
| Payment Platform       | M-Pesa            |
| API                    | REST + JSON       |
| Containers             | Docker            |
| Version Control        | Git + GitHub      |
| CI/CD                  | GitHub Actions    |

The technology stack is intentionally separated into phases. MikroTik, FreeRADIUS, and M-Pesa will not be introduced until the core application is stable.

---

# Architecture

The final architecture is expected to evolve into:

```text
                         AFRITECH ONLINE
                               │
                ┌──────────────┴──────────────┐
                │                             │
                ▼                             ▼
        ┌─────────────────┐           ┌─────────────────┐
        │     Next.js     │           │ Customer Portal │
        │    Frontend     │           │                 │
        └────────┬────────┘           └────────┬────────┘
                 │                             │
                 └──────────────┬──────────────┘
                                │
                              HTTPS
                                │
                                ▼
                       ┌─────────────────┐
                       │   Go + Gin API  │
                       │     Backend     │
                       └────────┬────────┘
                                │
               ┌────────────────┼────────────────┐
               │                │                │
               ▼                ▼                ▼
          ┌─────────┐     ┌─────────────┐  ┌──────────────┐
          │  MySQL  │     │   M-Pesa    │  │   Network    │
          │         │     │             │  │   Services   │
          └─────────┘     └─────────────┘  └──────┬───────┘
                                                  │
                                      ┌───────────┴───────────┐
                                      │                       │
                                      ▼                       ▼
                               ┌─────────────┐         ┌─────────────┐
                               │ FreeRADIUS  │         │   MikroTik   │
                               └─────────────┘         └─────────────┘
```

The frontend will never communicate directly with MySQL, MikroTik, FreeRADIUS, or M-Pesa.

All business operations will pass through the Go backend.

```text
Next.js
   │
   ▼
Go API
   │
   ├── MySQL
   ├── M-Pesa
   ├── MikroTik
   └── FreeRADIUS
```

---

# Development Phases

Afritech Online will be developed incrementally.

Each phase must produce a working and testable capability before the project moves to the next phase.

---

# Phase 0: Project Foundation

## Objective

Establish the project structure, development environment, configuration system, database connection, logging, and basic API infrastructure.

### Tasks

* [ ] Establish repository structure
* [ ] Initialize Go backend
* [ ] Configure Gin
* [ ] Configure MySQL
* [ ] Configure environment variables
* [ ] Add database migration system
* [ ] Add structured logging
* [ ] Add centralized error handling
* [ ] Add API response structure
* [ ] Add `/health` endpoint
* [ ] Create Docker development environment
* [ ] Create `.env.example`
* [ ] Configure `.gitignore`

### Expected Structure

```text
afritechonline/
│
├── backend/
├── frontend/
├── infrastructure/
├── docs/
├── scripts/
│
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md
```

### Exit Criteria

Phase 0 is complete when:

```text
Go application starts
        ↓
MySQL connects
        ↓
Configuration loads
        ↓
Database migrations run
        ↓
/health returns OK
```

---

# Phase 1: Go + MySQL Backend

## Objective

Build the core ISP management backend without depending on the frontend, MikroTik, FreeRADIUS, or M-Pesa.

This is the current primary development phase.

### Core Modules

```text
backend/
│
├── auth/
├── users/
├── customers/
├── packages/
├── subscriptions/
├── billing/
├── payments/
└── audit/
```

### Authentication

Implement:

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

Initial roles:

```text
SUPER_ADMIN
ADMIN
STAFF
CUSTOMER
```

### Customers

Implement:

```text
POST   /api/v1/customers
GET    /api/v1/customers
GET    /api/v1/customers/:id
PUT    /api/v1/customers/:id
DELETE /api/v1/customers/:id
```

### Internet Packages

Packages should support:

* Package name
* Price
* Duration
* Download speed
* Upload speed
* Data limit
* Status

Example:

```text
Home Basic
KES 1,000
30 Days
10 Mbps Download
5 Mbps Upload
```

### Subscriptions

Subscription states:

```text
PENDING
ACTIVE
EXPIRED
SUSPENDED
CANCELLED
```

The backend will manage:

```text
Customer
    │
    ▼
Package
    │
    ▼
Subscription
    │
    ├── Start Date
    ├── Expiry Date
    └── Status
```

### Billing

Implement:

```text
Invoices
Payment Records
Payment Status
Subscription Charges
```

Payment states:

```text
PENDING
COMPLETED
FAILED
CANCELLED
```

At this phase, payments can be recorded manually.

**M-Pesa is deliberately excluded from Phase 1.**

### Phase 1 Exit Criteria

The following must work entirely through the Go API:

```text
Create Administrator
        ↓
Login
        ↓
Create Customer
        ↓
Create Internet Package
        ↓
Create Subscription
        ↓
Generate Invoice
        ↓
Record Payment
        ↓
Activate Subscription
```

Automated tests must cover the critical business logic.

---

# Phase 2: Next.js Frontend

## Objective

Build the web interface on top of the completed Go API.

The frontend will consume the API rather than implementing business logic itself.

```text
Next.js
   │
   │ REST / JSON
   ▼
Go Backend
   │
   ▼
MySQL
```

### Frontend Foundation

* [ ] Next.js
* [ ] TypeScript
* [ ] UI system
* [ ] API client
* [ ] Authentication
* [ ] Form validation
* [ ] Error handling
* [ ] Loading states
* [ ] Protected routes

### Main Interfaces

```text
/login

/dashboard

/customers
/customers/new
/customers/[id]

/packages
/packages/new
/packages/[id]

/subscriptions
/subscriptions/[id]

/billing
/invoices
/payments

/settings
```

### Admin Dashboard

The dashboard should eventually provide:

```text
Total Customers
Active Customers
Expired Customers
Active Packages
Active Subscriptions
Pending Payments
Revenue
```

### Phase 2 Exit Criteria

An administrator must be able to perform the Phase 1 workflow through the browser:

```text
Login
  ↓
Dashboard
  ↓
Create Customer
  ↓
Create Package
  ↓
Create Subscription
  ↓
Generate Invoice
  ↓
Record Payment
  ↓
View Active Subscription
```

---

# Phase 3: MikroTik Integration

## Objective

Connect the Go backend to the physical MikroTik router and introduce network automation.

The MikroTik router has already been connected and verified through WinBox, providing the initial hardware environment for this phase.

MikroTik integration will be introduced gradually.

---

## Phase 3A: MikroTik Connectivity

First establish communication between Go and the router.

```text
Go Backend
     │
     ▼
MikroTik API
     │
     ▼
Router
```

Initial capabilities:

* [ ] Connect to router
* [ ] Authenticate
* [ ] Get router identity
* [ ] Get RouterOS version
* [ ] Check connection status
* [ ] Handle connection failures
* [ ] Disconnect safely

Example endpoint:

```text
GET /api/v1/routers/:id/status
```

### Exit Criteria

The backend can reliably connect to the MikroTik and retrieve basic router information.

No customer provisioning yet.

---

# Phase 3B: Router Management

Create router management inside Afritech Online.

```text
routers
```

Potential information:

```text
Router ID
Router Name
IP Address
API Port
Username
Encrypted Credentials
Status
Created At
Updated At
```

API:

```text
POST   /api/v1/routers
GET    /api/v1/routers
GET    /api/v1/routers/:id
PUT    /api/v1/routers/:id
DELETE /api/v1/routers/:id
```

The administrator should be able to test a router connection from the Afritech Online dashboard.

---

# Phase 3C: MikroTik User Management

The backend will begin managing actual network users.

Capabilities:

```text
Create User
Get User
Enable User
Disable User
Delete User
```

The Go network service will translate Afritech Online operations into the appropriate RouterOS operations.

---

# Phase 3D: Bandwidth Profiles

Connect internet packages with network profiles.

Example:

```text
Afritech 5 Mbps
Afritech 10 Mbps
Afritech 20 Mbps
```

Architecture:

```text
Internet Package
       │
       ▼
Network Profile
       │
       ▼
MikroTik
```

A customer's purchased package must determine the appropriate network configuration.

---

# Phase 3E: Automatic Network Provisioning

Connect subscription activation to MikroTik provisioning.

```text
Customer
    ↓
Subscription
    ↓
Payment
    ↓
ACTIVE
    ↓
Network Provisioning
    ↓
MikroTik
    ↓
Internet Access
```

The goal is to eliminate manual router configuration for every customer.

---

# Phase 3F: Automatic Expiration

When a subscription expires:

```text
Subscription
      │
      ▼
Expiry Date Reached
      │
      ▼
Subscription = EXPIRED
      │
      ▼
Network Provisioning Service
      │
      ▼
Disable Customer
```

This operation should eventually run automatically without administrator intervention.

---

# Phase 4: FreeRADIUS

## Objective

Introduce FreeRADIUS for network authentication, authorization, and accounting.

The initial architecture is:

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
Internet Access
```

FreeRADIUS will eventually handle:

* User authentication
* Authorization
* Session accounting
* Bandwidth profiles
* Access policies
* Session start/stop information

FreeRADIUS is intentionally introduced after direct MikroTik integration so that network problems can be isolated during development.

---

# Phase 5: M-Pesa Integration

## Objective

Introduce real payment processing after the billing and subscription systems are stable.

Architecture:

```text
Customer
    │
    ▼
Next.js
    │
    ▼
Go Backend
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
Payment Verification
    │
    ▼
Payment COMPLETED
    │
    ▼
Subscription ACTIVE
    │
    ▼
Network Provisioning
```

### Requirements

* [ ] STK Push
* [ ] Callback endpoint
* [ ] Transaction verification
* [ ] Payment reconciliation
* [ ] Duplicate callback protection
* [ ] Failed payment handling
* [ ] Payment audit trail
* [ ] Subscription activation

### Critical Rule

The frontend must never determine that a payment succeeded.

Only the backend can confirm payment status.

Payment processing must also be **idempotent**.

If the same callback is received twice, the system must not:

```text
Create two payments
Extend a subscription twice
Provision a customer twice
```

---

# Phase 6: Automation Engine

## Objective

Move recurring operations into background workers.

Potential automated tasks:

```text
Subscription Expiration
Payment Reconciliation
Network Provisioning
Provisioning Retries
Invoice Generation
Notifications
```

Architecture:

```text
                 Go Backend
                     │
             ┌───────┴───────┐
             │               │
          API Server       Workers
                             │
                ┌────────────┼────────────┐
                ▼            ▼            ▼
           Expiration     Payments   Provisioning
```

Redis or another queue system can be introduced when the workload requires it.

---

# Phase 7: Monitoring and Reporting

## Objective

Provide administrators with visibility into the ISP operation.

### Dashboard

```text
Customers
Active Customers
Online Customers
Expired Customers
Revenue
Failed Payments
Active Subscriptions
Routers
Router Status
```

### Reports

```text
Daily Revenue
Monthly Revenue
Payment History
Active Subscriptions
Expired Subscriptions
Customer Growth
Package Performance
Network Sessions
```

### Router Monitoring

Potential metrics:

```text
Router Status
CPU Usage
Memory Usage
Uptime
Active Sessions
Connectivity
```

---

# Phase 8: Production Hardening

Before declaring Afritech Online production-ready, implement:

* [ ] HTTPS
* [ ] Secure secret management
* [ ] Database backups
* [ ] Rate limiting
* [ ] Input validation
* [ ] RBAC
* [ ] Audit logging
* [ ] Security headers
* [ ] API documentation
* [ ] Monitoring
* [ ] Error tracking
* [ ] Automated testing
* [ ] CI/CD
* [ ] Docker deployment
* [ ] Disaster recovery procedures

The system must also be tested against:

```text
Payment failures
Router failures
Database failures
Network outages
Duplicate callbacks
Authentication failures
Expired subscriptions
Partial provisioning failures
```

---

# MVP Definition

The first serious MVP is **not** the entire platform.

The MVP is the successful completion of the core customer lifecycle:

```text
Administrator
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
Payment
      │
      ▼
Payment Verified
      │
      ▼
Subscription Activated
      │
      ▼
Network Account Provisioned
      │
      ▼
Customer Authenticated
      │
      ▼
Internet Access
      │
      ▼
Subscription Expires
      │
      ▼
Internet Access Disabled
```

---

# Development Rules

The following rules apply throughout development.

## 1. Backend First

The Go backend must be functional independently of the frontend.

```text
Go API
   ↓
MySQL
```

The frontend comes afterward.

## 2. Frontend Never Accesses MySQL

Never:

```text
Next.js → MySQL
```

Always:

```text
Next.js → Go API → MySQL
```

## 3. Frontend Never Controls MikroTik

Never:

```text
Next.js → MikroTik
```

Always:

```text
Next.js
   ↓
Go API
   ↓
Network Service
   ↓
MikroTik
```

## 4. Separate Network Logic

MikroTik-specific functionality must be isolated behind a network service or adapter.

```text
Business Logic
      │
      ▼
Network Provisioning Service
      │
      ├── MikroTik
      ├── FreeRADIUS
      └── Future Network Providers
```

This allows additional network technologies to be introduced without rewriting the billing system.

## 5. Payment Verification Is Backend Responsibility

Never trust the frontend to report a successful payment.

The backend must verify the payment.

## 6. Secrets Never Enter Git

Never commit:

```text
.env
Passwords
JWT secrets
M-Pesa credentials
MikroTik credentials
RADIUS secrets
API keys
```

Use `.env.example` for required configuration names.

---

# Security

Afritech Online will handle customer information, payment information, network credentials, and infrastructure access. Security therefore needs to be part of the architecture from the beginning.

The system must:

* Hash passwords
* Never store M-Pesa PINs
* Never expose MikroTik credentials to the frontend
* Store secrets outside source control
* Validate external callbacks
* Enforce authorization on the backend
* Use HTTPS in production
* Validate incoming data
* Implement rate limiting
* Maintain audit logs
* Avoid logging credentials and tokens
* Use database transactions for financial operations

Example environment configuration:

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

**Never commit real credentials to GitHub.**

---

# Planned Repository Structure

```text
afritechonline/
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── customers/
│   │   ├── packages/
│   │   ├── subscriptions/
│   │   ├── billing/
│   │   ├── payments/
│   │   ├── mikrotik/
│   │   ├── radius/
│   │   ├── sessions/
│   │   ├── notifications/
│   │   └── reports/
│   │
│   ├── database/
│   │   ├── migrations/
│   │   └── seeds/
│   │
│   ├── middleware/
│   ├── routes/
│   ├── config/
│   ├── tests/
│   │
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── app/
│   ├── components/
│   ├── hooks/
│   ├── lib/
│   ├── services/
│   ├── types/
│   ├── public/
│   └── package.json
│
├── infrastructure/
│   ├── docker/
│   ├── freeradius/
│   └── mikrotik/
│
├── docs/
│   ├── architecture.md
│   ├── database.md
│   ├── api.md
│   ├── deployment.md
│   └── phases.md
│
├── scripts/
│
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md
```

---

# Development Roadmap

```text
┌──────────────────────────────────┐
│ PHASE 0                          │
│ Project Foundation               │
│ Go + MySQL + Configuration       │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 1                          │
│ Go + MySQL Backend               │
│ Auth + Customers + Billing       │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 2                          │
│ Next.js Frontend                 │
│ Dashboard + Customer Portal      │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 3                          │
│ MikroTik Integration             │
│ Connectivity + Provisioning      │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 4                          │
│ FreeRADIUS                       │
│ Authentication + Accounting      │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 5                          │
│ M-Pesa                           │
│ Payments + Verification          │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 6                          │
│ Automation Engine                │
│ Workers + Scheduled Operations   │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 7                          │
│ Monitoring + Reporting           │
└───────────────┬──────────────────┘
                │
                ▼
┌──────────────────────────────────┐
│ PHASE 8                          │
│ Production Hardening             │
│ Security + CI/CD + Deployment    │
└───────────────┬──────────────────┘
                │
                ▼
          AFRITECH ONLINE
             v1.0.0
```

---

# Current Development Target

The immediate target is **Phase 1**.

Do not jump ahead to MikroTik, FreeRADIUS, or M-Pesa before the core backend is reliable.

The current milestone is:

```text
Go
 │
 ▼
Gin
 │
 ▼
MySQL
 │
 ▼
Migrations
 │
 ▼
Authentication
 │
 ▼
Customers
 │
 ▼
Internet Packages
 │
 ▼
Subscriptions
 │
 ▼
Invoices
 │
 ▼
Payments
 │
 ▼
Automated Tests
 │
 ▼
PHASE 1 COMPLETE
```

Once this works reliably, Phase 2 begins.

---

# Future Features

Future versions may include:

* Multiple MikroTik routers
* Multiple ISP locations
* Multi-tenant ISP management
* Reseller accounts
* Agent accounts
* Customer self-service
* SMS notifications
* Email notifications
* M-Pesa PayBill
* M-Pesa Till
* Airtel Money
* Card payments
* Voucher generation
* Hotspot management
* PPPoE management
* Bandwidth monitoring
* Usage analytics
* Revenue reports
* Automatic invoices
* Router health monitoring
* Mobile application

These features will only be prioritized after the core billing and network provisioning workflow is stable.

---

# Contributing

Development should follow the phased architecture.

Create a feature branch:

```bash
git checkout -b feature/customer-management
```

Make changes, test them, and commit using a consistent convention:

```bash
git commit -m "feat: add customer management"
```

Other examples:

```text
feat: add subscription management
feat: implement payment records
fix: prevent duplicate payment processing
refactor: separate network provisioning service
test: add subscription service tests
docs: update API documentation
```

Push the branch:

```bash
git push origin feature/customer-management
```

Then open a pull request.

---

# License

License information will be added when the project's licensing model is finalized.

---

# Project Status

Afritech Online is under active development.

The project has moved from initial architecture and planning toward implementation.

The development strategy is intentionally incremental:

**Foundation → Backend → Frontend → MikroTik → FreeRADIUS → M-Pesa → Automation → Monitoring → Production**

The immediate objective is to build a reliable Go + MySQL backend before introducing external network and payment dependencies.

---

# Repository

**GitHub:** https://github.com/Thorium234/afritechonline

**Project:** Afritech Online

**Purpose:** ISP management, billing, payment, authentication, and network automation.
