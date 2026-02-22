package job_category

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, category *JobCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*JobCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JobCategory), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobCategory, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*JobCategory), args.Int(1), args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, category *JobCategory) error {
	args := m.Called(ctx, category)
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
