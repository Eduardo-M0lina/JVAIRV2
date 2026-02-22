package task_status

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/task_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) task_status.Repository {
	return &Repository{db: db}
}
