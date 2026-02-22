package quote_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/quote_status"
)

func (r *Repository) Update(ctx context.Context, qs *quote_status.QuoteStatus) error {
	query := `
		UPDATE quote_statuses
		SET label = ?, class = ?, ` + "`order`" + ` = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		qs.Label,
		qs.Class,
		qs.Order,
		qs.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update quote status",
			slog.String("error", err.Error()),
			slog.Int64("id", qs.ID))
		return err
	}

	return nil
}
