package warranty_claim_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*warranty_claim_status.WarrantyClaimStatus, error) {
	query := `
		SELECT id, label, class, ` + "`order`" + `, is_active, created_at, updated_at
		FROM warranty_claim_statuses
		WHERE id = ?
	`

	wcs := &warranty_claim_status.WarrantyClaimStatus{}
	var class sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wcs.ID,
		&wcs.Label,
		&class,
		&wcs.Order,
		&wcs.IsActive,
		&wcs.CreatedAt,
		&wcs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, warranty_claim_status.ErrWarrantyClaimStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty claim status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if class.Valid {
		wcs.Class = &class.String
	}

	return wcs, nil
}
