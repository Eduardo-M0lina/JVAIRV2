package warranty_claim

import (
	"database/sql"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domainWC.Repository {
	return &Repository{db: db}
}
