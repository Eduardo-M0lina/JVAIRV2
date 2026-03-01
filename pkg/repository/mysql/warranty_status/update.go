package warranty_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/warranty_status"
)

func (r *Repository) Update(ctx context.Context, ws *warranty_status.WarrantyStatus) error {
	query := `
		UPDATE warranty_statuses
		SET label = ?, class = ?, ` + "`order`" + ` = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		ws.Label,
		ws.Class,
		ws.Order,
		ws.IsActive,
		ws.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty status",
			slog.String("error", err.Error()),
			slog.Int64("id", ws.ID))
		return err
	}

	return nil
}
