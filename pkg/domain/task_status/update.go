package task_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, status *TaskStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	// Validar que el estado existe
	_, err := uc.repo.GetByID(ctx, status.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get task status for update",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", status.ID))
		return err
	}

	if err := uc.repo.Update(ctx, status); err != nil {
		slog.ErrorContext(ctx, "Failed to update task status",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", status.ID))
		return err
	}

	slog.InfoContext(ctx, "Task status updated successfully",
		slog.Int64("task_status_id", status.ID),
		slog.String("label", status.Label))

	return nil
}
