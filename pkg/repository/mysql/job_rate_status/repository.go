package job_rate_status

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_rate_status"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_rate_status.Repository {
	return &repository{db: db}
}
