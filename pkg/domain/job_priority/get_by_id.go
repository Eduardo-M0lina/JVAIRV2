package job_priority

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*JobPriority, error) {
	priority, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job priority by ID",
			slog.String("error", err.Error()),
			slog.Int64("job_priority_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Job priority retrieved successfully",
		slog.Int64("job_priority_id", id))

	return priority, nil
}
