package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/jvairv2/pkg/domain/ability"
	"github.com/your-org/jvairv2/pkg/domain/role"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_GetByID(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	roleID := "admin"
	now := time.Now()

	expectedUser := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// Ejecutar la función que estamos probando
	user, err := useCase.GetByID(ctx, userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByEmail(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	email := "john@example.com"
	roleID := "admin"
	now := time.Now()

	expectedUser := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     email,
		Password:  "hashed_password",
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEmail", ctx, email).Return(expectedUser, nil)

	// Ejecutar la función que estamos probando
	user, err := useCase.GetByEmail(ctx, email)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := "admin"

	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRepo.On("GetByEmail", ctx, user.Email).Return(nil, ErrUserNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*User)
		u.ID = 1 // Simular la asignación de ID

		// Verificar que la contraseña fue hasheada
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password123"))
		assert.NoError(t, err)

		// Verificar que se establecieron los valores predeterminados
		assert.NotNil(t, u.CreatedAt)
		assert.NotNil(t, u.UpdatedAt)
		assert.True(t, u.IsActive)
	})

	mockRoleRepo.On("GetByName", mock.Anything, roleID).Return(&role.Role{ID: 1, Name: roleID}, nil)
	mockAssignedRoleRepo.On("Assign", mock.Anything, mock.AnythingOfType("*assigned_role.AssignedRole")).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, user)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockAssignedRoleRepo.AssertExpectations(t)
}

