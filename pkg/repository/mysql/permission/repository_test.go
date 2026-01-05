package permission

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/permission"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}

	repo := &Repository{db: db}
	return db, mock, repo
}

func TestGetByID_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = ?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	permission, err := repo.GetByID(context.Background(), permissionID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, permissionID, permission.ID)
	assert.Equal(t, abilityID, permission.AbilityID)
	assert.Equal(t, entityID, permission.EntityID)
	assert.Equal(t, entityType, permission.EntityType)
	assert.Equal(t, false, permission.Forbidden)
	assert.Equal(t, conditions, *permission.Conditions)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(999)

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = ?").
		WithArgs(permissionID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	permission, err := repo.GetByID(context.Background(), permissionID)

	// Verificar que haya un error de permiso no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrPermissionNotFound, err)
	assert.Nil(t, permission)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEntity_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	entityType := "App\\Models\\User"
	entityID := int64(10)
	now := time.Now()
	conditions1 := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		1, 2, entityID, entityType, false, conditions1, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE entity_type = \\? AND entity_id = \\?").
		WithArgs(entityType, entityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	permissions, err := repo.GetByEntity(context.Background(), entityType, entityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, int64(1), permissions[0].ID)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByAbility_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		1, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE ability_id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	permissions, err := repo.GetByAbility(context.Background(), abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, int64(1), permissions[0].ID)
	assert.Equal(t, abilityID, permissions[0].AbilityID)
	assert.Equal(t, entityID, permissions[0].EntityID)
	assert.Equal(t, entityType, permissions[0].EntityType)
	assert.Equal(t, false, permissions[0].Forbidden)
	assert.Equal(t, conditions, *permissions[0].Conditions)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExists_True(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	exists, err := repo.Exists(context.Background(), abilityID, entityID, entityType)

	// Verificar que no haya errores y que el resultado sea true
	assert.NoError(t, err)
	assert.True(t, exists)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExists_False(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	exists, err := repo.Exists(context.Background(), abilityID, entityID, entityType)

	// Verificar que no haya errores y que el resultado sea false
	assert.NoError(t, err)
	assert.False(t, exists)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM permissions").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		1, 2, 10, "App\\Models\\User", false, "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}", time.Now(), time.Now(),
	).AddRow(
		2, 3, 5, "App\\Models\\Role", true, nil, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions ORDER BY id LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	permissions, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, permissions, 2)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	filters := map[string]interface{}{
		"ability_id":  int64(2),
		"entity_type": "App\\Models\\User",
		"entity_id":   int64(10),
		"forbidden":   false,
	}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM permissions WHERE").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		1, 2, 10, "App\\Models\\User", false, "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}", time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	permissions, total, err := repo.List(context.Background(), filters, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, permissions, 1)

	// Verificar valores específicos del resultado
	assert.Equal(t, int64(1), permissions[0].ID)
	assert.Equal(t, int64(2), permissions[0].AbilityID)
	assert.Equal(t, int64(10), permissions[0].EntityID)
	assert.Equal(t, "App\\Models\\User", permissions[0].EntityType)
	assert.Equal(t, false, permissions[0].Forbidden)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM permissions").
		WillReturnRows(countRows)

	// Ejecutar la función que estamos probando
	permissions, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, permissions, 0)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE
	mock.ExpectExec("DELETE FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), permissionID)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(999)

	// Configurar la expectativa para la consulta GetByID
	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), permissionID)

	// Verificar que haya un error de permiso no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrPermissionNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_ErrorOnDelete(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE con error
	mock.ExpectExec("DELETE FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), permissionID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec("UPDATE permissions SET").
		WithArgs(abilityID, entityID, entityType, false, conditions, sqlmock.AnyArg(), permissionID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), permission)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, permission.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(999)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta GetByID
	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), permission)

	// Verificar que haya un error de permiso no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrPermissionNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DuplicatePermission(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		ID:         permissionID,
		AbilityID:  3, // Cambiamos el ability_id para que se verifique la existencia
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(permission.AbilityID, permission.EntityID, permission.EntityType).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), permission)

	// Verificar que haya un error de permiso duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicatePermission, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository(t *testing.T) {
	// Crear una base de datos mock
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Ejecutar la función que estamos probando
	repo := NewRepository(db)

	// Verificar que el repositorio se haya creado correctamente
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestUpdate_WithNilConditions(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: nil, // Condiciones nulas
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec("UPDATE permissions SET").
		WithArgs(abilityID, entityID, entityType, false, nil, sqlmock.AnyArg(), permissionID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), permission)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, permission.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	permissionID := int64(1)
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	now := time.Now()
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		ID:         permissionID,
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "ability_id", "entity_id", "entity_type", "forbidden", "conditions", "created_at", "updated_at",
	}).AddRow(
		permissionID, abilityID, entityID, entityType, false, conditions, now, now,
	)

	mock.ExpectQuery("SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at FROM permissions WHERE id = \\?").
		WithArgs(permissionID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE con error
	mock.ExpectExec("UPDATE permissions SET").
		WithArgs(abilityID, entityID, entityType, false, conditions, sqlmock.AnyArg(), permissionID).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), permission)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO permissions").
		WithArgs(abilityID, entityID, entityType, false, conditions, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), permission.ID)
	assert.NotNil(t, permission.CreatedAt)
	assert.NotNil(t, permission.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DuplicatePermission(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que haya un error de permiso duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicatePermission, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_WithNilConditions(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: nil, // Condiciones nulas
	}

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO permissions").
		WithArgs(abilityID, entityID, entityType, false, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), permission.ID)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrorOnExists(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta Exists con error
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrorOnInsert(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT con error
	mock.ExpectExec("INSERT INTO permissions").
		WithArgs(abilityID, entityID, entityType, false, conditions, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrorOnLastInsertID(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	conditions := "{\"field\":\"user_id\",\"operator\":\"=\",\"value\":10}"
	permission := &permission.Permission{
		AbilityID:  abilityID,
		EntityID:   entityID,
		EntityType: entityType,
		Forbidden:  false,
		Conditions: &conditions,
	}

	// Configurar la expectativa para la consulta Exists
	mock.ExpectQuery("SELECT 1 FROM permissions WHERE ability_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(abilityID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT con error en LastInsertId
	mock.ExpectExec("INSERT INTO permissions").
		WithArgs(abilityID, entityID, entityType, false, conditions, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), permission)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
