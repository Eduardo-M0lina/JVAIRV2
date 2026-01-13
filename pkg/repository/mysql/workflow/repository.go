package workflow

import (
	"database/sql"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// Repository implementa el repositorio de workflows usando MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de workflows
func NewRepository(db *sql.DB) domainWorkflow.Repository {
	return &Repository{
		db: db,
	}
}
