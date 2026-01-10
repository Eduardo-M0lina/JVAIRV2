package permission

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/ability"
)

func TestUseCase_GetByID(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	now := time.Now()

	expectedPermission := &Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, permissionID).Return(expectedPermission, nil)

	// Ejecutar la función que estamos probando
	permission, err := useCase.GetByID(ctx, permissionID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedPermission, permission)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByID_Error(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	permissionID := int64(999)
	expectedError := errors.New("permiso no encontrado")

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, permissionID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	permission, err := useCase.GetByID(ctx, permissionID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, permission)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByEntity(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	entityType := "App\\Models\\User"
	entityID := int64(10)
	now := time.Now()
	conditions1 := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	expectedPermissions := []*Permission{
		{
			ID:         1,
			AbilityID:  2,
			EntityID:   entityID,
			EntityType: entityType,
			Forbidden:  false,
			Conditions: &conditions1,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
		{
			ID:         2,
			AbilityID:  3,
			EntityID:   entityID,
			EntityType: entityType,
			Forbidden:  true,
			Conditions: nil,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEntity", ctx, entityType, entityID).Return(expectedPermissions, nil)

	// Ejecutar la función que estamos probando
	permissions, err := useCase.GetByEntity(ctx, entityType, entityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByAbility(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(2)
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	expectedPermissions := []*Permission{
		{
			ID:         1,
			AbilityID:  abilityID,
			EntityID:   10,
			EntityType: "App\\Models\\User",
			Forbidden:  false,
			Conditions: &conditions,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
		{
			ID:         2,
			AbilityID:  abilityID,
			EntityID:   5,
			EntityType: "App\\Models\\Role",
			Forbidden:  true,
			Conditions: nil,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByAbility", ctx, abilityID).Return(expectedPermissions, nil)

	// Ejecutar la función que estamos probando
	permissions, err := useCase.GetByAbility(ctx, abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create_Success(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	permission := &Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar el comportamiento esperado de los mocks
	mockAbilityRepo.On("GetByID", ctx, abilityID).Return(&ability.Ability{ID: abilityID, Name: "create_user"}, nil)
	mockRepo.On("Create", ctx, permission).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, permission)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockAbilityRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create_AbilityNotFound(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(999)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	expectedError := errors.New("ability no encontrada")

	permission := &Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar el comportamiento esperado de los mocks
	mockAbilityRepo.On("GetByID", ctx, abilityID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, permission)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockAbilityRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUseCase_Update_Success(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	permission := &Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar el comportamiento esperado de los mocks
	mockAbilityRepo.On("GetByID", ctx, abilityID).Return(&ability.Ability{ID: abilityID, Name: "create_user"}, nil)
	mockRepo.On("Update", ctx, permission).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Update(ctx, permission)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockAbilityRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Delete(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	permissionID := int64(1)

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Delete", ctx, permissionID).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Delete(ctx, permissionID)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Exists(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Exists", ctx, abilityID, entityID, entityType).Return(true, nil)

	// Ejecutar la función que estamos probando
	exists, err := useCase.Exists(ctx, abilityID, entityID, entityType)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.True(t, exists)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_List(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAbilityRepo := new(MockAbilityRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAbilityRepo)

	// Datos de prueba
	ctx := context.Background()
	filters := map[string]interface{}{"entity_type": "App\\Models\\User"}
	page := 1
	pageSize := 10

	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	expectedPermissions := []*Permission{
		{
			ID:         1,
			AbilityID:  2,
			EntityID:   10,
			EntityType: "App\\Models\\User",
			Forbidden:  false,
			Conditions: &conditions,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
		{
			ID:         2,
			AbilityID:  3,
			EntityID:   20,
			EntityType: "App\\Models\\User",
			Forbidden:  true,
			Conditions: nil,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
	}
	expectedTotal := 2

	// Configurar el comportamiento esperado del mock
	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedPermissions, expectedTotal, nil)

	// Ejecutar la función que estamos probando
	permissions, total, err := useCase.List(ctx, filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)
	assert.Equal(t, expectedTotal, total)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}
