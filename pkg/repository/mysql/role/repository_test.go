package role

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainRole "github.com/your-org/jvairv2/pkg/domain/role"
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
	roleID := int64(1)
	now := time.Now()
	title := "Administrator"
	scope := 1

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "admin", title, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	role, err := repo.GetByID(context.Background(), roleID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, "admin", role.Name)
	assert.Equal(t, title, *role.Title)
	assert.Equal(t, scope, *role.Scope)
	assert.NotNil(t, role.CreatedAt)
	assert.NotNil(t, role.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(999)

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	role, err := repo.GetByID(context.Background(), roleID)

	// Verificar que haya un error de rol no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrRoleNotFound, err)
	assert.Nil(t, role)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByName_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleName := "admin"
	roleID := int64(1)
	now := time.Now()
	title := "Administrator"
	scope := 1

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, roleName, title, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = ?").
		WithArgs(roleName).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	role, err := repo.GetByName(context.Background(), roleName)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, roleName, role.Name)
	assert.Equal(t, title, *role.Title)
	assert.Equal(t, scope, *role.Scope)
	assert.NotNil(t, role.CreatedAt)
	assert.NotNil(t, role.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleName := "new_role"
	title := "New Role"
	scope := 2
	role := &domainRole.Role{
		Name:  roleName,
		Title: &title,
		Scope: &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = ?").
		WithArgs(roleName).
		WillReturnError(ErrRoleNotFound)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO roles").
		WithArgs(roleName, title, scope, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), role.ID)
	assert.NotNil(t, role.CreatedAt)
	assert.NotNil(t, role.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DuplicateName(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleName := "existing_role"
	title := "Existing Role"
	scope := 2
	role := &domainRole.Role{
		Name:  roleName,
		Title: &title,
		Scope: &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		1, roleName, title, scope, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = ?").
		WithArgs(roleName).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

	// Verificar que haya un error de nombre duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateName, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_AlreadyExists(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Administrator"
	role := &domainRole.Role{
		Name:  "admin",
		Title: &title,
	}

	// Configurar la expectativa para la consulta GetByName
	rows := sqlmock.NewRows([]string{"id", "name", "title", "scope", "created_at", "updated_at"}).
		AddRow(1, "admin", "Administrator", 1, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = \\?").
		WithArgs(role.Name).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

	// Verificar que haya un error de nombre duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateName, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrorOnGetByName(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Administrator"
	role := &domainRole.Role{
		Name:  "admin",
		Title: &title,
	}

	// Configurar la expectativa para la consulta GetByName con error
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = \\?").
		WithArgs(role.Name).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

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
	title := "Administrator"
	role := &domainRole.Role{
		Name:  "admin",
		Title: &title,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = \\?").
		WithArgs(role.Name).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT con error
	mock.ExpectExec("INSERT INTO roles").
		WithArgs(role.Name, role.Title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

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
	title := "Administrator"
	role := &domainRole.Role{
		Name:  "admin",
		Title: &title,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = \\?").
		WithArgs(role.Name).
		WillReturnError(sql.ErrNoRows)

	// Configurar la expectativa para la consulta INSERT con error en LastInsertId
	mock.ExpectExec("INSERT INTO roles").
		WithArgs(role.Name, role.Title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), role)

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
	roleID := int64(1)
	roleName := "updated_role"
	title := "Updated Role"
	scope := 3
	role := &domainRole.Role{
		ID:    roleID,
		Name:  roleName,
		Title: &title,
		Scope: &scope,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "old_name", "Old Title", 2, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name = ?").
		WithArgs(roleName).
		WillReturnError(ErrRoleNotFound)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec("UPDATE roles SET name = \\?, title = \\?, scope = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(roleName, title, scope, sqlmock.AnyArg(), roleID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), role)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, role.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(1)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "role_to_delete", "Role to Delete", 1, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE
	mock.ExpectExec("DELETE FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), roleID)

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
	roleID := int64(999)

	// Configurar la expectativa para la consulta GetByID
	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), roleID)

	// Verificar que haya un error de rol no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrRoleNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_ErrorOnExec(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(1)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "role_to_delete", "Role to Delete", 1, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE con error
	mock.ExpectExec("DELETE FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), roleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_NoRowsAffected(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(1)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "role_to_delete", "Role to Delete", 1, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE sin filas afectadas
	mock.ExpectExec("DELETE FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), roleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_ErrorOnRowsAffected(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := int64(1)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		roleID, "role_to_delete", "Role to Delete", 1, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE con error en RowsAffected
	mock.ExpectExec("DELETE FROM roles WHERE id = ?").
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), roleID)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	now := time.Now()
	title1 := "Administrator"
	scope1 := 1
	title2 := "User"
	scope2 := 2

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM roles").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "admin", title1, scope1, now, now,
	).AddRow(
		2, "user", title2, scope2, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles ORDER BY name LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	roles, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, roles, 2)
	assert.Equal(t, int64(1), roles[0].ID)
	assert.Equal(t, "admin", roles[0].Name)
	assert.Equal(t, title1, *roles[0].Title)
	assert.Equal(t, scope1, *roles[0].Scope)
	assert.Equal(t, int64(2), roles[1].ID)
	assert.Equal(t, "user", roles[1].Name)
	assert.Equal(t, title2, *roles[1].Title)
	assert.Equal(t, scope2, *roles[1].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	now := time.Now()
	title := "Administrator"
	scope := 1
	filters := map[string]interface{}{
		"name":  "admin",
		"scope": 1,
	}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM roles WHERE name LIKE \\? AND scope = \\?").
		WithArgs("%admin%", 1).
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "admin", title, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, scope, created_at, updated_at FROM roles WHERE name LIKE \\? AND scope = \\? ORDER BY name LIMIT 10 OFFSET 0").
		WithArgs("%admin%", 1).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	roles, total, err := repo.List(context.Background(), filters, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, roles, 1)
	assert.Equal(t, int64(1), roles[0].ID)
	assert.Equal(t, "admin", roles[0].Name)
	assert.Equal(t, title, *roles[0].Title)
	assert.Equal(t, scope, *roles[0].Scope)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM roles").
		WillReturnRows(countRows)

	// Ejecutar la función que estamos probando
	roles, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, roles)

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
