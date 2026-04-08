package job_priority

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que la prioridad existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job priority for deletion",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", id))
		return err
	}

	// Verificar que no tenga jobs asociados
	hasJobs, err := uc.repo.HasJobs(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check job priority jobs",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", id))
		return err
	}

	if hasJobs {
		slog.WarnContext(ctx, "Cannot delete job priority with jobs",
			slog.Int64("job_priority_id", id))
		return ErrJobPriorityInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job priority",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", id))
		return err
	}

	slog.InfoContext(ctx, "Job priority deleted successfully",
		slog.Int64("job_priority_id", id))

	return nil
}
