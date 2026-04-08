package quote

import (
	"context"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Quote, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	quotes, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list quotes",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("pageSize", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Quotes listed successfully",
		slog.Int64("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return quotes, total, nil
}
