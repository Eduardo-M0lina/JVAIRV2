package settings

import (
	"database/sql"

	domainSettings "github.com/your-org/jvairv2/pkg/domain/settings"
)

// Repository implementa el repositorio de configuraciones usando MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de configuraciones
func NewRepository(db *sql.DB) domainSettings.Repository {
	return &Repository{
		db: db,
	}
}
