package warranty_type

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/warranty_type"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*warranty_type.WarrantyType, error) {
	query := `
		SELECT id, label, label_plural, is_active, created_at, updated_at
		FROM warranty_types
		WHERE id = ?
	`

	wt := &warranty_type.WarrantyType{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wt.ID,
		&wt.Label,
		&wt.LabelPlural,
		&wt.IsActive,
		&wt.CreatedAt,
		&wt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, warranty_type.ErrWarrantyTypeNotFound
		}
		slog.ErrorContext(ctx, "Failed to get warranty type by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	return wt, nil
}
