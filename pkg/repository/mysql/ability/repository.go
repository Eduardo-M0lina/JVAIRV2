package ability

import (
	"database/sql"
	"errors"
)

// Errores comunes del repositorio
var (
	ErrAbilityNotFound = errors.New("ability no encontrada")
	ErrDuplicateName   = errors.New("nombre de ability ya está en uso")
)

// Repository implementa la interfaz ability.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de abilities
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
