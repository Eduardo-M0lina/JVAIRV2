package assigned_role

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/your-org/jvairv2/pkg/domain/role"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*AssignedRole, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AssignedRole), args.Error(1)
}

func (m *MockRepository) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*AssignedRole, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*AssignedRole), args.Error(1)
}

func (m *MockRepository) Assign(ctx context.Context, assignedRole *AssignedRole) error {
	args := m.Called(ctx, assignedRole)
	return args.Error(0)
}

func (m *MockRepository) Revoke(ctx context.Context, roleID, entityID int64, entityType string) error {
	args := m.Called(ctx, roleID, entityID, entityType)
	return args.Error(0)
}

func (m *MockRepository) HasRole(ctx context.Context, roleID, entityID int64, entityType string) (bool, error) {
	args := m.Called(ctx, roleID, entityID, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*AssignedRole, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*AssignedRole), args.Int(1), args.Error(2)
}

// MockRoleRepository es un mock de la interfaz role.Repository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id int64) (*role.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*role.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*role.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(ctx context.Context, role *role.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *role.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*role.Role, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*role.Role), args.Int(1), args.Error(2)
}
