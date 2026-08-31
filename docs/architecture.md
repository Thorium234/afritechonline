# Architecture

## Overview

Afritech Online follows a clean layered architecture with clear separation of concerns.

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                    │
│                    TypeScript + Tailwind CSS                  │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                │ REST API
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go + Gin)                      │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌──────────────────────┐  │
│  │  Handlers   │→│  Services  │→│    Repositories       │  │
│  │  (HTTP)     │  │ (Business) │  │   (Database)         │  │
│  └────────────┘  └────────────┘  └──────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                    Middleware                         │   │
│  │  Auth, Rate Limiting, CORS, Security Headers          │   │
│  └──────────────────────────────────────────────────────┘   │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                │
        ┌───────────────────────┼───────────────────────┐
        │                       │                       │
        ▼                       ▼                       ▼
┌──────────────┐      ┌──────────────────┐    ┌───────────────┐
│   MySQL      │      │  MikroTik SSH    │    │  M-Pesa API   │
│  (Database)  │      │  (Network)       │    │  (Payments)   │
└──────────────┘      └──────────────────┘    └───────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │  FreeRADIUS      │
                       │  (Authentication)│
                       └──────────────────┘
```

## Design Principles

1. **Dependency Injection**: Services depend on interfaces, not concrete implementations.
2. **Idempotency**: Payment completion, subscription activation are idempotent.
3. **Security First**: JWT with refresh tokens, role-based access, rate limiting.
4. **Testability**: Clear interfaces enable unit testing with mocks.
5. **Transaction Safety**: Critical flows use database transactions.

## Package Structure

```
backend/
├── cmd/server/          # Application entrypoint
├── config/              # Configuration management
├── internal/
│   ├── auth/            # Authentication (register, login, tokens)
│   ├── customers/       # Customer management
│   ├── packages/        # Internet packages
│   ├── subscriptions/   # Subscriptions lifecycle
│   ├── invoices/        # Invoice generation
│   ├── payments/        # Payment processing
│   ├── mikrotik/        # MikroTik integration
│   ├── radius/          # FreeRADIUS integration
│   ├── mpesa/           # M-Pesa integration
│   ├── workers/         # Background jobs
│   ├── reports/         # Reporting and analytics
│   ├── database/        # DB connection, migrations, seeding
│   └── models/          # Domain models
├── middleware/          # HTTP middleware
├── pkg/                 # Shared utilities
│   ├── token/           # JWT management
│   ├── validator/       # Input validation
│   ├── response/        # JSON response helpers
│   ├── logger/          # Logging setup
│   └── contextutil/     # Context utilities
└── routes/              # Route definitions
```

## Data Flow

```
Customer
    │
    ▼
Select Package
    │
    ▼
Create Subscription (PENDING)
    │
    ▼
Generate Invoice
    │
    ▼
Record Payment (PENDING)
    │
    ▼
Complete Payment
    │
    ├── Mark Invoice PAID
    ├── Activate Subscription
    └── Provision Network Account
            │
            ▼
       MikroTik / FreeRADIUS
            │
            ▼
       Customer Gets Internet
```
