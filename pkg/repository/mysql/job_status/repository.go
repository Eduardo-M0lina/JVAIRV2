package job_status

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_status.Repository {
	return &Repository{db: db}
}
