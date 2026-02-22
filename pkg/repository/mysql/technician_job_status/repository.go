package technician_job_status

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/technician_job_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) technician_job_status.Repository {
	return &Repository{db: db}
}
