package quote

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de cotizaciones
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, q *Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Quote), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Quote, int64, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*Quote), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, q *Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockJobChecker es un mock del checker de jobs
type MockJobChecker struct {
	mock.Mock
}

func (m *MockJobChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// MockQuoteStatusChecker es un mock del checker de estados de cotización
type MockQuoteStatusChecker struct {
	mock.Mock
}

func (m *MockQuoteStatusChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}
