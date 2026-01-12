package user

import (
	"database/sql"

	domainUser "github.com/your-org/jvairv2/pkg/domain/user"
)

// Usar los errores del dominio
var (
	ErrUserNotFound       = domainUser.ErrUserNotFound
	ErrInvalidCredentials = domainUser.ErrInvalidCredentials
	ErrDuplicateEmail     = domainUser.ErrDuplicateEmail
)

// Repository implementa la interfaz domainUser.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de usuarios
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
