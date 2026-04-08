package warranty_equipment

import (
	"context"
	"database/sql"

	domainWE "github.com/angumol/jvairv2/pkg/domain/warranty_equipment"
)

// WarrantyCheckerAdapter adapta la verificación de warranties para el use case de warranty equipment
type WarrantyCheckerAdapter struct {
	db *sql.DB
}

func NewWarrantyCheckerAdapter(db *sql.DB) domainWE.WarrantyChecker {
	return &WarrantyCheckerAdapter{db: db}
}

func (a *WarrantyCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranties WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWE.ErrInvalidWarranty
	}
	return true, nil
}
