package quote_status

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/your-org/jvairv2/pkg/domain/quote_status"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*quote_status.QuoteStatus, error) {
	query := `
		SELECT id, label, class, ` + "`order`" + `, created_at, updated_at
		FROM quote_statuses
		WHERE id = ?
	`

	qs := &quote_status.QuoteStatus{}
	var class sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&qs.ID,
		&qs.Label,
		&class,
		&qs.Order,
		&qs.CreatedAt,
		&qs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, quote_status.ErrQuoteStatusNotFound
		}
		slog.ErrorContext(ctx, "Failed to get quote status by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if class.Valid {
		qs.Class = &class.String
	}

	return qs, nil
}
