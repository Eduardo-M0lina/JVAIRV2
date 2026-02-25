package invoice

import (
	"database/sql"
)

// Repository implementa el repositorio MySQL para invoices
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de invoices
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
