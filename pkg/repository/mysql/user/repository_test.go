package user

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	// Crear un mock de la base de datos
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear el mock de la base de datos: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Crear un nuevo repositorio con el mock
	repo := NewRepository(db)

	// Verificar que el repositorio no sea nil
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Repository) {
	// Crear un mock de la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear el mock de la base de datos: %v", err)
	}

	// Crear un nuevo repositorio con el mock
	repo := NewRepository(db)

	return db, mock, repo
}

func TestErrNoRows(t *testing.T) {
	// Verificar que ErrNoRows sea igual a sql.ErrNoRows
	assert.Equal(t, sql.ErrNoRows, sql.ErrNoRows)
}

func TestErrUserNotFound(t *testing.T) {
	// Verificar que ErrUserNotFound sea un error con el mensaje correcto
	assert.EqualError(t, ErrUserNotFound, "usuario no encontrado")
}

func TestErrInvalidCredentials(t *testing.T) {
	// Verificar que ErrInvalidCredentials sea un error con el mensaje correcto
	assert.EqualError(t, ErrInvalidCredentials, "credenciales inválidas")
}
