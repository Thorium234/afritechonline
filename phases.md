Yes. This is the right way to approach it. But I would **not** make the phases simply “backend, frontend, MikroTik.” That is too shallow. You need to separate the project around **working capabilities**, with each phase producing something testable.

Since your MikroTik is already reachable through WinBox, we have a real network device available for Phase 3 testing. That is valuable, but we should deliberately keep it out of the early phases so we don't mix network problems with application problems.

# Afritech Online Development Plan

## Overall Architecture

The final system should evolve toward:

```text
                         AFRITECH ONLINE
                               │
              ┌────────────────┴────────────────┐
              │                                 │
              ▼                                 ▼
       ┌─────────────┐                   ┌─────────────┐
       │   Next.js   │                   │   Customer  │
       │   Frontend  │                   │    Portal   │
       └──────┬──────┘                   └──────┬──────┘
              │                                 │
              └──────────────┬──────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   Go Backend    │
                    │      API        │
                    └────────┬────────┘
                             │
             ┌───────────────┼────────────────┐
             │               │                │
             ▼               ▼                ▼
        ┌─────────┐    ┌─────────────┐  ┌──────────────┐
        │  MySQL  │    │    Redis    │  │ Integrations │
        └─────────┘    └─────────────┘  └──────┬───────┘
                                                │
                              ┌─────────────────┼──────────────┐
                              │                 │              │
                              ▼                 ▼              ▼
                         MikroTik          FreeRADIUS       M-Pesa
```

Redis is **not mandatory for Phase 1**. It becomes useful later for queues, caching, sessions, rate limiting, and background jobs.

---

# Phase 0: Project Foundation

Before writing business functionality, establish the engineering foundation.

### Objectives

Create a clean repository structure and development environment.

```text
afritechonline/
│
├── backend/
├── frontend/
├── infrastructure/
├── docs/
├── scripts/
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md
```

### Backend foundation

Set up:

```text
Go
Gin
MySQL driver
Configuration management
Logging
Error handling
Database connection
Migration system
```

Create:

```text
GET /health
```

Expected:

```json
{
  "status": "ok"
}
```

### Deliverable

You should be able to run:

```bash
cd backend
go run ./cmd/server
```

and get:

```text
Server running on :8080
Database connected
```

### Exit Criteria

Phase 0 is complete when:

* Go application starts
* MySQL connects
* Configuration loads from `.env`
* `/health` works
* Repository structure is established
* No credentials are hardcoded

---

# Phase 1: Go + MySQL Core Backend

This is the **first real development phase**.

Do not touch MikroTik yet.

Do not build the fancy frontend yet.

We first make the backend a real application.

## 1.1 Database Design

Start with the core entities.

```text
users
customers
internet_packages
subscriptions
invoices
payments
routers
audit_logs
```

Initially, `routers` can exist in the database without actually connecting to a router.

Example relationship:

```text
Customer
    │
    └── Subscription
            │
            └── Internet Package
```

Later:

```text
Customer
    │
    └── Subscription
            │
            ├── Package
            └── Network Account
                       │
                       └── Router
```

## 1.2 Authentication

Implement:

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

Implement password hashing.

Implement roles:

```text
SUPER_ADMIN
ADMIN
STAFF
CUSTOMER
```

## 1.3 Customer Management

Implement:

```text
POST   /api/v1/customers
GET    /api/v1/customers
GET    /api/v1/customers/:id
PUT    /api/v1/customers/:id
DELETE /api/v1/customers/:id
```

Customer information might include:

```text
ID
Full Name
Phone
Email
Username
Status
Created At
Updated At
```

## 1.4 Internet Packages

Implement:

```text
POST   /api/v1/packages
GET    /api/v1/packages
GET    /api/v1/packages/:id
PUT    /api/v1/packages/:id
DELETE /api/v1/packages/:id
```

Example:

```text
10 Mbps
30 Days
KES 1,000
```

## 1.5 Subscriptions

Implement:

```text
POST /api/v1/subscriptions
GET  /api/v1/subscriptions
GET  /api/v1/subscriptions/:id
```

Subscription states:

```text
PENDING
ACTIVE
EXPIRED
SUSPENDED
CANCELLED
```

The backend should calculate:

```text
start_date
expiry_date
status
```

## 1.6 Billing

Create invoices.

```text
POST /api/v1/invoices
GET  /api/v1/invoices
GET  /api/v1/invoices/:id
```

At this stage, payments can simply be recorded manually.

No M-Pesa yet.

## 1.7 Payments

Create:

```text
POST /api/v1/payments
GET  /api/v1/payments
GET  /api/v1/payments/:id
```

Payment states:

```text
PENDING
COMPLETED
FAILED
CANCELLED
```

