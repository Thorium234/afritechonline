package invoices

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

type MockInvoiceRepo struct {
	mock.Mock
}

func (m *MockInvoiceRepo) NextNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockInvoiceRepo) Create(ctx context.Context, inv *models.Invoice) (*models.Invoice, error) {
	args := m.Called(ctx, inv)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Invoice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceRepo) GetByID(ctx context.Context, id uint64) (*models.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Invoice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceRepo) GetBySubscription(ctx context.Context, subID uint64) (*models.Invoice, error) {
	args := m.Called(ctx, subID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Invoice), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceRepo) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Invoice, int64, error) {
	args := m.Called(ctx, customerID, status, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Invoice), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *MockInvoiceRepo) MarkPaid(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSubscriptionProvider struct {
	mock.Mock
}

func (m *MockSubscriptionProvider) Get(ctx context.Context, id uint64) (*models.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestServiceCreateForSubscription(t *testing.T) {
	repo := new(MockInvoiceRepo)
	subSvc := new(MockSubscriptionProvider)
	svc := NewService(repo, subSvc)

	sub := &models.Subscription{ID: 1, CustomerID: 1, Amount: 1000, Currency: "KES", ExpiryDate: time.Now()}
	inv := &models.Invoice{ID: 1, InvoiceNo: "INV-20260101-000001", SubscriptionID: 1, CustomerID: 1, Status: StatusPending}

	subSvc.On("Get", mock.Anything, uint64(1)).Return(sub, nil)
	repo.On("NextNumber", mock.Anything).Return("INV-20260101-000001", nil)
	repo.On("Create", mock.Anything, mock.Anything).Return(inv, nil)

	got, err := svc.CreateForSubscription(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "INV-20260101-000001", got.InvoiceNo)
	assert.Equal(t, uint64(1), got.SubscriptionID)
	subSvc.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestServiceGet(t *testing.T) {
	repo := new(MockInvoiceRepo)
	subSvc := new(MockSubscriptionProvider)
	svc := NewService(repo, subSvc)

	inv := &models.Invoice{ID: 1, InvoiceNo: "INV-001", Status: StatusPending}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(inv, nil)

	got, err := svc.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "INV-001", got.InvoiceNo)
	repo.AssertExpectations(t)
}

func TestServiceGetNotFound(t *testing.T) {
	repo := new(MockInvoiceRepo)
	subSvc := new(MockSubscriptionProvider)
	svc := NewService(repo, subSvc)

	repo.On("GetByID", mock.Anything, uint64(99)).Return(nil, ErrNotFound)

	got, err := svc.Get(context.Background(), 99)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, ErrNotFound, err)
	repo.AssertExpectations(t)
}

func TestServiceMarkPaid(t *testing.T) {
	repo := new(MockInvoiceRepo)
	subSvc := new(MockSubscriptionProvider)
	svc := NewService(repo, subSvc)

	repo.On("MarkPaid", mock.Anything, uint64(1)).Return(nil)

	err := svc.MarkPaid(context.Background(), 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
