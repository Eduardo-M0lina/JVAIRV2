package role

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
	roleID := int64(1)
	title := "Administrator"
	scope := 1
	now := time.Now()

	expectedRole := &Role{
		ID:        roleID,
		Name:      "admin",
		Title:     &title,
		Scope:     &scope,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, roleID).Return(expectedRole, nil)

	// Ejecutar la función que estamos probando
	role, err := useCase.GetByID(ctx, roleID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)

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
	roleID := int64(999)
	expectedError := errors.New("rol no encontrado")

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, roleID).Return(nil, expectedError)

	// Ejecutar la función que estamos probando
	role, err := useCase.GetByID(ctx, roleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, role)

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
	roleName := "admin"
	title := "Administrator"
	scope := 1
	now := time.Now()

	expectedRole := &Role{
		ID:        1,
		Name:      roleName,
		Title:     &title,
		Scope:     &scope,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByName", ctx, roleName).Return(expectedRole, nil)

	// Ejecutar la función que estamos probando
	role, err := useCase.GetByName(ctx, roleName)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)

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
	title := "Administrator"
	scope := 1

	role := &Role{
		Name:  "admin",
		Title: &title,
		Scope: &scope,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Create", ctx, role).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, role)

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
	title := "Administrator"
	scope := 1

	role := &Role{
		ID:    1,
		Name:  "admin",
		Title: &title,
		Scope: &scope,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Update", ctx, role).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Update(ctx, role)

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
	roleID := int64(1)

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Delete", ctx, roleID).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Delete(ctx, roleID)

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
	filters := map[string]interface{}{"name": "admin"}
	page := 1
	pageSize := 10

	title1 := "Administrator"
	title2 := "User"
	scope1 := 1
	scope2 := 2
	now := time.Now()

	expectedRoles := []*Role{
		{
			ID:        1,
			Name:      "admin",
			Title:     &title1,
			Scope:     &scope1,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:        2,
			Name:      "user",
			Title:     &title2,
			Scope:     &scope2,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
	expectedTotal := 2

	// Configurar el comportamiento esperado del mock
	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedRoles, expectedTotal, nil)

	// Ejecutar la función que estamos probando
	roles, total, err := useCase.List(ctx, filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)
	assert.Equal(t, expectedTotal, total)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}
