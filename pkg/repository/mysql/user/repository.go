package user

import (
	"database/sql"
	"errors"
)

// Errores comunes del repositorio
var (
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrDuplicateEmail     = errors.New("email ya está en uso")
)

// Repository implementa la interfaz user.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de usuarios
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
