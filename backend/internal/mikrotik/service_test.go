package mikrotik

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

type MockRouterRepo struct {
	mock.Mock
}

func (m *MockRouterRepo) Create(ctx context.Context, r *models.Router) (*models.Router, error) {
	args := m.Called(ctx, r)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Router), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRouterRepo) GetByID(ctx context.Context, id uint64) (*models.Router, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Router), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRouterRepo) List(ctx context.Context) ([]*models.Router, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Router), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRouterRepo) Update(ctx context.Context, r *models.Router) (*models.Router, error) {
	args := m.Called(ctx, r)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Router), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRouterRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRouterRepo) UpdateStatus(ctx context.Context, id uint64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestServiceRegister(t *testing.T) {
	repo := new(MockRouterRepo)
	svc := NewService(repo, nil)

	router := &models.Router{Name: "Main", Host: "192.168.1.1", Username: "admin", PasswordEnc: "secret"}
	repo.On("Create", mock.Anything, router).Return(router, nil)

	created, err := svc.Register(context.Background(), router)
	assert.NoError(t, err)
	assert.Equal(t, "Main", created.Name)
	assert.Equal(t, 8728, created.APIPort)
	repo.AssertExpectations(t)
}

func TestServiceGet(t *testing.T) {
	repo := new(MockRouterRepo)
	svc := NewService(repo, nil)

	router := &models.Router{ID: 1, Name: "Main", Host: "192.168.1.1"}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(router, nil)

	got, err := svc.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Main", got.Name)
	repo.AssertExpectations(t)
}

func TestServiceList(t *testing.T) {
	repo := new(MockRouterRepo)
	svc := NewService(repo, nil)

	routers := []*models.Router{{ID: 1, Name: "Main"}}
	repo.On("List", mock.Anything).Return(routers, nil)

	got, err := svc.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestServiceDelete(t *testing.T) {
	repo := new(MockRouterRepo)
	svc := NewService(repo, nil)

	repo.On("Delete", mock.Anything, uint64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestServiceDefaultAPIPort(t *testing.T) {
	assert.Equal(t, 8728, defaultAPIPort(0))
	assert.Equal(t, 8728, defaultAPIPort(-1))
	assert.Equal(t, 8080, defaultAPIPort(8080))
}
