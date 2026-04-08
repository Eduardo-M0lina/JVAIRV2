package quote

import (
	"context"
	"log/slog"
)

// Create crea una nueva cotización
func (uc *UseCase) Create(ctx context.Context, q *Quote) error {
	if err := q.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que el job existe
	if _, err := uc.jobRepo.GetByID(ctx, q.JobID); err != nil {
		slog.ErrorContext(ctx, "Invalid job",
			slog.Int64("jobId", q.JobID),
			slog.String("error", err.Error()))
		return ErrInvalidJob
	}

	// Verificar que el quote status existe
	if _, err := uc.quoteStatusRepo.GetByID(ctx, q.QuoteStatusID); err != nil {
		slog.ErrorContext(ctx, "Invalid quote status",
			slog.Int64("quoteStatusId", q.QuoteStatusID),
			slog.String("error", err.Error()))
		return ErrInvalidQuoteStatus
	}

	if err := uc.repo.Create(ctx, q); err != nil {
		slog.ErrorContext(ctx, "Failed to create quote",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Quote created successfully",
		slog.Int64("id", q.ID))

	return nil
}
