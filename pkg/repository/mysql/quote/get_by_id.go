package quote

import (
	"context"
	"database/sql"
	"log/slog"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*domainQuote.Quote, error) {
	query := `
		SELECT id, job_id, quote_number, quote_status_id, amount, description, notes, created_at, updated_at
		FROM quotes
		WHERE id = ? AND deleted_at IS NULL
	`

	q := &domainQuote.Quote{}
	var description, notes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&q.ID,
		&q.JobID,
		&q.QuoteNumber,
		&q.QuoteStatusID,
		&q.Amount,
		&description,
		&notes,
		&q.CreatedAt,
		&q.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainQuote.ErrQuoteNotFound
		}
		slog.ErrorContext(ctx, "Failed to get quote by ID",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return nil, err
	}

	if description.Valid {
		q.Description = &description.String
	}
	if notes.Valid {
		q.Notes = &notes.String
	}

	return q, nil
}
