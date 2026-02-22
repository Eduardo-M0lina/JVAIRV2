package quote

import (
	"database/sql"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domainQuote.Repository {
	return &Repository{db: db}
}
