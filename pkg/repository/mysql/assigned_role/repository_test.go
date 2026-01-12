package assigned_role

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainAssignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

func TestAssign_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1
	assignedRole := &domainAssignedRole.AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta de verificación de existencia
	mock.ExpectQuery("SELECT id FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO assigned_roles").
		WithArgs(roleID, entityID, entityType, nil, nil, scope).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Assign(context.Background(), assignedRole)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), assignedRole.ID)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAssign_DuplicateAssignment(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1
	assignedRole := &domainAssignedRole.AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta de verificación de existencia
	rows := sqlmock.NewRows([]string{"id"}).AddRow(5)
	mock.ExpectQuery("SELECT id FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	err := repo.Assign(context.Background(), assignedRole)

	// Verificar que haya un error de asignación duplicada
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateAssignment, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAssign_AdditionalCases agrega casos de prueba adicionales para aumentar la cobertura
func TestAssign_AdditionalCases(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1
	assignedRole := &domainAssignedRole.AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: true, // Probamos con restricted = true
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta de verificación de existencia
	mock.ExpectQuery("SELECT id FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO assigned_roles").
		WithArgs(roleID, entityID, entityType, nil, nil, scope).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Assign(context.Background(), assignedRole)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), assignedRole.ID)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAssign_WithNilScope prueba el caso donde scope es nil
func TestAssign_WithNilScope(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	assignedRole := &domainAssignedRole.AssignedRole{
		RoleID:     roleID,
		EntityID:   entityID,
		EntityType: entityType,
		Restricted: false,
		Scope:      nil, // Scope es nil
	}

	// Configurar la expectativa para la consulta de verificación de existencia
	mock.ExpectQuery("SELECT id FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO assigned_roles").
		WithArgs(roleID, entityID, entityType, nil, nil, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Assign(context.Background(), assignedRole)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), assignedRole.ID)

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
	scope1 := 1
	scope2 := 2

	// Configurar la expectativa para la consulta
	restrictedToID := int64(5)
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		1, 2, entityID, entityType, nil, nil, scope1,
	).AddRow(
		2, 3, entityID, entityType, restrictedToID, "App\\Models\\Customer", scope2,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE entity_type = \\? AND entity_id = \\?").
		WithArgs(entityType, entityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, err := repo.GetByEntity(context.Background(), entityType, entityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, assignedRoles, 2)
	assert.Equal(t, int64(1), assignedRoles[0].ID)
	assert.Equal(t, int64(2), assignedRoles[0].RoleID)
	assert.Equal(t, entityID, assignedRoles[0].EntityID)
	assert.Equal(t, entityType, assignedRoles[0].EntityType)
	assert.Equal(t, false, assignedRoles[0].Restricted)
	assert.Equal(t, scope1, *assignedRoles[0].Scope)
	assert.Equal(t, int64(2), assignedRoles[1].ID)
	assert.Equal(t, int64(3), assignedRoles[1].RoleID)
	assert.Equal(t, true, assignedRoles[1].Restricted)
	assert.Equal(t, scope2, *assignedRoles[1].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEntity_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	entityType := "App\\Models\\User"
	entityID := int64(999)

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	})

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE entity_type = \\? AND entity_id = \\?").
		WithArgs(entityType, entityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, err := repo.GetByEntity(context.Background(), entityType, entityID)

	// Verificar que no haya errores y que el resultado sea un slice vacío
	assert.NoError(t, err)
	assert.Empty(t, assignedRoles)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByEntity_WithNilScope prueba el caso donde scope es nil
func TestGetByEntity_WithNilScope(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	entityType := "App\\Models\\User"
	entityID := int64(10)

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		1, 2, entityID, entityType, nil, nil, nil,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE entity_type = \\? AND entity_id = \\?").
		WithArgs(entityType, entityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, err := repo.GetByEntity(context.Background(), entityType, entityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Len(t, assignedRoles, 1)
	assert.Equal(t, int64(1), assignedRoles[0].ID)
	assert.Equal(t, int64(2), assignedRoles[0].RoleID)
	assert.Equal(t, entityID, assignedRoles[0].EntityID)
	assert.Equal(t, entityType, assignedRoles[0].EntityType)
	assert.Equal(t, false, assignedRoles[0].Restricted)
	assert.Nil(t, assignedRoles[0].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByEntity_ErrorCase prueba el caso de error en la consulta
func TestGetByEntity_ErrorCase(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	entityType := "App\\Models\\User"
	entityID := int64(10)

	// Configurar la expectativa para la consulta con error
	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE entity_type = \\? AND entity_id = \\?").
		WithArgs(entityType, entityID).
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	assignedRoles, err := repo.GetByEntity(context.Background(), entityType, entityID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Nil(t, assignedRoles)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	assignedRoleID := int64(1)
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"
	scope := 1

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		assignedRoleID, roleID, entityID, entityType, nil, nil, scope,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE id = \\?").
		WithArgs(assignedRoleID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRole, err := repo.GetByID(context.Background(), assignedRoleID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, assignedRole)
	assert.Equal(t, assignedRoleID, assignedRole.ID)
	assert.Equal(t, roleID, assignedRole.RoleID)
	assert.Equal(t, entityID, assignedRole.EntityID)
	assert.Equal(t, entityType, assignedRole.EntityType)
	assert.Equal(t, false, assignedRole.Restricted)
	assert.Equal(t, scope, *assignedRole.Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	assignedRoleID := int64(999)

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE id = \\?").
		WithArgs(assignedRoleID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	assignedRole, err := repo.GetByID(context.Background(), assignedRoleID)

	// Verificar que haya un error de asignación de rol no encontrada
	assert.Error(t, err)
	assert.Equal(t, ErrAssignedRoleNotFound, err)
	assert.Nil(t, assignedRole)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByID_WithNilScope prueba el caso donde scope es nil
func TestGetByID_WithNilScope(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	assignedRoleID := int64(1)
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		assignedRoleID, roleID, entityID, entityType, nil, nil, nil,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE id = \\?").
		WithArgs(assignedRoleID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRole, err := repo.GetByID(context.Background(), assignedRoleID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, assignedRole)
	assert.Equal(t, assignedRoleID, assignedRole.ID)
	assert.Equal(t, roleID, assignedRole.RoleID)
	assert.Equal(t, entityID, assignedRole.EntityID)
	assert.Equal(t, entityType, assignedRole.EntityType)
	assert.Equal(t, false, assignedRole.Restricted)
	assert.Nil(t, assignedRole.Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByID_ErrorCase prueba el caso de error en la consulta
func TestGetByID_ErrorCase(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	assignedRoleID := int64(1)

	// Configurar la expectativa para la consulta con error
	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE id = \\?").
		WithArgs(assignedRoleID).
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	assignedRole, err := repo.GetByID(context.Background(), assignedRoleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Nil(t, assignedRole)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHasRole_True(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(roleID, entityID, entityType).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	hasRole, err := repo.HasRole(context.Background(), roleID, entityID, entityType)

	// Verificar que no haya errores y que el resultado sea true
	assert.NoError(t, err)
	assert.True(t, hasRole)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHasRole_False(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT 1 FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	hasRole, err := repo.HasRole(context.Background(), roleID, entityID, entityType)

	// Verificar que no haya errores y que el resultado sea false
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHasRole_ErrorCase prueba el caso de error en la consulta
func TestHasRole_ErrorCase(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta con error
	mock.ExpectQuery("SELECT 1 FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\? LIMIT 1").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	hasRole, err := repo.HasRole(context.Background(), roleID, entityID, entityType)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.False(t, hasRole)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	scope1 := 1
	scope2 := 2

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	restrictedToID := int64(5)
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		1, 2, 10, "App\\Models\\User", nil, nil, scope1,
	).AddRow(
		2, 3, 20, "App\\Models\\User", restrictedToID, "App\\Models\\Customer", scope2,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles ORDER BY id ASC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, assignedRoles, 2)
	assert.Equal(t, int64(1), assignedRoles[0].ID)
	assert.Equal(t, int64(2), assignedRoles[0].RoleID)
	assert.Equal(t, int64(10), assignedRoles[0].EntityID)
	assert.Equal(t, "App\\Models\\User", assignedRoles[0].EntityType)
	assert.Equal(t, false, assignedRoles[0].Restricted)
	assert.Equal(t, scope1, *assignedRoles[0].Scope)
	assert.Equal(t, int64(2), assignedRoles[1].ID)
	assert.Equal(t, int64(3), assignedRoles[1].RoleID)
	assert.Equal(t, int64(20), assignedRoles[1].EntityID)
	assert.Equal(t, true, assignedRoles[1].Restricted)
	assert.Equal(t, scope2, *assignedRoles[1].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	scope := 1
	filters := map[string]interface{}{
		"role_id":     int64(2),
		"entity_type": "App\\Models\\User",
		"entity_id":   int64(10),
		"restricted":  false,
	}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles WHERE").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		1, 2, 10, "App\\Models\\User", nil, nil, scope,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := repo.List(context.Background(), filters, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, assignedRoles, 1)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles").
		WillReturnRows(countRows)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, assignedRoles, 0)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestList_CountError prueba el caso de error en la consulta COUNT
func TestList_CountError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT con error
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles").
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, assignedRoles)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestList_QueryError prueba el caso de error en la consulta SELECT
func TestList_QueryError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT con error
	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles ORDER BY id ASC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	assignedRoles, _, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Nil(t, assignedRoles)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestList_ScanError prueba el caso de error al escanear los resultados
func TestList_WithPartialFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	scope := 1
	filters := map[string]interface{}{
		"role_id": int64(2),
		// Solo incluimos un filtro para probar la construcción de la consulta con filtros parciales
	}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM assigned_roles WHERE role_id = \\?").
		WithArgs(int64(2)).
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "role_id", "entity_id", "entity_type", "restricted_to_id", "restricted_to_type", "scope",
	}).AddRow(
		1, 2, 10, "App\\Models\\User", nil, nil, scope,
	)

	mock.ExpectQuery("SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope FROM assigned_roles WHERE role_id = \\? ORDER BY id ASC LIMIT \\? OFFSET \\?").
		WithArgs(int64(2), 10, 0).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	assignedRoles, total, err := repo.List(context.Background(), filters, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, assignedRoles, 1)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRevoke_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta DELETE
	mock.ExpectExec("DELETE FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Revoke(context.Background(), roleID, entityID, entityType)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRevoke_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta DELETE
	mock.ExpectExec("DELETE FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Ejecutar la función que estamos probando
	err := repo.Revoke(context.Background(), roleID, entityID, entityType)

	// Verificar que haya un error de asignación de rol no encontrada
	assert.Error(t, err)
	assert.Equal(t, ErrAssignedRoleNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRevoke_ErrorCase prueba el caso de error en la consulta DELETE
func TestRevoke_ErrorCase(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta DELETE con error
	mock.ExpectExec("DELETE FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnError(sqlmock.ErrCancelled)

	// Ejecutar la función que estamos probando
	err := repo.Revoke(context.Background(), roleID, entityID, entityType)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.NotEqual(t, ErrAssignedRoleNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRevoke_ResultError prueba el caso de error al obtener el resultado de la consulta
func TestRevoke_ResultError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(2)
	entityID := int64(10)
	entityType := "App\\Models\\User"

	// Configurar la expectativa para la consulta DELETE con un resultado que dará error
	mock.ExpectExec("DELETE FROM assigned_roles WHERE role_id = \\? AND entity_id = \\? AND entity_type = \\?").
		WithArgs(roleID, entityID, entityType).
		WillReturnResult(sqlmock.NewErrorResult(sqlmock.ErrCancelled))

	// Ejecutar la función que estamos probando
	err := repo.Revoke(context.Background(), roleID, entityID, entityType)

	// Verificar que haya un error
	assert.Error(t, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
