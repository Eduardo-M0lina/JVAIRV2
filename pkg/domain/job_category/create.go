package job_category

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, category *JobCategory) error {
	if err := category.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, category); err != nil {
		slog.ErrorContext(ctx, "Failed to create job category",
			slog.String("error", err.Error()),
			slog.String("label", category.Label))
		return err
	}

	slog.InfoContext(ctx, "Job category created successfully",
		slog.Int64("job_category_id", category.ID),
		slog.String("label", category.Label))

	return nil
}
