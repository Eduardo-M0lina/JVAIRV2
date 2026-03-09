package job_activity_log

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_activity_log"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_activity_log.Repository {
	return &repository{db: db}
}
