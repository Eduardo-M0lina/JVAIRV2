package warranty_claim_status

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_claim_status"
)

func (r *Repository) Update(ctx context.Context, wcs *warranty_claim_status.WarrantyClaimStatus) error {
	query := `
		UPDATE warranty_claim_statuses
		SET label = ?, class = ?, ` + "`order`" + ` = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		wcs.Label,
		wcs.Class,
		wcs.Order,
		wcs.IsActive,
		wcs.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim status",
			slog.String("error", err.Error()),
			slog.Int64("id", wcs.ID))
		return err
	}

	return nil
}
