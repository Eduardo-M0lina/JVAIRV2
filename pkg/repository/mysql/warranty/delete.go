package warranty

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE warranties SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to soft delete warranty",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}
