package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserAbilities_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"
	now := time.Now()

	// Variables para los campos opcionales
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "create_user", title, nil, entityType, false, options, scope, now, now,
	).AddRow(
		2, "edit_user", "Edit User", nil, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery(`SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id WHERE p.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\User' UNION SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id INNER JOIN assigned_roles ar ON p.entity_id = ar.role_id WHERE ar.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\Role' AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123, 123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, err := repo.GetUserAbilities(context.Background(), userID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, abilities, 2)
	assert.Equal(t, int64(1), abilities[0].ID)
	assert.Equal(t, "create_user", abilities[0].Name)
	assert.Equal(t, title, *abilities[0].Title)
	assert.Equal(t, entityType, *abilities[0].EntityType)
	assert.Equal(t, options, *abilities[0].Options)
	assert.Equal(t, scope, *abilities[0].Scope)
	assert.NotNil(t, abilities[0].CreatedAt)
	assert.NotNil(t, abilities[0].UpdatedAt)
	assert.Equal(t, int64(2), abilities[1].ID)
	assert.Equal(t, "edit_user", abilities[1].Name)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAbilities_InvalidID(t *testing.T) {
	// Configurar el mock de la base de datos
	db, _, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba con un ID inválido
	userID := "invalid_id"

	// Ejecutar la función que estamos probando
	abilities, err := repo.GetUserAbilities(context.Background(), userID)

	// Verificar que haya un error de ID inválido
	assert.Error(t, err)
	assert.Nil(t, abilities)
	assert.Equal(t, "ID de usuario inválido", err.Error())
}

func TestGetUserAbilities_QueryError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta que falla
	mock.ExpectQuery(`SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id WHERE p.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\User' UNION SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id INNER JOIN assigned_roles ar ON p.entity_id = ar.role_id WHERE ar.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\Role' AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123, 123).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	abilities, err := repo.GetUserAbilities(context.Background(), userID)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Nil(t, abilities)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAbilities_ScanError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta con un error de escaneo
	// (devolvemos menos columnas de las esperadas)
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", // Faltan columnas
	}).AddRow(
		1, "create_user", "Create User",
	)

	mock.ExpectQuery(`SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id WHERE p.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\User' UNION SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id INNER JOIN assigned_roles ar ON p.entity_id = ar.role_id WHERE ar.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\Role' AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123, 123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, err := repo.GetUserAbilities(context.Background(), userID)

	// Verificar que haya un error de escaneo
	assert.Error(t, err)
	assert.Nil(t, abilities)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAbilities_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta que devuelve un conjunto vacío
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	})

	mock.ExpectQuery(`SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id WHERE p.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\User' UNION SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at FROM abilities a INNER JOIN permissions p ON a.id = p.ability_id INNER JOIN assigned_roles ar ON p.entity_id = ar.role_id WHERE ar.entity_id = \? AND p.entity_type = 'App\\\\Models\\\\Role' AND ar.entity_type = 'App\\\\Models\\\\User'`).WithArgs(123, 123).WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, err := repo.GetUserAbilities(context.Background(), userID)

	// Verificar que no haya errores y que el resultado sea un slice vacío
	assert.NoError(t, err)
	assert.NotNil(t, abilities)
	assert.Len(t, abilities, 0)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
