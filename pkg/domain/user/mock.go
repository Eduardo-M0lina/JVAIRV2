package user

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/your-org/jvairv2/pkg/domain/ability"
	"github.com/your-org/jvairv2/pkg/domain/assigned_role"
	"github.com/your-org/jvairv2/pkg/domain/role"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*User, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*User), args.Int(1), args.Error(2)
}

func (m *MockRepository) VerifyCredentials(ctx context.Context, email, password string) (*User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetUserRoles(ctx context.Context, userID string) ([]*role.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*role.Role), args.Error(1)
}

func (m *MockRepository) GetUserAbilities(ctx context.Context, userID string) ([]*ability.Ability, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ability.Ability), args.Error(1)
}

// MockAssignedRoleRepository es un mock de la interfaz assigned_role.Repository
type MockAssignedRoleRepository struct {
	mock.Mock
}

func (m *MockAssignedRoleRepository) GetByID(ctx context.Context, id int64) (*assigned_role.AssignedRole, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*assigned_role.AssignedRole), args.Error(1)
}

func (m *MockAssignedRoleRepository) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*assigned_role.AssignedRole, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*assigned_role.AssignedRole), args.Error(1)
}

func (m *MockAssignedRoleRepository) Assign(ctx context.Context, assignedRole *assigned_role.AssignedRole) error {
	args := m.Called(ctx, assignedRole)
	return args.Error(0)
}

func (m *MockAssignedRoleRepository) Revoke(ctx context.Context, roleID, entityID int64, entityType string) error {
	args := m.Called(ctx, roleID, entityID, entityType)
	return args.Error(0)
}

func (m *MockAssignedRoleRepository) HasRole(ctx context.Context, roleID, entityID int64, entityType string) (bool, error) {
	args := m.Called(ctx, roleID, entityID, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockAssignedRoleRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*assigned_role.AssignedRole, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*assigned_role.AssignedRole), args.Int(1), args.Error(2)
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
