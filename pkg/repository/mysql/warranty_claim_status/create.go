package warranty_claim_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/warranty_claim_status"
)

func (r *Repository) Create(ctx context.Context, wcs *warranty_claim_status.WarrantyClaimStatus) error {
	query := `
		INSERT INTO warranty_claim_statuses (label, class, ` + "`order`" + `, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		wcs.Label,
		wcs.Class,
		wcs.Order,
		wcs.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty claim status query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	wcs.ID = id
	return nil
}
