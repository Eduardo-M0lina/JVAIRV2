package property

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/property"
)

func (r *Repository) Update(ctx context.Context, p *property.Property) error {
	query := `
		UPDATE properties
		SET customer_id = ?, property_code = ?, street = ?, city = ?, state = ?,
		    zip = ?, notes = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		p.CustomerID,
		p.PropertyCode,
		p.Street,
		p.City,
		p.State,
		p.Zip,
		p.Notes,
		p.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute update property query",
			slog.String("error", err.Error()))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get rows affected",
			slog.String("error", err.Error()))
		return err
	}

	if rowsAffected == 0 {
		slog.WarnContext(ctx, "No property updated",
			slog.Int64("property_id", p.ID))
	}

	return nil
}
