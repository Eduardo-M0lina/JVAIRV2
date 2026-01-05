package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserRoles_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"
	now := time.Now()

	// Valores para campos opcionales
	title1 := "Administrator"
	title2 := "Regular User"
	scope1 := 1
	scope2 := 2

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "Admin", title1, scope1, now, now,
	).AddRow(
		2, "User", title2, scope2, now, now,
	)

	mock.ExpectQuery(`SELECT r.id, r.name, r.title, r.scope, r.created_at, r.updated_at FROM roles r INNER JOIN assigned_roles ar ON r.id = ar.role_id WHERE ar.entity_id = \? AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	roles, err := repo.GetUserRoles(context.Background(), userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, int64(1), roles[0].ID)
	assert.Equal(t, "Admin", roles[0].Name)
	assert.Equal(t, title1, *roles[0].Title)
	assert.Equal(t, &scope1, roles[0].Scope)
	assert.NotNil(t, roles[0].CreatedAt)
	assert.NotNil(t, roles[0].UpdatedAt)
	assert.Equal(t, int64(2), roles[1].ID)
	assert.Equal(t, "User", roles[1].Name)
	assert.Equal(t, title2, *roles[1].Title)
	assert.Equal(t, &scope2, roles[1].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRoles_InvalidID(t *testing.T) {
	// Configurar el mock de la base de datos
	db, _, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba con un ID inválido
	userID := "invalid_id"

	// Ejecutar la función que estamos probando
	roles, err := repo.GetUserRoles(context.Background(), userID)

	// Verificar que haya un error de ID inválido
	assert.Error(t, err)
	assert.Nil(t, roles)
	assert.Equal(t, "ID de usuario inválido", err.Error())
}

func TestGetUserRoles_QueryError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta que falla
	mock.ExpectQuery(`SELECT r.id, r.name, r.title, r.scope, r.created_at, r.updated_at FROM roles r INNER JOIN assigned_roles ar ON r.id = ar.role_id WHERE ar.entity_id = \? AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	roles, err := repo.GetUserRoles(context.Background(), userID)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Nil(t, roles)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRoles_ScanError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta con un error de escaneo
	// (devolvemos menos columnas de las esperadas)
	rows := sqlmock.NewRows([]string{
		"id", "name", // Faltan columnas
	}).AddRow(
		1, "Admin",
	)

	mock.ExpectQuery(`SELECT r.id, r.name, r.title, r.scope, r.created_at, r.updated_at FROM roles r INNER JOIN assigned_roles ar ON r.id = ar.role_id WHERE ar.entity_id = \? AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	roles, err := repo.GetUserRoles(context.Background(), userID)

	// Verificar que haya un error de escaneo
	assert.Error(t, err)
	assert.Nil(t, roles)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRoles_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta que devuelve un conjunto vacío
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	})

	mock.ExpectQuery(`SELECT r.id, r.name, r.title, r.scope, r.created_at, r.updated_at FROM roles r INNER JOIN assigned_roles ar ON r.id = ar.role_id WHERE ar.entity_id = \? AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	roles, err := repo.GetUserRoles(context.Background(), userID)

	// Verificar que no haya errores y que el resultado sea un slice vacío
	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 0)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
