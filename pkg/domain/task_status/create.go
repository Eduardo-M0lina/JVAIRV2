package task_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, status *TaskStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, status); err != nil {
		slog.ErrorContext(ctx, "Failed to create task status",
			slog.String("error", err.Error()),
			slog.String("label", status.Label))
		return err
	}

	slog.InfoContext(ctx, "Task status created successfully",
		slog.Int64("task_status_id", status.ID),
		slog.String("label", status.Label))

	return nil
}
