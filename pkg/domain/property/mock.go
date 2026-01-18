package property

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de propiedades
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, property *Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Property, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Property), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Property, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*Property), args.Int(1), args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, property *Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) HasJobs(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
