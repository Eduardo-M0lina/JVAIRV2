package payroll

import (
	"database/sql"

	"github.com/angumol/jvairv2/pkg/domain/payroll"
)

type repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio de payroll
func NewRepository(db *sql.DB) payroll.Repository {
	return &repository{db: db}
}
