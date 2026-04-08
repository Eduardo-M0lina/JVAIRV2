package job_visit

import (
	"context"
	"log/slog"
)

// Delete elimina una visita de trabajo (soft delete)
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Verificar que la visita existe
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrJobVisitNotFound {
			return err
		}
		slog.ErrorContext(ctx, "Failed to get job visit for delete",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	if existing.IsDeleted() {
		return ErrJobVisitNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job visit",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return err
	}

	slog.InfoContext(ctx, "Job visit deleted successfully",
		slog.Int64("id", id))

	return nil
}
