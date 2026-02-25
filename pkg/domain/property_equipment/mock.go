package property_equipment

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock del repositorio de equipos de propiedad
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, equipment *PropertyEquipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*PropertyEquipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PropertyEquipment), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, propertyID int64) ([]*PropertyEquipment, error) {
	args := m.Called(ctx, propertyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*PropertyEquipment), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, equipment *PropertyEquipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
