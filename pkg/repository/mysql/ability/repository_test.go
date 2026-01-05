package ability

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainAbility "github.com/your-org/jvairv2/pkg/domain/ability"
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
	abilityID := int64(1)
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	ability, err := repo.GetByID(context.Background(), abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, abilityID, ability.ID)
	assert.Equal(t, "create_user", ability.Name)
	assert.Equal(t, title, *ability.Title)
	assert.Equal(t, entityID, *ability.EntityID)
	assert.Equal(t, entityType, *ability.EntityType)
	assert.Equal(t, false, ability.OnlyOwned)
	assert.Equal(t, options, *ability.Options)
	assert.Equal(t, scope, *ability.Scope)
	assert.NotNil(t, ability.CreatedAt)
	assert.NotNil(t, ability.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(999)

	// Configurar la expectativa para la consulta
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	ability, err := repo.GetByID(context.Background(), abilityID)

	// Verificar que haya un error de habilidad no encontrada
	assert.Error(t, err)
	assert.Equal(t, ErrAbilityNotFound, err)
	assert.Nil(t, ability)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_WithNullFields(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	now := time.Now()

	// Configurar la expectativa para la consulta con campos nulos
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", nil, nil, nil, false, nil, nil, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	ability, err := repo.GetByID(context.Background(), abilityID)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, ability)
	assert.Equal(t, abilityID, ability.ID)
	assert.Equal(t, "create_user", ability.Name)
	assert.Nil(t, ability.Title)
	assert.Nil(t, ability.EntityID)
	assert.Nil(t, ability.EntityType)
	assert.Equal(t, false, ability.OnlyOwned)
	assert.Nil(t, ability.Options)
	assert.Nil(t, ability.Scope)
	assert.NotNil(t, ability.CreatedAt)
	assert.NotNil(t, ability.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnError(ErrAbilityNotFound)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO abilities").
		WithArgs(ability.Name, title, entityID, entityType, false, options, scope, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ability.ID)
	assert.NotNil(t, ability.CreatedAt)
	assert.NotNil(t, ability.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DuplicateName(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "create_user", title, entityID, entityType, false, options, scope, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

	// Verificar que haya un error de nombre duplicado
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateName, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_WithNilFields(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      nil,
		EntityID:   nil,
		EntityType: nil,
		OnlyOwned:  false,
		Options:    nil,
		Scope:      nil,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnError(ErrAbilityNotFound)

	// Configurar la expectativa para la consulta INSERT
	mock.ExpectExec("INSERT INTO abilities").
		WithArgs(ability.Name, nil, nil, nil, false, nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ability.ID)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_ErrorOnGetByName(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByName con error
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

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
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnError(ErrAbilityNotFound)

	// Configurar la expectativa para la consulta INSERT con error
	mock.ExpectExec("INSERT INTO abilities").
		WithArgs(ability.Name, title, entityID, entityType, false, options, scope, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

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
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByName
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE name = \\?").
		WithArgs(ability.Name).
		WillReturnError(ErrAbilityNotFound)

	// Configurar la expectativa para la consulta INSERT con error en LastInsertId
	mock.ExpectExec("INSERT INTO abilities").
		WithArgs(ability.Name, title, entityID, entityType, false, options, scope, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), ability)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE
	mock.ExpectExec("DELETE FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), abilityID)

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
	abilityID := int64(999)

	// Configurar la expectativa para la consulta GetByID
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), abilityID)

	// Verificar que haya un error de habilidad no encontrada
	assert.Error(t, err)
	assert.Equal(t, ErrAbilityNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_ErrorOnExec(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE con error
	mock.ExpectExec("DELETE FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), abilityID)

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
	abilityID := int64(1)
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE sin filas afectadas
	mock.ExpectExec("DELETE FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), abilityID)

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
	abilityID := int64(1)
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta DELETE con error en RowsAffected
	mock.ExpectExec("DELETE FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), abilityID)

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
	title1 := "Create User"
	title2 := "Update User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "create_user", title1, entityID, entityType, false, options, scope, now, now,
	).AddRow(
		2, "update_user", title2, entityID, entityType, true, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, abilities, 2)
	assert.Equal(t, int64(1), abilities[0].ID)
	assert.Equal(t, "create_user", abilities[0].Name)
	assert.Equal(t, title1, *abilities[0].Title)
	assert.Equal(t, int64(2), abilities[1].ID)
	assert.Equal(t, "update_user", abilities[1].Name)
	assert.Equal(t, title2, *abilities[1].Title)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	now := time.Now()
	title := "Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	filters := map[string]interface{}{
		"name":        "create",
		"entity_type": entityType,
	}

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities WHERE").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		1, "create_user", title, entityID, entityType, false, options, scope, now, now,
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), filters, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, abilities, 1)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_EmptyResult(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities").
		WillReturnRows(countRows)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, abilities)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_CountError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT con error
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities").
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, abilities)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_QueryError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT con error
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities").
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, abilities)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_ScanError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Configurar la expectativa para la consulta COUNT
	countRows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM abilities").
		WillReturnRows(countRows)

	// Configurar la expectativa para la consulta SELECT con columnas incorrectas
	rows := sqlmock.NewRows([]string{
		"id", // Faltan columnas
	}).AddRow(
		1, // Solo un valor
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities").
		WillReturnRows(rows)

	// Ejecutar la función que estamos probando
	abilities, total, err := repo.List(context.Background(), nil, 1, 10)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, 0, total)
	assert.Nil(t, abilities)

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

func TestUpdate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	title := "Updated Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		ID:         abilityID,
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", "Original Title", entityID, entityType, false, options, scope, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec("UPDATE abilities SET").
		WithArgs(ability.Name, title, entityID, entityType, false, options, scope, sqlmock.AnyArg(), abilityID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), ability)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, ability.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	title := "Updated Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		ID:         999,
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByID
	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(ability.ID).
		WillReturnError(sql.ErrNoRows)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), ability)

	// Verificar que haya un error de habilidad no encontrada
	assert.Error(t, err)
	assert.Equal(t, ErrAbilityNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_WithNilFields(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	ability := &domainAbility.Ability{
		ID:         abilityID,
		Name:       "create_user",
		Title:      nil,
		EntityID:   nil,
		EntityType: nil,
		OnlyOwned:  false,
		Options:    nil,
		Scope:      nil,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", "Original Title", 10, "App\\Models\\User", false, "{}", 1, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec("UPDATE abilities SET").
		WithArgs(ability.Name, nil, nil, nil, false, nil, nil, sqlmock.AnyArg(), abilityID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), ability)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.NotNil(t, ability.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_ErrorOnExec(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	abilityID := int64(1)
	title := "Updated Create User"
	entityType := "App\\Models\\User"
	options := "{}"
	scope := 1
	entityID := int64(10)
	ability := &domainAbility.Ability{
		ID:         abilityID,
		Name:       "create_user",
		Title:      &title,
		EntityID:   &entityID,
		EntityType: &entityType,
		OnlyOwned:  false,
		Options:    &options,
		Scope:      &scope,
	}

	// Configurar la expectativa para la consulta GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "title", "entity_id", "entity_type", "only_owned", "options", "scope", "created_at", "updated_at",
	}).AddRow(
		abilityID, "create_user", "Original Title", entityID, entityType, false, options, scope, time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at FROM abilities WHERE id = \\?").
		WithArgs(abilityID).
		WillReturnRows(rows)

	// Configurar la expectativa para la consulta UPDATE con error
	mock.ExpectExec("UPDATE abilities SET").
		WithArgs(ability.Name, title, entityID, entityType, false, options, scope, sqlmock.AnyArg(), abilityID).
		WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Update(context.Background(), ability)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
