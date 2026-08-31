package payments

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, p *models.Payment) (*models.Payment, error) {
	args := m.Called(ctx, p)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id uint64) (*models.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepo) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Payment, int64, error) {
	args := m.Called(ctx, customerID, status, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Payment), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepo) MarkCompleted(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepo) MarkFailed(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepo) Begin(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*sql.Tx), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentRepo) WithTx(tx *sql.Tx) *Repository {
	args := m.Called(tx)
	if args.Get(0) != nil {
		return args.Get(0).(*Repository)
	}
	return nil
}

type MockInvoiceService struct {
	mock.Mock
}

func (m *MockInvoiceService) Get(ctx context.Context, id uint64) (*models.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Invoice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceService) MarkPaid(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Get(ctx context.Context, id uint64) (*models.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionService) Activate(ctx context.Context, subID uint64) (*models.Subscription, error) {
	args := m.Called(ctx, subID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestServiceComplete(t *testing.T) {
	repo := new(MockPaymentRepo)
	invSvc := new(MockInvoiceService)
	subSvc := new(MockSubscriptionService)
	svc := NewService(repo, invSvc, subSvc)

	payment := &models.Payment{ID: 1, InvoiceID: 10, CustomerID: 1, Status: StatusPending}
	invoice := &models.Invoice{ID: 10, SubscriptionID: 20, Status: StatusPending}
	subscription := &models.Subscription{ID: 20, CustomerID: 1, Status: StatusPending}

	repo.On("GetByID", mock.Anything, uint64(1)).Return(payment, nil)
	repo.On("Begin", mock.Anything).Return(&sql.Tx{}, nil)
	repo.On("WithTx", mock.Anything).Return(repo)
	repo.On("MarkCompleted", mock.Anything, uint64(1)).Return(nil)
	invSvc.On("Get", mock.Anything, uint64(10)).Return(invoice, nil)
	invSvc.On("MarkPaid", mock.Anything, uint64(10)).Return(nil)
	subSvc.On("Activate", mock.Anything, uint64(20)).Return(subscription, nil)
	repo.On("GetByID", mock.Anything, uint64(1)).Return(payment, nil)

	got, err := svc.Complete(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, StatusPending, got.Status) // original object not modified by mock
	repo.AssertExpectations(t)
	invSvc.AssertExpectations(t)
	subSvc.AssertExpectations(t)
}

func TestServiceCompleteIdempotent(t *testing.T) {
	repo := new(MockPaymentRepo)
	invSvc := new(MockInvoiceService)
	subSvc := new(MockSubscriptionService)
	svc := NewService(repo, invSvc, subSvc)

	payment := &models.Payment{ID: 1, InvoiceID: 10, CustomerID: 1, Status: StatusCompleted}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(payment, nil)

	got, err := svc.Complete(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, got.Status)
	repo.AssertExpectations(t)
}

func TestServiceCompleteInvalidState(t *testing.T) {
	repo := new(MockPaymentRepo)
	invSvc := new(MockInvoiceService)
	subSvc := new(MockSubscriptionService)
	svc := NewService(repo, invSvc, subSvc)

	payment := &models.Payment{ID: 1, InvoiceID: 10, CustomerID: 1, Status: StatusFailed}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(payment, nil)

	_, err := svc.Complete(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidState, err)
	repo.AssertExpectations(t)
}

func TestServiceCreate(t *testing.T) {
	repo := new(MockPaymentRepo)
	invSvc := new(MockInvoiceService)
	subSvc := new(MockSubscriptionService)
	svc := NewService(repo, invSvc, subSvc)

	payment := &models.Payment{ID: 1, InvoiceID: 10, CustomerID: 1, Status: StatusPending}
	repo.On("Create", mock.Anything, payment).Return(payment, nil)

	got, err := svc.Create(context.Background(), payment)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), got.ID)
	repo.AssertExpectations(t)
}

func TestServiceGet(t *testing.T) {
	repo := new(MockPaymentRepo)
	invSvc := new(MockInvoiceService)
	subSvc := new(MockSubscriptionService)
	svc := NewService(repo, invSvc, subSvc)

	payment := &models.Payment{ID: 1, InvoiceID: 10, CustomerID: 1, Status: StatusPending}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(payment, nil)

	got, err := svc.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), got.ID)
	repo.AssertExpectations(t)
}
