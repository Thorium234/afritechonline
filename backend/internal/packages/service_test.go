package packages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Thorium234/afritechonline/backend/internal/models"
)

type MockPackageRepo struct {
	mock.Mock
}

func (m *MockPackageRepo) Create(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	args := m.Called(ctx, p)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InternetPackage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPackageRepo) List(ctx context.Context, activeOnly bool) ([]*models.InternetPackage, error) {
	args := m.Called(ctx, activeOnly)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.InternetPackage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPackageRepo) GetByID(ctx context.Context, id uint64) (*models.InternetPackage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InternetPackage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPackageRepo) Update(ctx context.Context, p *models.InternetPackage) (*models.InternetPackage, error) {
	args := m.Called(ctx, p)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InternetPackage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPackageRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestServiceList(t *testing.T) {
	repo := new(MockPackageRepo)
	svc := NewService(repo)

	pkgs := []*models.InternetPackage{
		{ID: 1, Name: "Basic", Price: 1000, IsActive: true},
	}
	repo.On("List", mock.Anything, false).Return(pkgs, nil)

	items, err := svc.List(context.Background(), false)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	repo.AssertExpectations(t)
}

func TestServiceGet(t *testing.T) {
	repo := new(MockPackageRepo)
	svc := NewService(repo)

	pkg := &models.InternetPackage{ID: 1, Name: "Basic", Price: 1000, IsActive: true}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(pkg, nil)

	got, err := svc.Get(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Basic", got.Name)
	repo.AssertExpectations(t)
}

func TestServiceGetNotFound(t *testing.T) {
	repo := new(MockPackageRepo)
	svc := NewService(repo)

	repo.On("GetByID", mock.Anything, uint64(99)).Return(nil, ErrNotFound)

	got, err := svc.Get(context.Background(), 99)
	assert.Error(t, err)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestServiceGetActive(t *testing.T) {
	repo := new(MockPackageRepo)
	svc := NewService(repo)

	pkg := &models.InternetPackage{ID: 1, Name: "Basic", Price: 1000, IsActive: true}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(pkg, nil)

	got, err := svc.GetActive(context.Background(), 1)
	assert.NoError(t, err)
	assert.True(t, got.IsActive)
	repo.AssertExpectations(t)
}

func TestServiceGetActiveInactive(t *testing.T) {
	repo := new(MockPackageRepo)
	svc := NewService(repo)

	pkg := &models.InternetPackage{ID: 1, Name: "Old", Price: 1000, IsActive: false}
	repo.On("GetByID", mock.Anything, uint64(1)).Return(pkg, nil)

	got, err := svc.GetActive(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, ErrNotFound, err)
	repo.AssertExpectations(t)
}
