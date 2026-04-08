package warranty_claim_type

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_type"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*warranty_claim_type.WarrantyClaimType, error) {
	query := `
		SELECT id, label, label_plural, created_at, updated_at
		FROM warranty_claim_types
		WHERE id = ?
	`

	wct := &warranty_claim_type.WarrantyClaimType{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wct.ID,
		&wct.Label,
		&wct.LabelPlural,
		&wct.CreatedAt,
		&wct.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, warranty_claim_type.ErrWarrantyClaimTypeNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty claim type by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return wct, nil
}
