package technician_job_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*TechnicianJobStatus, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	statuses, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list technician job statuses",
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("pageSize", pageSize))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Technician job statuses listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return statuses, total, nil
}
