package job_category

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, category *JobCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockService) GetByID(ctx context.Context, id int64) (*JobCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JobCategory), args.Error(1)
}

func (m *MockService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobCategory, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*JobCategory), args.Int(1), args.Error(2)
}

func (m *MockService) Update(ctx context.Context, category *JobCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
