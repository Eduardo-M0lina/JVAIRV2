package job_category

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_category"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_category.Repository {
	return &Repository{db: db}
}
