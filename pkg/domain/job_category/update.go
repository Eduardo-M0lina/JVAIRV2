package job_category

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, category *JobCategory) error {
	if err := category.Validate(); err != nil {
		return err
	}

	// Validar que la categoría existe
	_, err := uc.repo.GetByID(ctx, category.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job category for update",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", category.ID))
		return err
	}

	if err := uc.repo.Update(ctx, category); err != nil {
		slog.ErrorContext(ctx, "Failed to update job category",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", category.ID))
		return err
	}

	slog.InfoContext(ctx, "Job category updated successfully",
		slog.Int64("job_category_id", category.ID),
		slog.String("label", category.Label))

	return nil
}
