package password_reset

import (
	"database/sql"

	domainPasswordReset "github.com/angumol/jvairv2/pkg/domain/password_reset"
)

// Repository implementa la interfaz domainPasswordReset.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de password reset
func NewRepository(db *sql.DB) domainPasswordReset.Repository {
	return &Repository{
		db: db,
	}
}
