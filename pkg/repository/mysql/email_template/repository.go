package email_template

import (
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/email_template"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) email_template.Repository {
	return &Repository{db: db}
}
