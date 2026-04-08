package job_visit

import (
	"context"
	"log/slog"
)

// Delete realiza un soft delete de una visita de trabajo
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE job_visits SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete job visit",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}
