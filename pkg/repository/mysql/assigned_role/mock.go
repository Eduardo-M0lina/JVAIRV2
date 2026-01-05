package assigned_role

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}

	repo := &Repository{db: db}
	return db, mock, repo
}
