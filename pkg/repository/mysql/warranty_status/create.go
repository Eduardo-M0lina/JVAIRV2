package warranty_status

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_status"
)

func (r *Repository) Create(ctx context.Context, ws *warranty_status.WarrantyStatus) error {
	query := `
		INSERT INTO warranty_statuses (label, class, ` + "`order`" + `, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		ws.Label,
		ws.Class,
		ws.Order,
		ws.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty status query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	ws.ID = id
	return nil
}
