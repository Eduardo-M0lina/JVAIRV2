package file

import (
	"context"
	"log/slog"
)

// Delete elimina un registro de archivo de la base de datos (hard delete)
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM files WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete file record",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}
