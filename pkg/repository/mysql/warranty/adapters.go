package warranty

import (
	"context"
	"database/sql"

	domainWarranty "github.com/your-org/jvairv2/pkg/domain/warranty"
)

// JobCheckerAdapter adapta la verificación de jobs para el use case de warranties
type JobCheckerAdapter struct {
	db *sql.DB
}

func NewJobCheckerAdapter(db *sql.DB) domainWarranty.JobChecker {
	return &JobCheckerAdapter{db: db}
}

func (a *JobCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWarranty.ErrInvalidJob
	}
	return true, nil
}

// WarrantyTypeCheckerAdapter adapta la verificación de tipos de garantía
type WarrantyTypeCheckerAdapter struct {
	db *sql.DB
}

func NewWarrantyTypeCheckerAdapter(db *sql.DB) domainWarranty.WarrantyTypeChecker {
	return &WarrantyTypeCheckerAdapter{db: db}
}

func (a *WarrantyTypeCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranty_types WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWarranty.ErrInvalidWarrantyType
	}
	return true, nil
}

// WarrantyStatusCheckerAdapter adapta la verificación de estados de garantía
type WarrantyStatusCheckerAdapter struct {
	db *sql.DB
}

func NewWarrantyStatusCheckerAdapter(db *sql.DB) domainWarranty.WarrantyStatusChecker {
	return &WarrantyStatusCheckerAdapter{db: db}
}

func (a *WarrantyStatusCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranty_statuses WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWarranty.ErrInvalidWarrantyStatus
	}
	return true, nil
}
