package job

import (
	"context"
	"log/slog"
)

// Delete elimina un job (soft delete)
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Job not found for delete",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return ErrJobNotFound
	}

	if existing.IsDeleted() {
		return ErrJobNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Job deleted successfully",
		slog.Int64("id", id))

	return nil
}
