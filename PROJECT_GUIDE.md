# Afritech Online - Complete Project Guide

A comprehensive ISP (Internet Service Provider) management platform with secure authentication, customer management, subscription handling, payment processing, and network management.

**Stack**: Go (backend) + Next.js (frontend) + MySQL (database)

---

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Getting Started](#getting-started)
4. [Backend Documentation](#backend-documentation)
5. [Frontend Documentation](#frontend-documentation)
6. [Security & Deployment](#security--deployment)
7. [API Reference](#api-reference)
8. [Development Guide](#development-guide)

---

## 🎯 Project Overview

**Afritech Online** is an end-to-end ISP management system that handles:

- **Authentication**: Secure JWT-based authentication with role-based access control
- **Customer Management**: Create, update, and manage subscriber accounts
- **Internet Packages**: Define and manage internet service packages
- **Subscriptions**: Handle customer subscriptions to packages
- **Billing & Payments**: Invoice generation and payment processing with M-Pesa integration
- **Network Management**: Router management via MikroTik RADIUS integration
- **Reporting**: Revenue, customer, and network analytics

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Frontend (Next.js)                         │
│  (React 18 + TypeScript + Tailwind CSS + Custom Components) │
└────────────────────────────┬────────────────────────────────┘
                             │ HTTP REST API
                             │
┌─────────────────────────────▼────────────────────────────────┐
│                Backend (Go + Gin)                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │  Auth Layer  │  │  Service     │  │  Repository  │       │
│  │  (Handlers)  │→ │  (Business   │→ │  (Database   │       │
│  │              │  │   Logic)     │  │   Access)    │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                                                               │
│  ┌─────────────────┐  ┌──────────────────────┐              │
│  │ Integrations    │  │ Security Middleware  │              │
│  │ - MikroTik      │  │ - Rate Limiting      │              │
│  │ - FreeRADIUS    │  │ - JWT Validation     │              │
│  │ - M-Pesa Daraja│  │ - CORS Protection    │              │
│  │ - Email Service │  │ - Audit Logging      │              │
│  └─────────────────┘  └──────────────────────┘              │
└───────────────────────────┬────────────────────────────────┘
                            │
          ┌─────────────────┴─────────────────┐
          │                                   │
    ┌─────▼─────┐                      ┌─────▼──────┐
    │   MySQL   │                      │  Redis(*)  │
    │ (Database)│                      │  (Cache)   │
    └───────────┘                      └────────────┘
    
(*) Optional for caching and rate limiting
```

**Data Flow**:
1. Frontend sends authenticated requests with JWT token
2. Backend validates token and checks permissions
3. Service layer applies business logic
4. Repository layer persists/retrieves data
5. Security middleware logs all sensitive operations
6. External integrations handle network and payment operations

---

## 🚀 Getting Started

### Backend Setup

```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your configuration

# Setup database
go run cmd/server/main.go -migrate

# Run server
go run cmd/server/main.go
# Server starts on http://localhost:8000
```

### Frontend Setup

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Configure environment
cp .env.example .env.local
# Edit .env.local with API URL

# Run development server
npm run dev
# Frontend available at http://localhost:3000
```

### Docker Deployment

```bash
# Run entire stack with Docker Compose
docker-compose up -d

# Access:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8000
# - Nginx Proxy: http://localhost:80
```

---

## 📚 Backend Documentation

### Key Files & Structure

```
backend/
├── cmd/server/main.go           # Application entry point
├── config/config.go              # Configuration management
├── middleware/                   # HTTP middleware
│   ├── auth.go                   # JWT validation
│   ├── logger.go                 # Request logging
│   ├── ratelimit.go              # Rate limiting (NEW)
│   └── security.go               # Security headers & CORS
├── internal/
│   ├── auth/                     # Authentication service
│   ├── customers/                # Customer management
│   ├── invoices/                 # Invoice management
│   ├── payments/                 # Payment processing
│   ├── packages/                 # Internet packages
│   ├── subscriptions/            # Subscription management
│   ├── routers/                  # Router management
│   ├── radius/                   # RADIUS integration
│   ├── mpesa/                    # M-Pesa integration
│   ├── mikrotik/                 # MikroTik integration
│   ├── reports/                  # Analytics & reports
│   └── database/                 # Database layer
├── pkg/
│   ├── audit/audit.go            # Audit logging (NEW)
│   ├── validator/password.go     # Password validation (NEW)
│   ├── logger/                   # Structured logging
│   ├── token/                    # JWT token handling
│   └── response/                 # Response formatting
├── routes/routes.go              # Route configuration
├── Dockerfile                    # Container configuration
├── go.mod                        # Dependency management
└── Makefile                      # Build commands
```

### API Endpoints

See [SECURITY_FIXES_SUMMARY.md](SECURITY_FIXES_SUMMARY.md) for complete API reference.

**Authentication**:
```
POST   /api/v1/auth/login        # User login
POST   /api/v1/auth/register     # User registration
POST   /api/v1/auth/refresh      # Refresh JWT token
POST   /api/v1/auth/logout       # User logout
GET    /api/v1/auth/me           # Get current user
```

**Customers**:
```
GET    /api/v1/customers         # List all customers
POST   /api/v1/customers         # Create customer
GET    /api/v1/customers/:id     # Get customer details
PUT    /api/v1/customers/:id     # Update customer
DELETE /api/v1/customers/:id     # Delete customer
```

**Packages, Subscriptions, Payments, etc.** - Similar CRUD patterns.

### Security Features (NEW)

1. **Rate Limiting** (`middleware/ratelimit.go`)
   - Per-IP throttling: 100 requests/min for general endpoints
   - Auth endpoint protection: 15 login attempts/15min
   - Automatic cleanup of expired entries

2. **Audit Logging** (`pkg/audit/audit.go`)
   - 16 event types tracked (USER_LOGIN, PAYMENT_COMPLETED, etc.)
   - Full context: userID, IP, resource details, timestamp
   - JSON structured logs for analysis

3. **Password Validation** (`pkg/validator/password.go`)
   - NIST SP 800-63B compliance
   - 14+ character minimum
   - Common password blacklist
   - Character variety requirements

4. **JWT Secret Validation** (`config/config.go`)
   - 32+ character minimum
   - Rejects known placeholder values
   - Helpful error messages

5. **Database Transactions** (`internal/database/transaction.go`)
   - Atomic operations for payment + invoice + subscription
   - Configurable isolation levels
   - Automatic rollback on errors

---

## 💻 Frontend Documentation

### Key Files & Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── (app)/                 # Authenticated routes
│   │   │   ├── customers/         # Customer pages
│   │   │   ├── subscriptions/     # Subscription pages
│   │   │   ├── payments/          # Payment pages
│   │   │   ├── billing/           # Invoice pages
│   │   │   ├── packages/          # Package pages
│   │   │   ├── routers/           # Router pages
│   │   │   ├── reports/           # Analytics pages
│   │   │   └── layout.tsx         # App layout with sidebar/topbar
│   │   ├── login/page.tsx         # Login page
│   │   ├── layout.tsx             # Root layout
│   │   └── page.tsx               # Dashboard
│   │
│   ├── components/                # Reusable UI components (NEW)
│   │   ├── Button.tsx             # Variant-based button
│   │   ├── Modal.tsx              # Dialog component
│   │   ├── Form.tsx               # Form input suite
│   │   ├── Alert.tsx              # Alert messages
│   │   ├── Toast.tsx              # Notifications
│   │   ├── Table.tsx              # Data table
│   │   ├── Common.tsx             # Common patterns
│   │   ├── PasswordInput.tsx       # Password with strength meter
│   │   ├── Skeleton.tsx           # Loading skeletons
│   │   ├── PageHeader.tsx         # Page titles
│   │   ├── Sidebar.tsx            # Navigation sidebar
│   │   ├── Topbar.tsx             # Header bar
│   │   ├── StatCard.tsx           # Metric cards
│   │   ├── StatusBadge.tsx        # Status indicators
│   │   └── EmptyState.tsx         # Empty list state
│   │
│   ├── lib/
│   │   ├── api.ts                 # HTTP client with JWT
│   │   ├── auth-context.tsx       # Authentication state
│   │   ├── useApi.ts              # API hooks (NEW)
│   │   ├── types.ts               # TypeScript types
│   │   ├── utils.ts               # Formatting utilities
│   │   └── passwordValidator.ts   # Password validation (NEW)
│   │
│   └── styles/
│       └── globals.css            # Global styles & theme
│
├── next.config.js                 # Next.js configuration
├── tsconfig.json                  # TypeScript config
├── tailwind.config.js             # Tailwind CSS config
├── package.json                   # Dependencies
└── Dockerfile                     # Container config
```

### Component Library (NEW)

See [FRONTEND_COMPONENTS_REFERENCE.md](FRONTEND_COMPONENTS_REFERENCE.md) for detailed documentation.

**Core Components**:
- `Button` - Flexible button with 5 variants
- `Modal` - Dialog component with flexible layout
- `Form` - Input, Select, TextArea, Checkbox, FormField
- `Alert` - Dismissible alerts with 4 types
- `Toast` - Auto-dismissing notifications
- `Table` - Data table with sorting/pagination
- `PasswordInput` - Password with strength meter

**Patterns**:
- `ConfirmDialog` - Confirmation modal
- `FormActions` - Standard form buttons
- `DataRow` - Key-value display with copy
- `Stat` - Metric card with trends
- `EmptyState` - Empty list placeholder

### Utility Hooks (NEW)

- `useApi(url, options)` - Fetch data with loading/error
- `useApiMutation(options)` - POST/PUT/PATCH/DELETE
- `useAuth()` - Current user & auth methods
- `useToast()` - Show notifications

### Styling System

Uses CSS variables for theming:
```css
--accent: Primary action color (blue)
--accent-bright: Brighter shade
--accent-dark: Darker shade
--bg: Base background
--bg-elev: Elevated surface (cards)
--bg-elev-hover: Hover state
--text: Primary text
--text-dim: Secondary text
--text-mute: Muted text
--border: Border color
--border-strong: Stronger border
--danger: Error/destructive color (red)
```

All components use these variables for consistent theming.

---

## 🔐 Security & Deployment

### Security Checklist

- [x] JWT authentication with refresh tokens
- [x] Role-based access control (RBAC)
- [x] Rate limiting per IP address
- [x] CORS protection with configurable origins
- [x] Security headers (CSP, HSTS, X-Frame-Options)
- [x] Password hashing with bcrypt
- [x] NIST-compliant password requirements
- [x] Audit logging for sensitive operations
- [x] Database transaction isolation
- [x] Environment-based configuration
- [x] Input validation and sanitization
- [ ] DDoS protection (WAF recommended)
- [ ] SSL/TLS encryption (NGINX)
- [ ] Secrets management (use environment variables)

### Environment Variables

**Backend** (`.env`):
```
# Server
PORT=8000
DEBUG=false

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secure_password
DB_NAME=afritech

# JWT
JWT_SECRET=generate_with_openssl_rand_base64_32
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=720h

# Security
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=60

# Integrations
MPESA_API_URL=https://sandbox.safaricom.co.ke
MPESA_CONSUMER_KEY=...
MPESA_CONSUMER_SECRET=...
MPESA_SHORTCODE=...
MPESA_PASSKEY=...

MIKROTIK_HOST=...
MIKROTIK_USERNAME=...
MIKROTIK_PASSWORD=...

# Email
SMTP_HOST=...
SMTP_PORT=587
SMTP_USER=...
SMTP_PASSWORD=...

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://example.com
```

**Frontend** (`.env.local`):
```
NEXT_PUBLIC_API_URL=http://localhost:8000
NEXT_PUBLIC_APP_NAME=Afritech Online
```

### Deployment Checklist

- [ ] Generate strong JWT secret: `openssl rand -base64 32`
- [ ] Set production database credentials
- [ ] Configure HTTPS/SSL certificates
- [ ] Set CORS origins for production domain
- [ ] Enable rate limiting threshold
- [ ] Configure audit logging retention
- [ ] Set up database backups
- [ ] Configure monitoring and alerting
- [ ] Review security headers for production
- [ ] Test payment gateway credentials
- [ ] Set up MikroTik production connection
- [ ] Configure email service for alerts
- [ ] Run security audit
- [ ] Load test endpoints
- [ ] Plan disaster recovery

---

## 📖 API Reference

### Request Format

All requests to `/api/v1/*` must include JWT token:

```bash
curl -H "Authorization: Bearer {token}" \
  http://localhost:8000/api/v1/customers
```

### Response Format

All responses follow this structure:

```json
{
  "data": { ... },
  "error": null,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Handling

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not found |
| 422 | Validation error |
| 429 | Rate limited |
| 500 | Server error |

### Complete API Reference

See [SECURITY_FIXES_SUMMARY.md](SECURITY_FIXES_SUMMARY.md) for endpoint documentation.

---

## 🛠️ Development Guide

### Adding a New Feature

1. **Backend**:
   - Create model in `internal/models/models.go`
   - Create repository in `internal/feature/repository.go`
   - Create service in `internal/feature/service.go`
   - Create handlers in `internal/feature/handlers.go`
   - Register routes in `routes/routes.go`
   - Add tests in `internal/feature/service_test.go`

2. **Frontend**:
   - Add types in `lib/types.ts`
   - Create component in `components/`
   - Create page in `app/(app)/feature/`
   - Use `useApi` or `useApiMutation` hooks
   - Add error and loading states
   - Write unit tests

### Code Style

**Backend (Go)**:
- Use `fmt` for formatting
- Follow Go idioms and conventions
- Use structured logging
- Add error context with `%w` wrapping
- Write tests for business logic

**Frontend (TypeScript)**:
- Use strict TypeScript mode
- Export types from modules
- Use React hooks patterns
- Add PropTypes or TypeScript props
- Write tests for components

### Testing

**Backend**:
```bash
cd backend
go test ./...           # Run all tests
go test -v ./...        # Verbose output
go test -cover ./...    # Coverage report
go test ./... -bench .  # Benchmark tests
```

**Frontend**:
```bash
cd frontend
npm test                # Run jest tests
npm run type-check      # TypeScript check
npm run lint            # ESLint check
```

### Database Migrations

```bash
# Create migration
cd backend
go run cmd/server/main.go -create-migration "add_feature"

# Run migrations
go run cmd/server/main.go -migrate

# Rollback
go run cmd/server/main.go -rollback
```

---

## 📊 Project Status

### Completed
- ✅ Backend core architecture with Go + Gin
- ✅ MySQL database with migrations
- ✅ JWT authentication system
- ✅ CRUD operations for all entities
- ✅ Security middleware (rate limiting, headers)
- ✅ Audit logging system
- ✅ Password validation (NIST)
- ✅ Frontend Next.js setup
- ✅ Component library (8+ components)
- ✅ API integration layer
- ✅ Authentication context

### In Progress
- 🔄 Enhanced form pages and validation
- 🔄 Data table with sorting/filtering
- 🔄 Payment processing UI
- 🔄 M-Pesa integration frontend
- 🔄 Dashboard analytics

### Planned
- 📅 Admin dashboard with charts
- 📅 Customer self-service portal
- 📅 Email notifications
- 📅 SMS alerts
- 📅 Mobile app (React Native)

---

## 🤝 Contributing

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes following code style
3. Write/update tests
4. Commit with clear message: `git commit -m "feat: add new feature"`
5. Push to GitHub: `git push origin feature/my-feature`
6. Create Pull Request with description

---

## 📞 Support & Issues

For bug reports and feature requests, open an issue on GitHub.

---

## 📄 License

[Specify your license here]

---

## 🙏 Acknowledgments

Built with:
- Go & Gin Framework
- Next.js & React
- MySQL & TypeScript
- Tailwind CSS
- Open source community

---

**Last Updated**: January 2024  
**Version**: 1.0.0  
**Status**: Active Development
