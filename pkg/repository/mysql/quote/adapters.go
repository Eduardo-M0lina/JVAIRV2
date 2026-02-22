package quote

import (
	"context"
	"database/sql"
	"log/slog"

	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
)

// JobCheckerAdapter adapta la verificación de existencia de jobs
type JobCheckerAdapter struct {
	db *sql.DB
}

func NewJobCheckerAdapter(db *sql.DB) domainQuote.JobChecker {
	return &JobCheckerAdapter{db: db}
}

func (a *JobCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ? AND deleted_at IS NULL)",
		id,
	).Scan(&exists)
	if err != nil || !exists {
		slog.ErrorContext(ctx, "Job not found",
			slog.Int64("jobId", id))
		return nil, domainQuote.ErrInvalidJob
	}
	return true, nil
}

// QuoteStatusCheckerAdapter adapta la verificación de existencia de estados de cotización
type QuoteStatusCheckerAdapter struct {
	db *sql.DB
}

func NewQuoteStatusCheckerAdapter(db *sql.DB) domainQuote.QuoteStatusChecker {
	return &QuoteStatusCheckerAdapter{db: db}
}

func (a *QuoteStatusCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM quote_statuses WHERE id = ?)",
		id,
	).Scan(&exists)
	if err != nil || !exists {
		slog.ErrorContext(ctx, "Quote status not found",
			slog.Int64("quoteStatusId", id))
		return nil, domainQuote.ErrInvalidQuoteStatus
	}
	return true, nil
}
