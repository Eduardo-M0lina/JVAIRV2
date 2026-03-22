package password_history

import (
	"database/sql"

	domainPasswordHistory "github.com/your-org/jvairv2/pkg/domain/password_history"
)

// Repository implementa la interfaz domainPasswordHistory.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de password history
func NewRepository(db *sql.DB) domainPasswordHistory.Repository {
	return &Repository{
		db: db,
	}
}
