package warranty

import (
	"database/sql"

	domainWarranty "github.com/your-org/jvairv2/pkg/domain/warranty"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domainWarranty.Repository {
	return &Repository{db: db}
}
