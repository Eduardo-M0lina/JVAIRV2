package job_category

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que la categoría existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job category for deletion",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", id))
		return err
	}

	// Verificar que no tenga jobs asociados
	hasJobs, err := uc.repo.HasJobs(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check job category jobs",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", id))
		return err
	}

	if hasJobs {
		slog.WarnContext(ctx, "Cannot delete job category with jobs",
			slog.Int64("job_category_id", id))
		return ErrJobCategoryInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job category",
			slog.String("error", err.Error()),
			slog.Int64("job_category_id", id))
		return err
	}

	slog.InfoContext(ctx, "Job category deleted successfully",
		slog.Int64("job_category_id", id))

	return nil
}
