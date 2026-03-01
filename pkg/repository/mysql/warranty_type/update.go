package warranty_type

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/warranty_type"
)

func (r *Repository) Update(ctx context.Context, wt *warranty_type.WarrantyType) error {
	query := `
		UPDATE warranty_types
		SET label = ?, label_plural = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		wt.Label,
		wt.LabelPlural,
		wt.IsActive,
		wt.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty type",
			slog.String("error", err.Error()),
			slog.Int64("id", wt.ID))
		return err
	}

	return nil
}
