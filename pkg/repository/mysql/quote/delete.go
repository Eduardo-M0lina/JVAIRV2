package quote

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE quotes SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete quote",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}
