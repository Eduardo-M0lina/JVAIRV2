package warranty_equipment

import (
	"database/sql"

	domainWE "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domainWE.Repository {
	return &Repository{db: db}
}
