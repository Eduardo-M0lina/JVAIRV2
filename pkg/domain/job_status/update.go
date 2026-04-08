package job_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, status *JobStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	// Validar que el estado existe
	_, err := uc.repo.GetByID(ctx, status.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job status for update",
			slog.String("error", err.Error()),
			slog.Int64("job_status_id", status.ID))
		return err
	}

	if err := uc.repo.Update(ctx, status); err != nil {
		slog.ErrorContext(ctx, "Failed to update job status",
			slog.String("error", err.Error()),
			slog.Int64("job_status_id", status.ID))
		return err
	}

	slog.InfoContext(ctx, "Job status updated successfully",
		slog.Int64("job_status_id", status.ID),
		slog.String("label", status.Label))

	return nil
}