func TestUseCase_Create_DuplicateEmail(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	roleID := "admin"
	now := time.Now()

	existingUser := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	newUser := &User{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Mismo email que el usuario existente
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRepo.On("GetByEmail", ctx, newUser.Email).Return(existingUser, nil)

	// Ejecutar la función que estamos probando
	err := useCase.Create(ctx, newUser)

	// Verificar que haya un error de email duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateEmail, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_VerifyCredentials_Success(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	email := "john@example.com"
	password := "password123"
	roleID := "admin"
	now := time.Now()

	// Crear un hash de la contraseña para simular la contraseña almacenada
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     email,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

	// Ejecutar la función que estamos probando
	authenticatedUser, err := useCase.VerifyCredentials(ctx, email, password)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, user, authenticatedUser)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_VerifyCredentials_InvalidCredentials(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	email := "john@example.com"
	password := "password123"
	wrongPassword := "wrongpassword"
	roleID := "admin"
	now := time.Now()

	// Crear un hash de la contraseña para simular la contraseña almacenada
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     email,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

	// Ejecutar la función que estamos probando
	authenticatedUser, err := useCase.VerifyCredentials(ctx, email, wrongPassword)

	// Verificar que haya un error de credenciales inválidas
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, authenticatedUser)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_VerifyCredentials_InactiveUser(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	email := "john@example.com"
	password := "password123"
	roleID := "admin"
	now := time.Now()

	// Crear un hash de la contraseña para simular la contraseña almacenada
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     email,
		Password:  string(hashedPassword),
		IsActive:  false, // Usuario inactivo
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

	// Ejecutar la función que estamos probando
	authenticatedUser, err := useCase.VerifyCredentials(ctx, email, password)

	// Verificar que haya un error de usuario inactivo
	assert.Error(t, err)
	assert.Equal(t, ErrUserInactive, err)
	assert.Nil(t, authenticatedUser)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	roleID := "admin"
	now := time.Now()

	existingUser := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	updatedUser := &User{
		ID:       1,
		RoleID:   &roleID,
		Name:     "John Updated",
		Email:    "john@example.com",
		Password: "hashed_password", // Misma contraseña
		IsActive: true,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*User)

		// Verificar que se actualizó el timestamp
		assert.NotNil(t, u.UpdatedAt)
		assert.Equal(t, "John Updated", u.Name)
	})

	// Ejecutar la función que estamos probando
	err := useCase.Update(ctx, updatedUser)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update_WithNewPassword(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	roleID := "admin"
	now := time.Now()

	existingUser := &User{
		ID:        1,
		RoleID:    &roleID,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "old_hashed_password",
		IsActive:  true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	updatedUser := &User{
		ID:       1,
		RoleID:   &roleID,
		Name:     "John Updated",
		Email:    "john@example.com",
		Password: "new_password", // Nueva contraseña
		IsActive: true,
	}

	// Configurar el comportamiento esperado de los mocks
	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*User)

		// Verificar que la contraseña fue hasheada
		assert.NotEqual(t, "new_password", u.Password)
		assert.NotEqual(t, "old_hashed_password", u.Password)

		// Verificar que se actualizó el timestamp
		assert.NotNil(t, u.UpdatedAt)
	})

	// Ejecutar la función que estamos probando
	err := useCase.Update(ctx, updatedUser)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamaron a los métodos de los repositorios con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Delete(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"

	// Configurar el comportamiento esperado del mock
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Ejecutar la función que estamos probando
	err := useCase.Delete(ctx, userID)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_List(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	filters := map[string]interface{}{"name": "John"}
	page := 1
	pageSize := 10
	roleID1 := "admin"
	roleID2 := "user"
	now := time.Now()

	expectedUsers := []*User{
		{
			ID:        1,
			RoleID:    &roleID1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "hashed_password1",
			IsActive:  true,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:        2,
			RoleID:    &roleID2,
			Name:      "John Smith",
			Email:     "john.smith@example.com",
			Password:  "hashed_password2",
			IsActive:  true,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
	expectedTotal := 2

	// Configurar el comportamiento esperado del mock
	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedUsers, expectedTotal, nil)

	// Ejecutar la función que estamos probando
	users, total, err := useCase.List(ctx, filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	assert.Equal(t, expectedTotal, total)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetUserRoles(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	title1 := "Administrator"
	title2 := "User"
	scope1 := 1
	scope2 := 2
	now := time.Now()

	expectedRoles := []*role.Role{
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

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetUserRoles", ctx, userID).Return(expectedRoles, nil)

	// Ejecutar la función que estamos probando
	roles, err := useCase.GetUserRoles(ctx, userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetUserAbilities(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	title1 := "Create User"
	title2 := "Edit User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope1 := 1
	scope2 := 2
	entityID := int64(10)
	now := time.Now()

	expectedAbilities := []*ability.Ability{
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

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetUserAbilities", ctx, userID).Return(expectedAbilities, nil)

	// Ejecutar la función que estamos probando
	abilities, err := useCase.GetUserAbilities(ctx, userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, expectedAbilities, abilities)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_HasAbility_True(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	abilityName := "create_user"
	title1 := "Create User"
	title2 := "Edit User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope1 := 1
	scope2 := 2
	entityID := int64(10)
	now := time.Now()

	abilities := []*ability.Ability{
		{
			ID:         1,
			Name:       "create_user", // Coincide con el nombre buscado
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

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetUserAbilities", ctx, userID).Return(abilities, nil)

	// Ejecutar la función que estamos probando
	hasAbility, err := useCase.HasAbility(ctx, userID, abilityName)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.True(t, hasAbility)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}

func TestUseCase_HasAbility_False(t *testing.T) {
	// Crear los mocks de los repositorios
	mockRepo := new(MockRepository)
	mockAssignedRoleRepo := new(MockAssignedRoleRepository)
	mockRoleRepo := new(MockRoleRepository)

	// Crear el caso de uso con los mocks
	useCase := NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)

	// Datos de prueba
	ctx := context.Background()
	userID := "1"
	abilityName := "delete_user" // No existe en las habilidades del usuario
	title1 := "Create User"
	title2 := "Edit User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope1 := 1
	scope2 := 2
	entityID := int64(10)
	now := time.Now()

	abilities := []*ability.Ability{
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

	// Configurar el comportamiento esperado del mock
	mockRepo.On("GetUserAbilities", ctx, userID).Return(abilities, nil)

	// Ejecutar la función que estamos probando
	hasAbility, err := useCase.HasAbility(ctx, userID, abilityName)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.False(t, hasAbility)

	// Verificar que se llamó al método del repositorio con los argumentos correctos
	mockRepo.AssertExpectations(t)
}
