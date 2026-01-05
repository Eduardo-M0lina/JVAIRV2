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

func TestCreate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "New Test User",
		Email:    "newuser@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta GetByEmail (que debe devolver ErrUserNotFound)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(testUser.Email).WillReturnError(ErrUserNotFound)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)).WithArgs(
		testUser.Name, testUser.Email, sqlmock.AnyArg(), roleID,
		nil, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(123, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(123), testUser.ID)
	assert.NotNil(t, testUser.CreatedAt)
	assert.NotNil(t, testUser.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DuplicateEmail(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "Duplicate User",
		Email:    "existing@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta GetByEmail (que debe devolver un usuario existente)
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		456, "Existing User", "existing@example.com", "hashed_password", roleID,
		nil, nil, time.Now(), time.Now(), nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(testUser.Email).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que haya un error de email duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateEmail, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_GetByEmailError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "Error User",
		Email:    "error@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta GetByEmail (que debe devolver un error de base de datos)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(testUser.Email).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_InsertError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "Insert Error User",
		Email:    "inserterror@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta GetByEmail (que debe devolver ErrUserNotFound)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, email, password, role_id,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`)).WithArgs(testUser.Email).WillReturnError(ErrUserNotFound)

	// Configurar la expectativa para la consulta INSERT (que debe devolver un error)
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)).WithArgs(
		testUser.Name, testUser.Email, sqlmock.AnyArg(), roleID,
		nil, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
