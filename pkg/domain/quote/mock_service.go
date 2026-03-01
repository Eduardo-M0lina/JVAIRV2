package quote

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, q *Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *MockService) GetByID(ctx context.Context, id int64) (*Quote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Quote), args.Error(1)
}

func (m *MockService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Quote, int64, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*Quote), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) Update(ctx context.Context, q *Quote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
