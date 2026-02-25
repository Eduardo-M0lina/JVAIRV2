package job_priority

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, priority *JobPriority) error {
	args := m.Called(ctx, priority)
	return args.Error(0)
}

func (m *MockService) GetByID(ctx context.Context, id int64) (*JobPriority, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JobPriority), args.Error(1)
}

func (m *MockService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobPriority, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*JobPriority), args.Int(1), args.Error(2)
}

func (m *MockService) Update(ctx context.Context, priority *JobPriority) error {
	args := m.Called(ctx, priority)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
