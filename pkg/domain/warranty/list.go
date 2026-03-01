package warranty

import (
	"context"
	"log/slog"
)

// List obtiene una lista paginada de garantías
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Warranty, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	warranties, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranties",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("pageSize", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Warranties listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return warranties, total, nil
}
