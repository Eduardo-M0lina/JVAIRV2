package job_priority

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, priority *JobPriority) error {
	if err := priority.Validate(); err != nil {
		return err
	}

	// Validar que la prioridad existe
	_, err := uc.repo.GetByID(ctx, priority.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job priority for update",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", priority.ID))
		return err
	}

	if err := uc.repo.Update(ctx, priority); err != nil {
		slog.ErrorContext(ctx, "Failed to update job priority",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", priority.ID))
		return err
	}

	slog.InfoContext(ctx, "Job priority updated successfully",
		slog.Int64("job_priority_id", priority.ID),
		slog.String("label", priority.Label))

	return nil
}
