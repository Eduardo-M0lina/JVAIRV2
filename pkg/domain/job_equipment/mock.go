package job_equipment

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de equipos de trabajo
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, equipment *JobEquipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*JobEquipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JobEquipment), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, jobID int64, equipmentType string) ([]*JobEquipment, error) {
	args := m.Called(ctx, jobID, equipmentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*JobEquipment), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, equipment *JobEquipment) error {
	args := m.Called(ctx, equipment)
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

func (m *MockJobChecker) JobExists(id int64) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
