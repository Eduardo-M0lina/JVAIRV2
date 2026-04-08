package quote_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, qs *QuoteStatus) error {
	if err := qs.Validate(); err != nil {
		return err
	}

	// Validar que el status existe
	_, err := uc.repo.GetByID(ctx, qs.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote status for update",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", qs.ID))
		return err
	}

	if err := uc.repo.Update(ctx, qs); err != nil {
		slog.ErrorContext(ctx, "Failed to update quote status",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", qs.ID))
		return err
	}

	slog.InfoContext(ctx, "Quote status updated successfully",
		slog.Int64("quote_status_id", qs.ID),
		slog.String("label", qs.Label))

	return nil
}
