package ability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_GetByID(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(1)
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	now := time.Now()

	expectedAbility := &Ability{
		ID:         abilityID,
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, abilityID).Return(expectedAbility, nil)

	// Ejecutar la función que estamos probando
	ability, err := useCase.GetByID(ctx, abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAbility, ability)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByID_Error(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(999)
	expectedError := errors.New("ability no encontrada")

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, abilityID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	ability, err := useCase.GetByID(ctx, abilityID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, ability)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByName(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityName := "create_user"
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	now := time.Now()

	expectedAbility := &Ability{
		ID:         1,
		Name:       abilityName,
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByName", ctx, abilityName).Return(expectedAbility, nil)

	// Ejecutar la función que estamos probando
	ability, err := useCase.GetByName(ctx, abilityName)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAbility, ability)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	ability := &Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Create", ctx, ability).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, ability)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	ability := &Ability{
		ID:         1,
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Update", ctx, ability).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Update(ctx, ability)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Delete(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	abilityID := int64(1)

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Delete", ctx, abilityID).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Delete(ctx, abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_List(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockRepository)

	// Crear el caso de uso con el mock
	useCase := NewUseCase(mockRepo)

	// Datos de prueba
	ctx := context.Background()
	filters := map[string]interface{}{"name": "create"}
	page := 1
	pageSize := 10

	title1 := "Create User"
	title2 := "Edit User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope1 := 1
	scope2 := 2
	entityID := int64(10)
	now := time.Now()

	expectedAbilities := []*Ability{
		{
			ID:         1,
			Name:       "create_user",
			Title:      &title1,
			EntityID:   &entityID,
			EntityType: &entityType,
			OnlyOwned:  false,
			Options:    &options,
			Scope:      &scope1,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
		{
			ID:         2,
			Name:       "edit_user",
			Title:      &title2,
			EntityID:   &entityID,
			EntityType: &entityType,
			OnlyOwned:  false,
			Options:    &options,
			Scope:      &scope2,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		},
	}
	expectedTotal := 2

	// Configurar el comportamiento esperado del mock
	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedAbilities, expectedTotal, nil)

	// Ejecutar la función que estamos probando
	abilities, total, err := useCase.List(ctx, filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAbilities, abilities)
	assert.Equal(t, expectedTotal, total)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}
