package warranty_claim_type

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/warranty_claim_type"
)

func (r *Repository) Update(ctx context.Context, wct *warranty_claim_type.WarrantyClaimType) error {
	query := `
		UPDATE warranty_claim_types
		SET label = ?, label_plural = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		wct.Label,
		wct.LabelPlural,
		wct.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim type",
			slog.String("error", err.Error()),
			slog.Int64("id", wct.ID))
		return err
	}

	return nil
}
