package email_template

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/email_template"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) email_template.Repository {
	return &Repository{db: db}
}
