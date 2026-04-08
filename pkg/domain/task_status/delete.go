package task_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que el estado existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get task status for deletion",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", id))
		return err
	}

	// Verificar que no tenga job_tasks asociados
	hasJobTasks, err := uc.repo.HasJobTasks(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check task status job tasks",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", id))
		return err
	}

	if hasJobTasks {
		slog.WarnContext(ctx, "Cannot delete task status with job tasks",
			slog.Int64("task_status_id", id))
		return ErrTaskStatusInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete task status",
			slog.String("error", err.Error()),
			slog.Int64("task_status_id", id))
		return err
	}

	slog.InfoContext(ctx, "Task status deleted successfully",
		slog.Int64("task_status_id", id))

	return nil
}
