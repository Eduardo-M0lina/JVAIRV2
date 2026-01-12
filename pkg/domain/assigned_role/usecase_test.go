package assigned_role

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/role"
)

func TestUseCase_GetByID(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	assignedRoleID := int64(1)
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1

	expectedAssignedRole := &AssignedRole{
		ID:         assignedRoleID,
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      &scope,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, assignedRoleID).Return(expectedAssignedRole, nil)

	// Ejecutar la función que estamos probando
	assignedRole, err := useCase.GetByID(ctx, assignedRoleID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignedRole, assignedRole)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByID_Error(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	assignedRoleID := int64(999)
	expectedError := errors.New("asignación de rol no encontrada")

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, assignedRoleID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	assignedRole, err := useCase.GetByID(ctx, assignedRoleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, assignedRole)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByEntity(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	entityType := "App\\Models\\User"
	entityID := int64(10)
	scope1 := 1
	scope2 := 2

	expectedAssignedRoles := []*AssignedRole{
		{
			ID:         1,
			RoleID:     2,
			EntityID:   entityID,
			EntityType: entityType,
			Restricted: false,
			Scope:      &scope1,
		},
		{
			ID:         2,
			RoleID:     3,
			EntityID:   entityID,
			EntityType: entityType,
			Restricted: true,
			Scope:      &scope2,
		},
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEntity", ctx, entityType, entityID).Return(expectedAssignedRoles, nil)

	// Ejecutar la función que estamos probando
	assignedRoles, err := useCase.GetByEntity(ctx, entityType, entityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignedRoles, assignedRoles)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Assign_Success(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1

	assignedRole := &AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      &scope,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRoleRepo.On("GetByID", ctx, roleID).Return(&role.Role{ID: roleID, Name: "admin"}, nil)
	mockRepo.On("Assign", ctx, assignedRole).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Assign(ctx, assignedRole)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRoleRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Assign_RoleNotFound(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := int64(999)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1
	expectedError := errors.New("rol no encontrado")

	assignedRole := &AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      &scope,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRoleRepo.On("GetByID", ctx, roleID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	err := useCase.Assign(ctx, assignedRole)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRoleRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Assign")
}

func TestUseCase_Revoke(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Revoke", ctx, roleID, entityID, entityType).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Revoke(ctx, roleID, entityID, entityType)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_HasRole(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar el comportamiento esperado del mock
	mockRepo.On("HasRole", ctx, roleID, entityID, entityType).Return(true, nil)

	// Ejecutar la función que estamos probando
	hasRole, err := useCase.HasRole(ctx, roleID, entityID, entityType)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.True(t, hasRole)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_List(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	filters := map[string]interface{}{"entity_type": "App\\Models\\User"}
	page := 1
	pageSize := 10

	scope1 := 1
	scope2 := 2

	expectedAssignedRoles := []*AssignedRole{
		{
			ID:         1,
			RoleID:     2,
			EntityID:   10,
			EntityType: "App\\Models\\User",
			Restricted: false,
			Scope:      &scope1,
		},
		{
			ID:         2,
			RoleID:     3,
			EntityID:   20,
			EntityType: "App\\Models\\User",
			Restricted: true,
			Scope:      &scope2,
		},
	}
	expectedTotal := 2

	// Configurar el comportamiento esperado del mock
	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedAssignedRoles, expectedTotal, nil)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := useCase.List(ctx, filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAssignedRoles, assignedRoles)
	assert.Equal(t, expectedTotal, total)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}
