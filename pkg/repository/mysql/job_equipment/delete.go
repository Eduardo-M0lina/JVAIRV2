package job_equipment

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM job_equipment WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute delete job_equipment query",
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
		slog.WarnContext(ctx, "No job equipment deleted",
			slog.Int64("equipment_id", id))
	}

	return nil
}