The important part is designing this correctly now because M-Pesa will eventually feed into the same payment system.

### Phase 1 Exit Criteria

You should be able to demonstrate:

```text
Create Admin
       ↓
Login
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
Subscription becomes ACTIVE
```

Everything works using **API requests only**.

No frontend.

No MikroTik.

If this doesn't work reliably, moving forward is stupid.

---

# Phase 2: Next.js Frontend

Now that the API works independently, build the interface on top of it.

This separation is important.

The frontend should consume the API rather than containing business logic.

```text
Next.js
    │
    │ REST API
    ▼
Go Backend
    │
    ▼
MySQL
```

## 2.1 Frontend Foundation

Set up:

```text
Next.js
TypeScript
Tailwind CSS
API client
Authentication
Form validation
Error handling
Loading states
```

## 2.2 Authentication UI

Build:

```text
/login
```

Then:

```text
/dashboard
```

Protect authenticated routes.

## 2.3 Admin Dashboard

Dashboard should show:

```text
Total Customers
Active Customers
Expired Customers
Active Packages
Pending Payments
Revenue
```

Don't waste time making it visually spectacular yet.

Make it functional.

## 2.4 Customer Management

Create:

```text
/customers
/customers/new
/customers/[id]
```

Functions:

```text
Create
View
Edit
Delete
Search
Filter
```

## 2.5 Package Management

Create:

```text
/packages
/packages/new
/packages/[id]
```

## 2.6 Subscription Management

Create:

```text
/subscriptions
/subscriptions/[id]
```

Display:

```text
Customer
Package
Start Date
Expiry Date
Status
Amount
```

## 2.7 Billing Interface

Create:

```text
/billing
/invoices
/payments
```

The frontend should show the backend's payment status.

It should **never decide whether money was received**.

### Phase 2 Exit Criteria

A user can open the web application and perform:

```text
Login
  ↓
Dashboard
  ↓
Create Customer
  ↓
Create Package
  ↓
Subscribe Customer
  ↓
Generate Invoice
  ↓
Record Payment
  ↓
See Active Subscription
```

At this point you have a functioning **ISP management application**, even though it isn't controlling the network yet.

---

# Phase 3: MikroTik Integration

Now we introduce the hardware.

This phase should itself be broken into multiple stages.

Don't immediately connect billing to router commands.

That's how you create an untestable mess.

---

## Phase 3A: MikroTik Connectivity

First prove that Go can communicate with your MikroTik.

You already have the router working through WinBox.

Now establish:

```text
Go Backend
     │
     ▼
MikroTik API
     │
     ▼
Router
```

Create:

```text
mikrotik/
├── client.go
├── connection.go
├── users.go
├── profiles.go
└── sessions.go
```

The first API should be something simple like:

```text
GET /api/v1/routers/:id/status
```

Expected:

```json
{
  "connected": true,
  "identity": "MikroTik",
  "version": "..."
}
```

### Exit Criteria

Go can:

```text
Connect
Authenticate
Get Router Identity
Get RouterOS Version
Disconnect
```

**Do not modify customers yet.**

---

# Phase 3B: Router Management

Now build router management.

Database:

```text
routers
```

Example:

```text
Router ID
Name
IP Address
API Port
Username
Encrypted Credentials
Status
Created At
```

Backend:

```text
POST /api/v1/routers
GET  /api/v1/routers
GET  /api/v1/routers/:id
PUT  /api/v1/routers/:id
DELETE /api/v1/routers/:id
```

Admin should be able to register a MikroTik router.

Then:

```text
Test Connection
```

The UI should tell the administrator:

```text
CONNECTED
```

or:

```text
CONNECTION FAILED
```

---

# Phase 3C: MikroTik User Management

Now interact with actual network users.

The backend should be able to:

```text
Create network user
Disable network user
Enable network user
Delete network user
Get network user
```

For example:

```text
POST /api/v1/routers/:id/users
```

The Go backend translates the request into the appropriate RouterOS operation.

This is where your WinBox knowledge becomes useful.

---

# Phase 3D: Bandwidth Profiles

Create ISP packages that correspond to actual router profiles.

For example:

```text
Afritech 5Mbps
Afritech 10Mbps
Afritech 20Mbps
```

Then:

```text
Internet Package
       │
       ▼
Network Profile
       │
       ▼
MikroTik
```

A customer purchasing a 10 Mbps package should receive the appropriate network configuration.

---

# Phase 3E: Subscription → MikroTik Automation

Now connect the systems.

This is the important part.

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

The backend should automatically provision the customer after successful activation.

---

# Phase 3F: Expiration Automation

Now implement the reverse.

