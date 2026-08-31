package subscriptions

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

type MockSubRepo struct {
	mock.Mock
}

func (m *MockSubRepo) Create(ctx context.Context, s *models.Subscription) (*models.Subscription, error) {
	args := m.Called(ctx, s)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubRepo) GetByID(ctx context.Context, id uint64) (*models.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubRepo) List(ctx context.Context, customerID uint64, status string, limit, offset int) ([]*models.Subscription, int64, error) {
	args := m.Called(ctx, customerID, status, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Subscription), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *MockSubRepo) UpdateStatus(ctx context.Context, id uint64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSubRepo) GetActiveForCustomer(ctx context.Context, customerID uint64) (*models.Subscription, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubRepo) ListExpired(ctx context.Context, cutoff interface{}) ([]uint64, error) {
	args := m.Called(ctx, cutoff)
	if args.Get(0) != nil {
		return args.Get(0).([]uint64), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockPackageProvider struct {
	mock.Mock
}

func (m *MockPackageProvider) GetActive(ctx context.Context, id uint64) (*models.InternetPackage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InternetPackage), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockCustomerProvider struct {
	mock.Mock
}

func (m *MockCustomerProvider) Get(ctx context.Context, id uint64) (*models.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCustomerProvider) SetActive(ctx context.Context, id uint64, active bool) error {
	args := m.Called(ctx, id, active)
	return args.Error(0)
}

func TestServiceCreate(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	pkg := &models.InternetPackage{ID: 1, Name: "Basic", Price: 1000, DurationDays: 30, Currency: "KES"}
	pkgSvc.On("GetActive", mock.Anything, uint64(1)).Return(pkg, nil)
	custSvc.On("Get", mock.Anything, uint64(1)).Return(&models.Customer{ID: 1}, nil)
	repo.On("Create", mock.Anything, mock.Anything).Return(&models.Subscription{ID: 1, CustomerID: 1, PackageID: 1}, nil)

	sub, err := svc.Create(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), sub.ID)
	assert.Equal(t, StatusPending, sub.Status)
	pkgSvc.AssertExpectations(t)
	custSvc.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestServiceActivate(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	sub := &models.Subscription{ID: 1, CustomerID: 1, Status: StatusPending}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)
	repo.On("UpdateStatus", mock.Anything, uint64(1), StatusActive).Return(nil)
	custSvc.On("SetActive", mock.Anything, uint64(1), true).Return(nil)
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)

	got, err := svc.Activate(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, StatusActive, got.Status)
	repo.AssertExpectations(t)
	custSvc.AssertExpectations(t)
}

func TestServiceActivateIdempotent(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	sub := &models.Subscription{ID: 1, CustomerID: 1, Status: StatusActive}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)

	got, err := svc.Activate(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, StatusActive, got.Status)
	repo.AssertExpectations(t)
}

func TestServiceActivateInvalidState(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	sub := &models.Subscription{ID: 1, CustomerID: 1, Status: StatusExpired}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)

	_, err := svc.Activate(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidState, err)
	repo.AssertExpectations(t)
}

func TestServiceExpire(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	sub := &models.Subscription{ID: 1, CustomerID: 1, Status: StatusActive}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)
	repo.On("UpdateStatus", mock.Anything, uint64(1), StatusExpired).Return(nil)
	custSvc.On("SetActive", mock.Anything, uint64(1), false).Return(nil)
	repo.On("GetByID", mock.Anything, uint64(1)).Return(sub, nil)

	got, err := svc.Expire(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, StatusExpired, got.Status)
	repo.AssertExpectations(t)
	custSvc.AssertExpectations(t)
}

func TestServiceRefreshExpiredStatus(t *testing.T) {
	repo := new(MockSubRepo)
	pkgSvc := new(MockPackageProvider)
	custSvc := new(MockCustomerProvider)
	svc := NewService(repo, pkgSvc, custSvc)

	cutoff := time.Now()
	repo.On("ListExpired", mock.Anything, cutoff).Return([]uint64{1, 2}, nil)
	repo.On("UpdateStatus", mock.Anything, uint64(1), StatusExpired).Return(nil)
	repo.On("UpdateStatus", mock.Anything, uint64(2), StatusExpired).Return(nil)

	count, err := svc.RefreshExpiredStatus(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	repo.AssertExpectations(t)
}
