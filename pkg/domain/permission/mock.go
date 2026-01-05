package permission

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/your-org/jvairv2/pkg/domain/ability"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Permission), args.Error(1)
}

func (m *MockRepository) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*Permission, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Permission), args.Error(1)
}

func (m *MockRepository) GetByAbility(ctx context.Context, abilityID int64) ([]*Permission, error) {
	args := m.Called(ctx, abilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Permission), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, permission *Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, permission *Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Exists(ctx context.Context, abilityID, entityID int64, entityType string) (bool, error) {
	args := m.Called(ctx, abilityID, entityID, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Permission, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*Permission), args.Int(1), args.Error(2)
}

// MockAbilityRepository es un mock de la interfaz ability.Repository
type MockAbilityRepository struct {
	mock.Mock
}

func (m *MockAbilityRepository) GetByID(ctx context.Context, id int64) (*ability.Ability, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ability.Ability), args.Error(1)
}

func (m *MockAbilityRepository) GetByName(ctx context.Context, name string) (*ability.Ability, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ability.Ability), args.Error(1)
}

func (m *MockAbilityRepository) Create(ctx context.Context, ability *ability.Ability) error {
	args := m.Called(ctx, ability)
	return args.Error(0)
}

func (m *MockAbilityRepository) Update(ctx context.Context, ability *ability.Ability) error {
	args := m.Called(ctx, ability)
	return args.Error(0)
}

func (m *MockAbilityRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAbilityRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*ability.Ability, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*ability.Ability), args.Int(1), args.Error(2)
}
