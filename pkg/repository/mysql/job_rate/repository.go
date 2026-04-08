package job_rate

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_rate"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_rate.Repository {
	return &repository{db: db}
}
