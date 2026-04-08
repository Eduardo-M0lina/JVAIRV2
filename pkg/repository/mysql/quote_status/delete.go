package quote_status

import (
	"context"
	"log/slog"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM quote_statuses WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete quote status",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	return nil
}

func (r *Repository) HasQuotes(ctx context.Context, id int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM quotes WHERE quote_status_id = ? AND deleted_at IS NULL",
		id,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
