# Afritech Online API Documentation

## Base URL

```
/api/v1
```

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "status": 400,
    "message": "Error description"
  }
}
```

## Endpoints

### POST /auth/register

Register a new user.

**Request:**
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "securepassword"
}
```

**Response:** 201 Created
```json
{
  "data": {
    "user": { "id": 1, "username": "admin", "email": "admin@example.com", "role": "CUSTOMER" },
    "tokens": { "access_token": "...", "refresh_token": "...", "token_type": "Bearer", "expires_in": 900 }
  }
}
```

### POST /auth/login

Authenticate and receive tokens.

**Request:**
```json
{
  "identifier": "admin",
  "password": "securepassword"
}
```

### POST /auth/refresh

Refresh access token.

**Request:**
```json
{
  "refresh_token": "..."
}
```

### GET /auth/me

Get current authenticated user.

### Customers

- `GET /customers` - List (staff+)
- `POST /customers` - Create (staff+)
- `GET /customers/:id` - Get (staff+)
- `PUT /customers/:id` - Update (staff+)
- `DELETE /customers/:id` - Delete (admin+)

### Packages

- `GET /packages` - List (public)
- `POST /packages` - Create (staff+)
- `GET /packages/:id` - Get (public)
- `PUT /packages/:id` - Update (staff+)
- `DELETE /packages/:id` - Delete (admin+)

### Subscriptions

- `GET /subscriptions` - List (protected)
- `POST /subscriptions` - Create (staff+)
- `GET /subscriptions/:id` - Get (protected)

### Invoices

- `GET /invoices` - List (protected)
- `POST /invoices` - Create (staff+)
- `GET /invoices/:id` - Get (protected)

### Payments

- `GET /payments` - List (protected)
- `POST /payments` - Create (staff+)
- `GET /payments/:id` - Get (protected)
- `POST /payments/:id/complete` - Complete (staff+)
- `POST /payments/:id/fail` - Mark failed (staff+)

### M-Pesa

- `POST /payments/mpesa/stkpush` - Initiate STK Push (staff+)
- `POST /payments/mpesa/callback` - Callback endpoint (public)

### Routers

- `GET /routers` - List (staff+)
- `POST /routers` - Create (staff+)
- `GET /routers/:id` - Get (staff+)
- `PUT /routers/:id` - Update (staff+)
- `DELETE /routers/:id` - Delete (admin+)
- `POST /routers/:id/test` - Test connection (staff+)
- `GET /routers/:id/status` - Get status (staff+)

### Reports

- `GET /reports/revenue?days=30` - Revenue summary (staff+)
- `GET /reports/customers` - Customer stats (staff+)
- `GET /reports/routers` - Active routers (staff+)
