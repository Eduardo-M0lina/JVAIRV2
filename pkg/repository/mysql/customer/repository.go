package customer

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/customer"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) customer.Repository {
	return &Repository{db: db}
}
