package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyCredentials_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	now := time.Now()

	// Configurar la expectativa para la consulta GetByEmail
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Test User", email, string(hashedPassword), 1,
		now, "token123", now, now, nil,
		1, "Admin", "Administrator",
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
		WHERE u.email = ? AND u.deleted_at IS NULL
	`)).WithArgs(email).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	u, err := repo.VerifyCredentials(context.Background(), email, password)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, int64(123), u.ID)
	assert.Equal(t, "Test User", u.Name)
	assert.Equal(t, email, u.Email)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyCredentials_UserNotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "nonexistent@example.com"
	password := "password123"

	// Configurar la expectativa para la consulta GetByEmail que no encuentra al usuario
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
		WHERE u.email = ? AND u.deleted_at IS NULL
	`)).WithArgs(email).WillReturnError(ErrUserNotFound)

	// Ejecutar la función que estamos probando
	u, err := repo.VerifyCredentials(context.Background(), email, password)

	// Verificar que haya un error de credenciales inválidas
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrInvalidCredentials, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyCredentials_InactiveUser(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "inactive@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	now := time.Now()
	deletedAt := now // Usuario inactivo (borrado)

	// Configurar la expectativa para la consulta GetByEmail
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Inactive User", email, string(hashedPassword), 1,
		now, "token123", now, now, deletedAt,
		1, "Admin", "Administrator",
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
		WHERE u.email = ? AND u.deleted_at IS NULL
	`)).WithArgs(email).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	u, err := repo.VerifyCredentials(context.Background(), email, password)

	// Verificar que haya un error de credenciales inválidas
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrInvalidCredentials, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyCredentials_WrongPassword(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	now := time.Now()

	// Configurar la expectativa para la consulta GetByEmail
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Test User", email, string(hashedPassword), 1,
		now, "token123", now, now, nil,
		1, "Admin", "Administrator",
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
		WHERE u.email = ? AND u.deleted_at IS NULL
	`)).WithArgs(email).WillReturnRows(rows)

	// Ejecutar la función que estamos probando con una contraseña incorrecta
	u, err := repo.VerifyCredentials(context.Background(), email, wrongPassword)

	// Verificar que haya un error de credenciales inválidas
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrInvalidCredentials, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerifyCredentials_DatabaseError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	email := "error@example.com"
	password := "password123"

	// Configurar la expectativa para la consulta GetByEmail que devuelve un error de base de datos
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
		WHERE u.email = ? AND u.deleted_at IS NULL
	`)).WithArgs(email).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	u, err := repo.VerifyCredentials(context.Background(), email, password)

	// Verificar que haya un error de credenciales inválidas
	// Nota: La función VerifyCredentials convierte cualquier error en ErrInvalidCredentials
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrInvalidCredentials, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
