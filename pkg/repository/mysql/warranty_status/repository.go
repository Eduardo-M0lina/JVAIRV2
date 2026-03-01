package warranty_status

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/warranty_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) warranty_status.Repository {
	return &Repository{db: db}
}
