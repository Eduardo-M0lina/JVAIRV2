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

func TestGetByEmail_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "test@example.com"
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
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		123, "Test User", "test@example.com", "hashed_password", roleID,
		time.Now(), "token123", time.Now(), time.Now(), nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(email).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByEmail(context.Background(), email)

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

func TestGetByEmail_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "nonexistent@example.com"

	// Configurar la expectativa para la consulta que no encuentra resultados
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(email).WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByEmail(context.Background(), email)

	// Verificar que haya un error de usuario no encontrado
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail_DatabaseError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "test@example.com"

	// Configurar la expectativa para la consulta que devuelve un error de base de datos
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(email).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	result, err := repo.GetByEmail(context.Background(), email)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
