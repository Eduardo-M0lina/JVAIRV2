package technician_job_status

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, status *TechnicianJobStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*TechnicianJobStatus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TechnicianJobStatus), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*TechnicianJobStatus, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*TechnicianJobStatus), args.Int(1), args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, status *TechnicianJobStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
