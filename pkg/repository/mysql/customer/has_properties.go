package customer

import (
	"context"
	"log/slog"
)

func (r *Repository) HasProperties(ctx context.Context, id int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM properties
		WHERE customer_id = ? AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check customer properties",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return false, err
	}

	return count > 0, nil
}
