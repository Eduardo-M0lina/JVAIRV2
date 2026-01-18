package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestList_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	page := 1
	pageSize := 10
	filters := map[string]interface{}{}
	now := time.Now()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
	`)).WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "User 1", "user1@example.com", "hashed_password", 1,
		now, "token123", now, now, nil,
		1, "Admin", "Administrator",
	).AddRow(
		456, "User 2", "user2@example.com", "hashed_password", 2,
		now, "token456", now, now, nil,
		2, "User", "Standard User",
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
		WHERE u.deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT ? OFFSET ?
	`)).WithArgs(pageSize, 0).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	users, total, err := repo.List(context.Background(), filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
	assert.Equal(t, int64(123), users[0].ID)
	assert.Equal(t, "User 1", users[0].Name)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.Equal(t, int64(456), users[1].ID)
	assert.Equal(t, "User 2", users[1].Name)
	assert.Equal(t, "user2@example.com", users[1].Email)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	page := 1
	pageSize := 10
	filters := map[string]interface{}{
		"name":    "User",
		"email":   "example",
		"role_id": "1",
	}
	now := time.Now()

	// Configurar la expectativa para la consulta COUNT con filtros
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
	`)).WithArgs(
		"%User%", "%example%", "1",
	).WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT con filtros
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "role_id",
		"email_verified_at", "remember_token", "created_at", "updated_at", "deleted_at",
		"role_id_int", "role_name", "role_title",
	}).AddRow(
		123, "User 1", "user1@example.com", "hashed_password", 1,
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
		WHERE u.deleted_at IS NULL
		 AND name LIKE ? AND email LIKE ? AND role_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`)).WithArgs(
		"%User%", "%example%", "1", pageSize, 0,
	).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	users, total, err := repo.List(context.Background(), filters, page, pageSize)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(123), users[0].ID)
	assert.Equal(t, "User 1", users[0].Name)
	assert.Equal(t, "user1@example.com", users[0].Email)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_CountError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	page := 1
	pageSize := 10
	filters := map[string]interface{}{}

	// Configurar la expectativa para la consulta COUNT que falla
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
	`)).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	users, total, err := repo.List(context.Background(), filters, page, pageSize)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, users)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_QueryError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	page := 1
	pageSize := 10
	filters := map[string]interface{}{}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
	`)).WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT que falla
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
		WHERE u.deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT ? OFFSET ?
	`)).WithArgs(pageSize, 0).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	users, total, err := repo.List(context.Background(), filters, page, pageSize)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, users)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_ScanError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	page := 1
	pageSize := 10
	filters := map[string]interface{}{}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM users u
		WHERE u.deleted_at IS NULL
	`)).WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT con un error de escaneo
	// (devolvemos menos columnas de las esperadas)
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", // Faltan columnas
	}).AddRow(
		123, "User 1", "user1@example.com",
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
		WHERE u.deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT ? OFFSET ?
	`)).WithArgs(pageSize, 0).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	users, total, err := repo.List(context.Background(), filters, page, pageSize)

	// Verificar que haya un error de escaneo
	assert.Error(t, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, users)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
