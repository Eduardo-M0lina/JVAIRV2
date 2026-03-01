package job_visit

import "database/sql"

// Repository implementa el repositorio MySQL para visitas de trabajo
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
