package alert

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/alert"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) alert.Repository {
	return &repository{db: db}
}
