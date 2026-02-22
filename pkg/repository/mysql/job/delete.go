package job

import (
	"context"
	"log/slog"
)

// Delete elimina un job (soft delete)
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE jobs SET deleted_at = NOW(), updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
