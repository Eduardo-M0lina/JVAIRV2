package warranty_type

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/warranty_type"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) warranty_type.Repository {
	return &Repository{db: db}
}
