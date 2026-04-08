package warranty_claim_status

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_status"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) warranty_claim_status.Repository {
	return &Repository{db: db}
}
