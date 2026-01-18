package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainUser "github.com/your-org/jvairv2/pkg/domain/user"
)

func TestUpdate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "2"
	testUser := &domainUser.User{
		ID:     123,
		Name:   "Updated User",
		Email:  "updated@example.com",
		RoleID: &roleID,
	}

	// Configurar la expectativa para la consulta GetByID
	originalRoleID := "1"
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Original User", "original@example.com", "hashed_password", originalRoleID,
		time.Now(), "token123", time.Now(), time.Now(), nil,
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
		WHERE u.id = ? AND u.deleted_at IS NULL
	`)).WithArgs(123).WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET name = ?, email = ?, role_id = ?,
		    email_verified_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		testUser.Name, testUser.Email, roleID,
		nil, sqlmock.AnyArg(), testUser.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), testUser)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, testUser.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_UserNotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &domainUser.User{
		ID:     999,
		Name:   "Non-existent User",
		Email:  "nonexistent@example.com",
		RoleID: &roleID,
	}

	// Configurar la expectativa para la consulta GetByID que no encuentra al usuario
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
	`)).WithArgs(999).WillReturnError(ErrUserNotFound)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), testUser)

	// Verificar que haya un error de usuario no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DatabaseError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &domainUser.User{
		ID:     123,
		Name:   "Error User",
		Email:  "error@example.com",
		RoleID: &roleID,
	}

	// Configurar la expectativa para la consulta GetByID
	originalRoleID := "1"
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "Original User", "original@example.com", "hashed_password", originalRoleID,
		time.Now(), "token123", time.Now(), time.Now(), nil,
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
		WHERE u.id = ? AND u.deleted_at IS NULL
	`)).WithArgs(123).WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE que falla
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET name = ?, email = ?, role_id = ?,
		    email_verified_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		testUser.Name, testUser.Email, roleID,
		nil, sqlmock.AnyArg(), testUser.ID,
	).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), testUser)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
