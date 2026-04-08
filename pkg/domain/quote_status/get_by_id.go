package quote_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*QuoteStatus, error) {
	qs, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote status by ID",
			slog.String("error", err.Error()),
			slog.Int64("quote_status_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Quote status retrieved successfully",
		slog.Int64("quote_status_id", id))

	return qs, nil
}
