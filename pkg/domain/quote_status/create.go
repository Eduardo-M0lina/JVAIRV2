package quote_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, qs *QuoteStatus) error {
	if err := qs.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, qs); err != nil {
		slog.ErrorContext(ctx, "Failed to create quote status",
			slog.String("error", err.Error()),
			slog.String("label", qs.Label))
		return err
	}

	slog.InfoContext(ctx, "Quote status created successfully",
		slog.Int64("quote_status_id", qs.ID),
		slog.String("label", qs.Label))

	return nil
}
