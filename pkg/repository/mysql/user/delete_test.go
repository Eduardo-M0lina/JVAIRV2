package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDelete_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta UPDATE (soft delete)
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		sqlmock.AnyArg(), 123,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), userID)

	// Verificar que no haya errores
	assert.NoError(t, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_InvalidID(t *testing.T) {
	// Configurar el mock de la base de datos
	db, _, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba con un ID inválido
	userID := "invalid_id"

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), userID)

	// Verificar que haya un error de ID inválido
	assert.Error(t, err)
	assert.Equal(t, "ID de usuario inválido", err.Error())
}

func TestDelete_UserNotFound(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "999"

	// Configurar la expectativa para la consulta UPDATE que no afecta a ninguna fila
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		sqlmock.AnyArg(), 999,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), userID)

	// Verificar que haya un error de usuario no encontrado
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_DatabaseError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Configurar la expectativa para la consulta UPDATE que falla
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		sqlmock.AnyArg(), 123,
	).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), userID)

	// Verificar que haya un error de base de datos
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_RowsAffectedError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	userID := "123"

	// Crear un resultado que devuelve un error al llamar a RowsAffected
	result := sqlmock.NewErrorResult(sql.ErrTxDone)

	// Configurar la expectativa para la consulta UPDATE
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`)).WithArgs(
		sqlmock.AnyArg(), 123,
	).WillReturnResult(result)

	// Ejecutar la función que estamos probando
	err := repo.Delete(context.Background(), userID)

	// Verificar que haya un error al obtener las filas afectadas
	assert.Error(t, err)
	assert.Equal(t, sql.ErrTxDone, err)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
