package quote

import (
	"context"
	"log/slog"
)

// Update actualiza una cotización existente
func (uc *UseCase) Update(ctx context.Context, q *Quote) error {
	if err := q.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que la cotización existe
	existing, err := uc.repo.GetByID(ctx, q.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote for update",
			slog.String("error", err.Error()),
			slog.Int64("quote_id", q.ID))
		return err
	}

	// Verificar que el job existe si cambió
	if q.JobID != existing.JobID {
		if _, err := uc.jobRepo.GetByID(ctx, q.JobID); err != nil {
			slog.ErrorContext(ctx, "Invalid job",
				slog.Int64("jobId", q.JobID),
				slog.String("error", err.Error()))
			return ErrInvalidJob
		}
	}

	// Verificar que el quote status existe si cambió
	if q.QuoteStatusID != existing.QuoteStatusID {
		if _, err := uc.quoteStatusRepo.GetByID(ctx, q.QuoteStatusID); err != nil {
			slog.ErrorContext(ctx, "Invalid quote status",
				slog.Int64("quoteStatusId", q.QuoteStatusID),
				slog.String("error", err.Error()))
			return ErrInvalidQuoteStatus
		}
	}

	if err := uc.repo.Update(ctx, q); err != nil {
		slog.ErrorContext(ctx, "Failed to update quote",
			slog.String("error", err.Error()),
			slog.Int64("quote_id", q.ID))
		return err
	}

	slog.InfoContext(ctx, "Quote updated successfully",
		slog.Int64("quote_id", q.ID))

	return nil
}
