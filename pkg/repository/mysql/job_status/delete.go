package job_status

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM job_statuses WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete job status",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}
