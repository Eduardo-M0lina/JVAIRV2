package job

import (
	"database/sql"
)

// Repository implementa el repositorio MySQL para jobs
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de jobs
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
