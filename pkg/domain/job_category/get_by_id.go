package job_category

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*JobCategory, error) {
	category, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job category by ID",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Job category retrieved successfully",
		slog.Int64("job_category_id", id))

	return category, nil
}
