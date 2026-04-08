package job_priority

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, priority *JobPriority) error {
	if err := priority.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, priority); err != nil {
		slog.ErrorContext(ctx, "Failed to create job priority",
			slog.String("error", err.Error()),
			slog.String("label", priority.Label))
		return err
	}

	slog.InfoContext(ctx, "Job priority created successfully",
		slog.Int64("job_priority_id", priority.ID),
		slog.String("label", priority.Label))

	return nil
}
