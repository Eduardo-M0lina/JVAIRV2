package job_priority

import (
	"context"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*JobPriority, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	priorities, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job priorities",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("pageSize", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Job priorities listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return priorities, total, nil
}
