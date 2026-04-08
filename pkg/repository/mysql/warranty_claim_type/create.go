package warranty_claim_type

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_type"
)

func (r *Repository) Create(ctx context.Context, wct *warranty_claim_type.WarrantyClaimType) error {
	query := `
		INSERT INTO warranty_claim_types (label, label_plural, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		wct.Label,
		wct.LabelPlural,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty claim type query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	wct.ID = id
	return nil
}
