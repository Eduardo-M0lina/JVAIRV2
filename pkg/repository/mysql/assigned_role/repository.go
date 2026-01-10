package assigned_role

import (
	"database/sql"
	"errors"
)

// Errores comunes del repositorio
var (
	ErrAssignedRoleNotFound = errors.New("asignación de rol no encontrada")
	ErrDuplicateAssignment  = errors.New("el rol ya está asignado a esta entidad")
)

// Repository implementa la interfaz assigned_role.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de asignaciones de roles
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