```text
Subscription
     │
     ▼
expiry_date reached
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

This should happen automatically.

You should not need an administrator to manually disable users every morning.

---

# Phase 4: FreeRADIUS

Only after direct MikroTik integration is working should we introduce FreeRADIUS.

This is important because otherwise you won't know whether a failure comes from:

```text
Go
MikroTik
RADIUS
Database
```

Build it independently.

Architecture:

```text
Customer Device
       ↓
MikroTik
       ↓
FreeRADIUS
       ↓
Authentication
```

Afritech backend:

```text
Go
 ↓
RADIUS Database / API
 ↓
User Authorization
```

Implement:

```text
Username
Password
Package
Speed
Expiration
Session
Accounting
```

---

# Phase 5: M-Pesa

Now introduce real money.

This should happen **after** the billing architecture is stable.

The architecture becomes:

```text
Customer
    ↓
Next.js
    ↓
Go API
    ↓
M-Pesa
    ↓
Callback
    ↓
Payment Verification
    ↓
Payment COMPLETED
    ↓
Subscription ACTIVE
    ↓
Network Provisioning
    ↓
MikroTik
```

Implement:

```text
POST /api/v1/payments/mpesa/stkpush
POST /api/v1/payments/mpesa/callback
GET  /api/v1/payments/:id
```

Critical requirement:

**M-Pesa callback processing must be idempotent.**

If Safaricom sends the same callback twice, you must not:

```text
create two payments
extend subscription twice
activate two accounts
```

---

# Phase 6: Automation Engine

At this point the platform becomes genuinely interesting.

Create background workers for:

```text
Subscription expiration
Payment reconciliation
Network provisioning
Failed provisioning retries
Notifications
Invoice generation
```

Architecture:

```text
                    Go Backend
                        │
                ┌───────┴───────┐
                │               │
             API Server      Workers
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
               Expiration   Payments   Provisioning
```

This is where Redis or a proper job queue becomes useful.

---

# Phase 7: Monitoring and Reporting

Once the core system is working:

### Dashboard

```text
Customers
Active Users
Online Users
Expired Users
Revenue
Failed Payments
Routers
Router Status
```

### Reports

```text
Daily Revenue
Monthly Revenue
Active Subscriptions
Expired Subscriptions
Payment History
Customer Growth
Bandwidth Packages
Router Sessions
```

### Network monitoring

```text
Router Online
Router Offline
CPU
Memory
Uptime
Active Sessions
```

---

# Phase 8: Production Hardening

Only after functionality works.

Implement:

```text
HTTPS
Secret management
Database backups
Rate limiting
Audit logging
Input validation
RBAC
Security headers
API documentation
Monitoring
Error tracking
Automated tests
CI/CD
Docker deployment
```

Then perform:

```text
Load testing
Failure testing
Payment failure testing
Router failure testing
Database failure testing
Network outage testing
Security testing
```

---

# The Complete Roadmap

I would put the project into this exact sequence:

```text
PHASE 0
Project Foundation
       │
       ▼
PHASE 1
Go + MySQL Backend
       │
       ▼
PHASE 2
Next.js Frontend
       │
       ▼
PHASE 3
MikroTik Integration
       │
       ├── 3A Connectivity
       ├── 3B Router Management
       ├── 3C User Management
       ├── 3D Bandwidth Profiles
       ├── 3E Automatic Provisioning
       └── 3F Automatic Expiration
       │
       ▼
PHASE 4
FreeRADIUS
       │
       ▼
PHASE 5
M-Pesa
       │
       ▼
PHASE 6
Automation / Workers
       │
       ▼
PHASE 7
Monitoring & Reporting
       │
       ▼
PHASE 8
Production Hardening
       │
       ▼
             AFRITECH ONLINE v1.0
```

## The most important rule

**Never move to the next phase because you are bored with the current one. Move because the exit criteria have been satisfied.**

For example, don't say:

> "The Go backend is mostly done, let's start MikroTik."

Instead:

> "Phase 1 is complete because authentication, customers, packages, subscriptions, invoices, payments, migrations, validation, error handling, and tests all work."

Then move.

That discipline matters even more because you're working with **payments and network infrastructure**. A bug in a normal CRUD application is annoying. A bug that marks an unpaid customer as paid and automatically gives them internet is a real operational problem.

### Recommended immediate target

Forget Phases 2 through 8 for now.

**Phase 0 → Phase 1.**

The first milestone should be:

```text
GitHub Repository
      ↓
Go + Gin
      ↓
MySQL
      ↓
Migrations
      ↓
Authentication
      ↓
Customers
      ↓
Packages
      ↓
Subscriptions
      ↓
Invoices
      ↓
Payments
      ↓
Automated Tests
      ↓
PHASE 1 COMPLETE
```

Once that is solid, we build Next.js on top of it. Then, and only then, we touch the MikroTik.
