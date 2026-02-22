package task_status

import (
	"context"
	"log/slog"
)

func (r *Repository) HasJobTasks(ctx context.Context, id int64) (bool, error) {
	query := `SELECT COUNT(*) FROM job_tasks WHERE task_status_id = ?`

	var count int
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		slog.ErrorContext(ctx, "Failed to check task status job tasks",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return false, err
	}

	return count > 0, nil
}
