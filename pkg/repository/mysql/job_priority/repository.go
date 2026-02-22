package job_priority

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_priority"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_priority.Repository {
	return &Repository{db: db}
}
