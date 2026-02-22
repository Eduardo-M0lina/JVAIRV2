package task_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*TaskStatus, error) {
	status, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get task status by ID",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Task status retrieved successfully",
		slog.Int64("task_status_id", id))

	return status, nil
}
