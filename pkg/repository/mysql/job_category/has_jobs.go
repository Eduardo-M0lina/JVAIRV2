package job_category

import (
	"context"
	"log/slog"
)

func (r *Repository) HasJobs(ctx context.Context, id int64) (bool, error) {
	query := `SELECT COUNT(*) FROM jobs WHERE job_category_id = ?`

	var count int
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		slog.ErrorContext(ctx, "Failed to check job category jobs",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return false, err
	}

	return count > 0, nil
}
