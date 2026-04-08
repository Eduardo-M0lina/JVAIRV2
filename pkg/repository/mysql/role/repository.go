package role

import (
	"database/sql"
	"errors"
)

// Errores comunes del repositorio
var (
	ErrRoleNotFound  = errors.New("rol no encontrado")
	ErrDuplicateName = errors.New("nombre de rol ya está en uso")
)

// Repository implementa la interfaz user.RoleRepository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de roles
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
