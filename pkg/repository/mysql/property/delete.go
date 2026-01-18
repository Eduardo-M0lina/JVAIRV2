package property

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE properties
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute delete property query",
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
		slog.WarnContext(ctx, "No property deleted",
			slog.Int64("property_id", id))
	}

	return nil
}
