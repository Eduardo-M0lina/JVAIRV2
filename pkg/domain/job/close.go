package job

import (
	"context"
	"log/slog"
)

// Close cierra un job con un status específico
func (uc *UseCase) Close(ctx context.Context, id int64, jobStatusID int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Job not found for close",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return ErrJobNotFound
	}

	if existing.IsDeleted() {
		return ErrJobNotFound
	}

	if existing.IsClosed() {
		return ErrJobAlreadyClosed
	}

	// Verificar que el status existe
	if jobStatusID > 0 {
		if _, err := uc.jobStatusRepo.GetByID(ctx, jobStatusID); err != nil {
			return ErrInvalidJobStatus
		}
	}

	if err := uc.repo.Close(ctx, id, jobStatusID); err != nil {
		slog.ErrorContext(ctx, "Failed to close job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Job closed successfully",
		slog.Int64("id", id))

	return nil
}
