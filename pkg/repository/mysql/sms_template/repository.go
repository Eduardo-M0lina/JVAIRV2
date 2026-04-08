package sms_template

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/sms_template"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) sms_template.Repository {
	return &Repository{db: db}
}
