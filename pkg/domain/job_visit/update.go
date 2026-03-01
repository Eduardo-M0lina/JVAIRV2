package job_visit

import (
	"context"
	"fmt"
	"log/slog"
)

// Update actualiza una visita de trabajo existente
func (uc *UseCase) Update(ctx context.Context, visit *JobVisit) error {
	if err := visit.ValidateUpdate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Verificar que la visita existe
	existing, err := uc.repo.GetByID(ctx, visit.ID)
	if err != nil {
		if err == ErrJobVisitNotFound {
			return err
		}
		slog.ErrorContext(ctx, "Failed to get job visit for update",
			slog.String("error", err.Error()),
			slog.Int64("id", visit.ID))
		return err
	}

	if existing.IsDeleted() {
		return ErrJobVisitNotFound
	}

	// Verificar que el usuario existe
	exists, err := uc.userChecker.UserExists(ctx, visit.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check user existence",
			slog.String("error", err.Error()),
			slog.Int64("userId", visit.UserID))
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	if err := uc.repo.Update(ctx, visit); err != nil {
		slog.ErrorContext(ctx, "Failed to update job visit",
			slog.String("error", err.Error()),
			slog.Int64("id", visit.ID))
		return err
	}

	slog.InfoContext(ctx, "Job visit updated successfully",
		slog.Int64("id", visit.ID))

	return nil
}
