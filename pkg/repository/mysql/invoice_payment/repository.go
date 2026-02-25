package invoice_payment

import (
	"database/sql"
)

// Repository implementa el repositorio MySQL para invoice payments
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de invoice payments
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
