package warranty_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*warranty_status.WarrantyStatus, error) {
	query := `
		SELECT id, label, class, ` + "`order`" + `, is_active, created_at, updated_at
		FROM warranty_statuses
		WHERE id = ?
	`

	ws := &warranty_status.WarrantyStatus{}
	var class sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ws.ID,
		&ws.Label,
		&class,
		&ws.Order,
		&ws.IsActive,
		&ws.CreatedAt,
		&ws.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, warranty_status.ErrWarrantyStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if class.Valid {
		ws.Class = &class.String
	}

	return ws, nil
}
