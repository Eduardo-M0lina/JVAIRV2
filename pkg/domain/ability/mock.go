package ability

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Ability, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Ability), args.Error(1)
}

func (m *MockRepository) GetByName(ctx context.Context, name string) (*Ability, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Ability), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, ability *Ability) error {
	args := m.Called(ctx, ability)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, ability *Ability) error {
	args := m.Called(ctx, ability)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Ability, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*Ability), args.Int(1), args.Error(2)
}
