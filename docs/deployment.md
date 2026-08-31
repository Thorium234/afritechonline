# Afritech Online Deployment Guide

## Prerequisites

- Docker and Docker Compose
- MySQL/MariaDB database
- M-Pesa developer account (for production)

## Environment Configuration

1. Copy `.env.example` to `.env`:
   ```bash
   cp backend/.env.example backend/.env
   ```

2. Configure the following variables:
   - `JWT_SECRET` - A strong random string for JWT signing
   - `DB_*` - Database connection details
   - `MPESA_*` - M-Pesa API credentials (sandbox or production)
   - `MIKROTIK_*` - MikroTik router credentials (optional)

## Deployment Steps

### 1. Build and Start Services

```bash
docker compose up -d --build
```

### 2. Run Migrations

Migrations run automatically on startup.

### 3. Seed Initial Data

```bash
docker compose exec backend go run ./cmd/server
```

The server will seed default packages and an admin user if `SEED_ADMIN_PASSWORD` is set.

### 4. Access the Application

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

## Production Checklist

- [ ] Set strong `JWT_SECRET`
- [ ] Configure `APP_ENV=production`
- [ ] Set `APP_BASE_URL` to your domain
- [ ] Configure M-Pesa production credentials
- [ ] Set up database backups (see `scripts/backup.sh`)
- [ ] Configure HTTPS reverse proxy (nginx)
- [ ] Set up monitoring and logging
- [ ] Review firewall rules
- [ ] Enable database encryption at rest
