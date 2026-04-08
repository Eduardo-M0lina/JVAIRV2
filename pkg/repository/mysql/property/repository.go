package property

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/property"
)

// Repository implementa el repositorio de propiedades usando MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de propiedades
func NewRepository(db *sql.DB) property.Repository {
	return &Repository{
		db: db,
	}
}
