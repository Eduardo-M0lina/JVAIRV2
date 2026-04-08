package warranty_type

import (
	"context"
	"log/slog"

	"github.com/angumol/jvairv2/pkg/domain/warranty_type"
)

func (r *Repository) Create(ctx context.Context, wt *warranty_type.WarrantyType) error {
	query := `
		INSERT INTO warranty_types (label, label_plural, is_active, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		wt.Label,
		wt.LabelPlural,
		wt.IsActive,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert warranty type query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	wt.ID = id
	return nil
}
