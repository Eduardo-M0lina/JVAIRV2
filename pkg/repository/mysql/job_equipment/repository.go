package job_equipment

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

// Repository implementa el repositorio de equipos de trabajo usando MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de equipos de trabajo
func NewRepository(db *sql.DB) job_equipment.Repository {
	return &Repository{
		db: db,
	}
}
