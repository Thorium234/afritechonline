# Afritech Online

**Afritech Online** is an ISP management and internet billing platform designed to automate customer management, internet package subscriptions, payments, authentication, and MikroTik router management.

The project is being designed primarily for small and medium-sized Internet Service Providers (ISPs), Wi-Fi hotspot operators, and community networks that need a centralized way to manage customers and control internet access.

> **Project Status:** 🚧 Under active development
> **Current stage:** Architecture and system design

---

## Overview

Managing an ISP manually becomes difficult as the number of customers increases. Customer registration, package assignment, payments, internet activation, session tracking, and account expiration can quickly become disconnected processes.

Afritech Online aims to bring these operations into a single platform.

The intended system connects:

* Customer management
* Internet packages
* Billing
* M-Pesa payments
* MikroTik RouterOS
* FreeRADIUS
* Internet authentication
* Session management
* Account expiration
* Notifications
* Administrative reporting

The long-term goal is to allow an ISP administrator to manage the entire customer lifecycle from one system.

---

## Core Concept

The basic customer lifecycle is:

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
Activate Internet Account
   │
   ▼
Authenticate Through RADIUS
   │
   ▼
MikroTik Provides Internet Access
   │
   ▼
Subscription Expires
   │
   ▼
Account Suspended
```

This removes the need for administrators to manually activate and deactivate customers every time a payment is received or a package expires.

---

## Planned Architecture

```text
                         ┌───────────────────────┐
                         │     Customer Portal    │
                         │                       │
                         │ Login                 │
                         │ Packages              │
                         │ Payments              │
                         │ Account Status        │
                         └───────────┬───────────┘
                                     │
                                     ▼
                         ┌───────────────────────┐
                         │      Backend API      │
                         │                       │
                         │ Authentication        │
                         │ Customers             │
                         │ Packages              │
                         │ Subscriptions         │
                         │ Billing               │
                         │ Payments              │
                         │ Notifications         │
                         └───────┬────────┬──────┘
                                 │        │
                    ┌────────────┘        └─────────────┐
                    ▼                                   ▼
          ┌───────────────────┐              ┌───────────────────┐
          │     FreeRADIUS    │              │   Payment System  │
          │                   │              │                   │
          │ Authentication    │              │ M-Pesa            │
          │ Authorization     │              │ Other gateways    │
          │ Accounting        │              │                   │
          └─────────┬─────────┘              └───────────────────┘
                    │
                    ▼
          ┌───────────────────┐
          │      MikroTik     │
          │     RouterOS      │
          │                   │
          │ Hotspot / PPPoE   │
          │ Bandwidth Control │
          │ User Management   │
          └───────────────────┘
```

---

## Main Components

### Customer Management

Administrators will be able to:

* Register customers
* Update customer information
* View customer status
* Suspend customers
* Reactivate customers
* View subscription history
* View payment history

### Internet Packages

The system will support configurable packages containing information such as:

* Package name
* Price
* Duration
* Download speed
* Upload speed
* Data limits
* Active/inactive status

Example:

```text
Package: Home Basic
Price: KES 1,000
Duration: 30 days
Download: 10 Mbps
Upload: 5 Mbps
```

### Billing

Afritech Online will automatically calculate subscription charges and track:

* Pending payments
* Completed payments
* Failed payments
* Expired subscriptions
* Active subscriptions
* Payment history

### M-Pesa Integration

The platform is intended to support M-Pesa payments, including an STK Push workflow.

```text
Customer
   │
   ▼
Select Package
   │
   ▼
Enter M-Pesa Number
   │
   ▼
STK Push
   │
   ▼
Customer Confirms Payment
   │
   ▼
Payment Callback
   │
   ▼
Verify Transaction
   │
   ▼
