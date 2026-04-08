package job_task

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/job_task"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) job_task.Repository {
	return &repository{db: db}
}
