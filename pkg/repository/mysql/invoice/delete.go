package invoice

import (
	"context"
	"log/slog"
)

// Delete elimina una factura (soft delete)
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE invoices SET deleted_at = NOW(), updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete invoice",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
