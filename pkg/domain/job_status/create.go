package job_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, status *JobStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, status); err != nil {
		slog.ErrorContext(ctx, "Failed to create job status",
			slog.String("error", err.Error()),
			slog.String("label", status.Label))
		return err
	}

	slog.InfoContext(ctx, "Job status created successfully",
		slog.Int64("job_status_id", status.ID),
		slog.String("label", status.Label))

	return nil
}
