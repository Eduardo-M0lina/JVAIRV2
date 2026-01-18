package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

func TestGetByID_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"
	roleID := "1"
	expectedUser := &user.User{
		ID:       123,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
		RoleID:   &roleID,
		IsActive: true,
	}

	// Configurar la expectativa para la consulta
	roleName := "Admin"
	roleTitle := "Administrator"
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Test User", "test@example.com", "hashed_password", roleID,
		time.Now(), "token123", time.Now(), time.Now(), nil,
		1, roleName, roleTitle,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT u.id, u.name, u.email, u.password, u.role_id,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at,
		       r.id as role_id_int, r.name as role_name, r.title as role_title
		FROM users u
		LEFT JOIN (
			SELECT ar.entity_id, ar.role_id
			FROM assigned_roles ar
			WHERE ar.entity_type = 'App\\Models\\User'
			GROUP BY ar.entity_id
		) ar ON ar.entity_id = u.id
		LEFT JOIN roles r ON r.id = ar.role_id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`)).WithArgs(123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByID(context.Background(), userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Password, result.Password)
	assert.Equal(t, expectedUser.RoleID, result.RoleID)
	assert.Equal(t, expectedUser.IsActive, result.IsActive)
	assert.NotNil(t, result.EmailVerifiedAt)
	assert.NotNil(t, result.RememberToken)
	assert.NotNil(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.Nil(t, result.DeletedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_InvalidID(t *testing.T) {
	// Configurar el mock de la base de datos
	db, _, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba con un ID inválido
	userID := "invalid_id"

	// Ejecutar la función que estamos probando
	result, err := repo.GetByID(context.Background(), userID)

	// Verificar que haya un error de ID inválido
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ID de usuario inválido", err.Error())
}

func TestGetByID_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "999"

	// Configurar la expectativa para la consulta que no encuentra resultados
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT u.id, u.name, u.email, u.password, u.role_id,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at,
		       r.id as role_id_int, r.name as role_name, r.title as role_title
		FROM users u
		LEFT JOIN (
			SELECT ar.entity_id, ar.role_id
			FROM assigned_roles ar
			WHERE ar.entity_type = 'App\\Models\\User'
			GROUP BY ar.entity_id
		) ar ON ar.entity_id = u.id
		LEFT JOIN roles r ON r.id = ar.role_id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`)).WithArgs(999).WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByID(context.Background(), userID)

	// Verificar que haya un error de usuario no encontrado
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_DatabaseError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta que devuelve un error de base de datos
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT u.id, u.name, u.email, u.password, u.role_id,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at,
		       r.id as role_id_int, r.name as role_name, r.title as role_title
		FROM users u
		LEFT JOIN (
			SELECT ar.entity_id, ar.role_id
			FROM assigned_roles ar
			WHERE ar.entity_type = 'App\\Models\\User'
			GROUP BY ar.entity_id
		) ar ON ar.entity_id = u.id
		LEFT JOIN roles r ON r.id = ar.role_id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`)).WithArgs(123).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByID(context.Background(), userID)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
