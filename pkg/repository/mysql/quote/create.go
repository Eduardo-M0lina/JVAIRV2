package quote

import (
	"context"
	"log/slog"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

func (r *Repository) Create(ctx context.Context, q *domainQuote.Quote) error {
	query := `
		INSERT INTO quotes (job_id, quote_number, quote_status_id, amount, description, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		q.JobID,
		q.QuoteNumber,
		q.QuoteStatusID,
		q.Amount,
		q.Description,
		q.Notes,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute insert quote query",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	q.ID = id
	return nil
}
