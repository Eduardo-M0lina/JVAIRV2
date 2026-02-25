package supervisor

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) supervisor.Repository {
	return &Repository{db: db}
}
