package warranty_claim_type

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_type"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) warranty_claim_type.Repository {
	return &Repository{db: db}
}