Activate Subscription
```

Payment activation must only happen after the backend has successfully verified the payment.

### FreeRADIUS

FreeRADIUS is intended to handle authentication and authorization.

It can be used to manage:

* User credentials
* Authentication
* Session accounting
* Bandwidth profiles
* Access policies
* Session start and stop information

### MikroTik Integration

Afritech Online is designed to integrate with MikroTik RouterOS.

Potential management capabilities include:

* Creating users
* Removing users
* Enabling users
* Disabling users
* Assigning speed profiles
* Managing hotspot users
* Managing PPPoE users
* Monitoring active sessions
* Disconnecting sessions
* Applying access restrictions

The MikroTik integration will be separated from the core business logic so that router-specific functionality does not contaminate the rest of the application.

---

## Technology Stack

The exact production stack is still being finalized.

The architecture is expected to use technologies in the following areas:

| Component              | Technology                         |
| ---------------------- | ---------------------------------- |
| Frontend               | Web application                    |
| Backend                | REST API                           |
| Database               | PostgreSQL                         |
| Authentication         | Token/session based authentication |
| Network Authentication | FreeRADIUS                         |
| Network Device         | MikroTik RouterOS                  |
| Payments               | M-Pesa                             |
| Deployment             | Docker                             |
| Version Control        | Git + GitHub                       |

Technology choices may change as implementation progresses.

---

## Project Structure

The project is currently in the design phase. The planned structure is:

```text
afritechonline/
│
├── backend/
│   ├── auth/
│   ├── customers/
│   ├── packages/
│   ├── subscriptions/
│   ├── billing/
│   ├── payments/
│   ├── mikrotik/
│   ├── radius/
│   └── notifications/
│
├── frontend/
│   ├── components/
│   ├── dashboard/
│   ├── customer/
│   └── portal/
│
├── infrastructure/
│   ├── freeradius/
│   ├── mikrotik/
│   └── docker/
│
├── docs/
│   ├── architecture.md
│   ├── database.md
│   ├── api.md
│   └── deployment.md
│
├── README.md
├── .env.example
└── docker-compose.yml
```

---

## Development Roadmap

### Phase 1: System Design

* [x] Initial project idea
* [x] Initial architecture
* [x] Use case definition
* [x] MikroTik configuration research
* [ ] Finalize database design
* [ ] Finalize API architecture
* [ ] Define authentication architecture

### Phase 2: Backend

* [ ] Project initialization
* [ ] Database implementation
* [ ] User authentication
* [ ] Customer management
* [ ] Internet package management
* [ ] Subscription management
* [ ] Billing system
* [ ] Payment records
* [ ] M-Pesa integration
* [ ] Payment callbacks
* [ ] Subscription activation logic

### Phase 3: Network Integration

* [ ] FreeRADIUS setup
* [ ] RADIUS authentication
* [ ] RADIUS accounting
* [ ] MikroTik integration
* [ ] User provisioning
* [ ] Speed profile management
* [ ] Session management
* [ ] Automatic account suspension
* [ ] Automatic account activation

### Phase 4: Frontend

* [ ] Admin dashboard
* [ ] Customer dashboard
* [ ] Package management interface
* [ ] Customer management interface
* [ ] Billing interface
* [ ] Payment interface
* [ ] Network monitoring interface

### Phase 5: Production

* [ ] Docker deployment
* [ ] Environment configuration
* [ ] Database backups
* [ ] Logging
* [ ] Monitoring
* [ ] Security hardening
* [ ] Automated testing
* [ ] CI/CD
* [ ] Production deployment

---

## Security

Security is a core requirement because the platform will handle customer information, payment information, network credentials, and infrastructure access.

The production system must:

* Never store M-Pesa PINs
* Never expose MikroTik credentials to the frontend
* Keep secrets in environment variables
* Validate payment callbacks
* Authenticate administrative requests
* Authorize users based on roles
* Protect API endpoints
* Use HTTPS in production
* Validate all incoming data
* Implement rate limiting
* Log security-sensitive operations
* Keep database credentials out of source control

Example environment configuration:

```env
DATABASE_URL=
SECRET_KEY=

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

## Development Requirements

Before development begins, install:

* Git
* Docker
* Docker Compose
* PostgreSQL
* A supported backend runtime
* Node.js if the frontend uses a JavaScript framework

Clone the repository:

```bash
git clone https://github.com/Thorium234/afritechonline.git
cd afritechonline
```

The development installation instructions will be updated once the backend and frontend have been implemented.

---

## Network Architecture

A typical deployment is expected to look similar to:

```text
                    INTERNET
                       │
                       ▼
                ┌──────────────┐
                │    Router    │
                │   MikroTik   │
                └──────┬───────┘
                       │
          ┌────────────┴────────────┐
          │                         │
          ▼                         ▼
   Customer Traffic          Afritech Server
                              │
                  ┌───────────┼───────────┐
                  │           │           │
                  ▼           ▼           ▼
              Backend     PostgreSQL   FreeRADIUS
                  │
                  ▼
             M-Pesa API
```

The exact network topology will depend on the ISP's infrastructure.

---

## Important Design Principle

Afritech Online should **not** directly tie every part of the application to MikroTik.

The business logic should operate independently:

```text
Subscription
Payment
Customer
Package
Billing
     │
     ▼
Network Provisioning Service
     │
     ├── MikroTik
     ├── FreeRADIUS
     └── Future Network Providers
```

This makes the platform easier to extend to other networking equipment in the future.

---

## MVP

The first working version should deliberately remain small.

The MVP target is:

```text
1. Administrator creates an internet package
             ↓
2. Administrator/customer creates an account
             ↓
3. Customer selects a package
             ↓
4. Customer makes an M-Pesa payment
             ↓
5. Backend verifies payment
             ↓
6. Subscription becomes active
             ↓
7. Network account is provisioned
             ↓
8. Customer authenticates
             ↓
9. Customer receives internet
             ↓
10. Subscription expires
             ↓
11. Internet access is automatically disabled
```

If this flow works reliably, Afritech Online has a legitimate technical foundation.

---

## Future Features

Potential future functionality includes:

* Multiple MikroTik routers
* Multi-ISP support
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
* Network monitoring
* Router health monitoring
* Multiple locations
* Multi-tenant ISP management
* Mobile application

These features should only be added after the core billing and network provisioning workflow is stable.

---

## Contributing

Contributions are welcome once the project reaches the implementation stage.

For major changes:

1. Fork the repository.
2. Create a feature branch.

```bash
git checkout -b feature/my-feature
```

3. Make your changes.
4. Test the changes.
5. Commit your work.

```bash
git commit -m "feat: add customer management"
```

6. Push the branch.

```bash
git push origin feature/my-feature
```

7. Open a pull request.

---

## License

License information will be added when the project's licensing model is finalized.

---

## Project Status

Afritech Online is currently under development.

The repository currently contains the initial system concept, architecture, use cases, and MikroTik configuration research. The production application has not yet been completed.

The immediate priority is to build and validate the core:

**Customer → Package → Payment → Subscription → RADIUS → MikroTik → Internet Access**

Everything else comes after that.

---

## Repository

**GitHub:** https://github.com/Thorium234/afritechonline

**Project:** Afritech Online

**Purpose:** ISP management, billing, payment, authentication, and network automation.
