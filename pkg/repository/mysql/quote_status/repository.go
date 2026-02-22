package quote_status

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/quote_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) quote_status.Repository {
	return &Repository{db: db}
}
