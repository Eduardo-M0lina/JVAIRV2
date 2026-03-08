package job_resident

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_resident"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_resident.Repository {
	return &repository{db: db}
}
