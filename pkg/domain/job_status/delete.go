package job_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que el estado existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job status for deletion",
			slog.String("error", err.Error()),
			slog.Int64("job_status_id", id))
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job status",
			slog.String("error", err.Error()),
			slog.Int64("job_status_id", id))
		return err
	}

	slog.InfoContext(ctx, "Job status deleted successfully",
		slog.Int64("job_status_id", id))

	return nil
}
