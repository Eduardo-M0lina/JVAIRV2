package workflow

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) List(ctx context.Context, filters Filters, page, pageSize int) ([]Workflow, int64, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]Workflow), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Workflow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Workflow), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, workflow *Workflow) error {
	args := m.Called(ctx, workflow)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, workflow *Workflow) error {
	args := m.Called(ctx, workflow)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Duplicate(ctx context.Context, id int64) (*Workflow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Workflow), args.Error(1)
}

func (m *MockRepository) GetWorkflowStatuses(ctx context.Context, workflowID int64) ([]WorkflowStatus, error) {
	args := m.Called(ctx, workflowID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]WorkflowStatus), args.Error(1)
}

func (m *MockRepository) SetWorkflowStatuses(ctx context.Context, workflowID int64, statuses []WorkflowStatus) error {
	args := m.Called(ctx, workflowID, statuses)
	return args.Error(0)
}
