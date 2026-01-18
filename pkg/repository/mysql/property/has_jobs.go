package property

import (
	"context"
	"log/slog"
)

func (r *Repository) HasJobs(ctx context.Context, id int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM jobs
		WHERE property_id = ? AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check property jobs",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return false, err
	}

	return count > 0, nil
}
