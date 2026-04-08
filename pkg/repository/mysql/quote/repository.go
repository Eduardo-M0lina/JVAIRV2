package quote

import (
	"database/sql"

	domainQuote "github.com/angumol/jvairv2/pkg/domain/quote"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domainQuote.Repository {
	return &Repository{db: db}
}
