package job_rate

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_rate"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_rate.Repository {
	return &repository{db: db}
}
