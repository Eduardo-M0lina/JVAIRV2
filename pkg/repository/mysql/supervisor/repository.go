package supervisor

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/supervisor"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) supervisor.Repository {
	return &Repository{db: db}
}
