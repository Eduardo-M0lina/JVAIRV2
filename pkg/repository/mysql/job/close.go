package job

import (
	"context"
	"log/slog"
)

// Close cierra un job actualizando closed=true y el job_status_id
func (r *Repository) Close(ctx context.Context, id int64, jobStatusID int64) error {
	query := `UPDATE jobs SET closed = 1, job_status_id = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, jobStatusID, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to close job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
