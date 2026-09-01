package mpesa

import (
	"context"
	"fmt"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/internal/payments"
)

// Service orchestrates M-Pesa operations.
type Service struct {
	client     *Client
	payments   PaymentRepository
	invoices   InvoiceRepository
	subs       SubscriptionRepository
}

// PaymentRepository is a minimal interface for payment operations.
type PaymentRepository interface {
	Create(ctx context.Context, p *models.Payment) (*models.Payment, error)
	GetByReference(ctx context.Context, reference, method string) (*models.Payment, error)
	MarkCompleted(ctx context.Context, id uint64) error
	MarkFailed(ctx context.Context, id uint64) error
}

// InvoiceRepository is a minimal interface for invoice operations.
type InvoiceRepository interface {
	GetByID(ctx context.Context, id uint64) (*models.Invoice, error)
	MarkPaid(ctx context.Context, id uint64) error
}

// SubscriptionRepository is a minimal interface for subscription operations.
type SubscriptionRepository interface {
	Get(ctx context.Context, id uint64) (*models.Subscription, error)
	Activate(ctx context.Context, subID uint64) (*models.Subscription, error)
}

// NewService creates an M-Pesa service.
func NewService(client *Client, payments PaymentRepository, invoices InvoiceRepository, subs SubscriptionRepository) *Service {
	return &Service{client: client, payments: payments, invoices: invoices, subs: subs}
}

// STKPush initiates an STK Push payment request.
func (s *Service) STKPush(ctx context.Context, invoiceID uint64, phoneNumber string, amount float64) (*STKPushResponse, error) {
	req := STKPushRequest{
		Amount:           amount,
		PhoneNumber:      phoneNumber,
		AccountReference: fmt.Sprintf("INV-%d", invoiceID),
		TransactionDesc:  "Internet subscription payment",
	}
	return s.client.InitiateSTKPush(ctx, req)
}

// ProcessCallback handles M-Pesa payment callbacks idempotently.
func (s *Service) ProcessCallback(ctx context.Context, body []byte) error {
	payment, err := s.client.ParseCallback(body)
	if err != nil {
		return err
	}

	// Find the payment by reference
	p, err := s.payments.GetByReference(ctx, payment.Reference, payments.PaymentMethodMPesa)
	if err != nil {
		return err
	}

	if p.Status == payments.StatusCompleted {
		return nil // idempotent
	}

	if payment.Status == payments.StatusCompleted {
		if err := s.payments.MarkCompleted(ctx, p.ID); err != nil {
			return err
		}

		inv, err := s.invoices.GetByID(ctx, p.InvoiceID)
		if err != nil {
			return err
		}
		if inv.Status != "PAID" {
			if err := s.invoices.MarkPaid(ctx, inv.ID); err != nil {
				return err
			}
		}
		if _, err := s.subs.Activate(ctx, inv.SubscriptionID); err != nil {
			return err
		}
	} else {
		if err := s.payments.MarkFailed(ctx, p.ID); err != nil {
			return err
		}
	}

	return nil
}
