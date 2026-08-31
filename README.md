# Afritech Online

Afritech Online is an ISP management and internet billing platform.

## Quick Start

```bash
cp backend/.env.example backend/.env
# Edit backend/.env with your database credentials

docker compose up -d
```

## API Documentation

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user

### Customers
- `GET /api/v1/customers` - List customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/:id` - Get customer
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

### Packages
- `GET /api/v1/packages` - List packages
- `POST /api/v1/packages` - Create package
- `GET /api/v1/packages/:id` - Get package
- `PUT /api/v1/packages/:id` - Update package
- `DELETE /api/v1/packages/:id` - Delete package

### Subscriptions
- `GET /api/v1/subscriptions` - List subscriptions
- `POST /api/v1/subscriptions` - Create subscription
- `GET /api/v1/subscriptions/:id` - Get subscription

### Invoices
- `GET /api/v1/invoices` - List invoices
- `POST /api/v1/invoices` - Create invoice
- `GET /api/v1/invoices/:id` - Get invoice

### Payments
- `GET /api/v1/payments` - List payments
- `POST /api/v1/payments` - Create payment
- `GET /api/v1/payments/:id` - Get payment
- `POST /api/v1/payments/:id/complete` - Complete payment
- `POST /api/v1/payments/:id/fail` - Mark payment failed
- `POST /api/v1/payments/mpesa/stkpush` - Initiate M-Pesa STK Push
- `POST /api/v1/payments/mpesa/callback` - M-Pesa callback

### Routers (MikroTik)
- `GET /api/v1/routers` - List routers
- `POST /api/v1/routers` - Create router
- `GET /api/v1/routers/:id` - Get router
- `PUT /api/v1/routers/:id` - Update router
- `DELETE /api/v1/routers/:id` - Delete router
- `POST /api/v1/routers/:id/test` - Test connection
- `GET /api/v1/routers/:id/status` - Get router status

### Reports
- `GET /api/v1/reports/revenue?days=30` - Revenue summary
- `GET /api/v1/reports/customers` - Customer statistics
- `GET /api/v1/reports/routers` - Active routers count

## License

MIT
