package job

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de jobs
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, j *Job) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Job), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Job, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*Job), args.Int(1), args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, j *Job) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Close(ctx context.Context, id int64, jobStatusID int64) error {
	args := m.Called(ctx, id, jobStatusID)
	return args.Error(0)
}

// MockJobCategoryChecker es un mock del checker de categorías
type MockJobCategoryChecker struct {
	mock.Mock
}

func (m *MockJobCategoryChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// MockJobPriorityChecker es un mock del checker de prioridades
type MockJobPriorityChecker struct {
	mock.Mock
}

func (m *MockJobPriorityChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// MockJobStatusChecker es un mock del checker de estados
type MockJobStatusChecker struct {
	mock.Mock
}

func (m *MockJobStatusChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// MockWorkflowChecker es un mock del checker de workflows
type MockWorkflowChecker struct {
	mock.Mock
}

func (m *MockWorkflowChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

func (m *MockWorkflowChecker) GetInitialStatusID(ctx context.Context, workflowID int64) (int64, error) {
	args := m.Called(ctx, workflowID)
	return args.Get(0).(int64), args.Error(1)
}

// MockPropertyChecker es un mock del checker de propiedades
type MockPropertyChecker struct {
	mock.Mock
}

func (m *MockPropertyChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

func (m *MockPropertyChecker) GetWorkflowID(ctx context.Context, propertyID int64) (int64, error) {
	args := m.Called(ctx, propertyID)
	return args.Get(0).(int64), args.Error(1)
}

// MockUserChecker es un mock del checker de usuarios
type MockUserChecker struct {
	mock.Mock
}

func (m *MockUserChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

// MockTechnicianJobStatusChecker es un mock del checker de estados de técnico
type MockTechnicianJobStatusChecker struct {
	mock.Mock
}

func (m *MockTechnicianJobStatusChecker) GetByID(ctx context.Context, id int64) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

func (m *MockTechnicianJobStatusChecker) GetLinkedJobStatusID(ctx context.Context, id int64) (*int64, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}
