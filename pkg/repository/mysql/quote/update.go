package quote

import (
	"context"
	"log/slog"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

func (r *Repository) Update(ctx context.Context, q *domainQuote.Quote) error {
	query := `
		UPDATE quotes
		SET job_id = ?, quote_number = ?, quote_status_id = ?, amount = ?,
		    description = ?, notes = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		q.JobID,
		q.QuoteNumber,
		q.QuoteStatusID,
		q.Amount,
		q.Description,
		q.Notes,
		q.ID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update quote",
			slog.String("error", err.Error()),
			slog.Int64("id", q.ID))
		return err
	}

	return nil
}
