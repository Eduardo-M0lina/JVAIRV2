package quote

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Quote, error) {
	q, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get quote by ID",
			slog.String("error", err.Error()),
			slog.Int64("quote_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Quote retrieved successfully",
		slog.Int64("quote_id", id))

	return q, nil
}
