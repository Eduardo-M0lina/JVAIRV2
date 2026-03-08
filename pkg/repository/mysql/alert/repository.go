package alert

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/alert"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) alert.Repository {
	return &repository{db: db}
}
