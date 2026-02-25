package property_equipment

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/property_equipment"
)

// Repository implementa el repositorio de equipos de propiedad usando MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de equipos de propiedad
func NewRepository(db *sql.DB) property_equipment.Repository {
	return &Repository{
		db: db,
	}
}
