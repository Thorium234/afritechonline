package customers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

// MockRepository is a test double for Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, c)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, search string, limit, offset int) ([]*models.Customer, int64, error) {
	args := m.Called(ctx, search, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Customer), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetByID(ctx context.Context, id uint64) (*models.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetByPhone(ctx context.Context, phone string) (*models.Customer, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, c)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestServiceListWithPagination(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	repo.On("List", mock.Anything, "", 20, 0).Return([]*models.Customer{}, int64(0), nil)

	customers, total, err := svc.ListWithPagination(context.Background(), "", 1, 20)
	assert.NoError(t, err)
	assert.Empty(t, customers)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

func TestServiceGet(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	expected := &models.Customer{ID: 1, FullName: "Alice", Phone: "254712345678", Username: "alice"}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

	c, err := svc.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", c.FullName)
	repo.AssertExpectations(t)
}

func TestServiceGetNotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	repo.On("GetByID", mock.Anything, uint64(99)).Return(nil, ErrNotFound)

	c, err := svc.Get(context.Background(), 99)
	assert.Error(t, err)
	assert.Nil(t, c)
	assert.Equal(t, ErrNotFound, err)
	repo.AssertExpectations(t)
}

func TestServiceCreate(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	input := &models.Customer{FullName: "Bob", Phone: "254712345679", Username: "bob"}
	created := &models.Customer{ID: 1, FullName: "Bob", Phone: "254712345679", Username: "bob"}
	repo.On("Create", mock.Anything, input).Return(created, nil)

	c, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), c.ID)
	repo.AssertExpectations(t)
}

func TestServiceCreateDefaultsStatus(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	input := &models.Customer{FullName: "Carol", Phone: "254712345680", Username: "carol"}
	created := &models.Customer{ID: 2, FullName: "Carol", Phone: "254712345680", Username: "carol", Status: "INACTIVE"}
	repo.On("Create", mock.Anything, input).Return(created, nil)

	c, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, "INACTIVE", c.Status)
	repo.AssertExpectations(t)
}

func TestServiceUpdate(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	existing := &models.Customer{ID: 1, FullName: "Old Name", Phone: "254712345678", Username: "old"}
	updated := &models.Customer{ID: 1, FullName: "New Name", Phone: "254712345678", Username: "old"}

	repo.On("GetByID", mock.Anything, uint64(1)).Return(existing, nil)
	repo.On("Update", mock.Anything, existing).Return(updated, nil)

	c, err := svc.Update(context.Background(), existing)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", c.FullName)
	repo.AssertExpectations(t)
}

func TestServiceDelete(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	repo.On("Delete", mock.Anything, uint64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestServiceDeleteNotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	repo.On("Delete", mock.Anything, uint64(99)).Return(ErrNotFound)

	err := svc.Delete(context.Background(), 99)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	repo.AssertExpectations(t)
}

func TestServiceSetActive(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	customer := &models.Customer{ID: 1, Status: "INACTIVE"}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(customer, nil)
	repo.On("Update", mock.Anything, customer).Return(customer, nil)

	err := svc.SetActive(context.Background(), 1, true)
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", customer.Status)
	repo.AssertExpectations(t)
}

func TestServiceSetActiveSuspended(t *testing.T) {
	repo := new(MockRepository)
	svc := NewService(repo)

	customer := &models.Customer{ID: 1, Status: "SUSPENDED"}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(customer, nil)

	err := svc.SetActive(context.Background(), 1, true)
	assert.NoError(t, err)
	assert.Equal(t, "SUSPENDED", customer.Status)
	repo.AssertExpectations(t)
}
