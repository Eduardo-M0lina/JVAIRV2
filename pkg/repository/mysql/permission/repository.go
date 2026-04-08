package permission

import (
	"database/sql"
	"errors"
)

// Errores comunes del repositorio
var (
	ErrPermissionNotFound  = errors.New("permiso no encontrado")
	ErrDuplicatePermission = errors.New("el permiso ya existe para esta entidad y ability")
)

// Repository implementa la interfaz permission.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de permisos
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
