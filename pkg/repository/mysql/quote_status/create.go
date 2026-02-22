package quote_status

import (
	"context"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/quote_status"
)

func (r *Repository) Create(ctx context.Context, qs *quote_status.QuoteStatus) error {
	query := `
		INSERT INTO quote_statuses (label, class, ` + "`order`" + `, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		qs.Label,
		qs.Class,
		qs.Order,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert quote status query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	qs.ID = id
	return nil
}
