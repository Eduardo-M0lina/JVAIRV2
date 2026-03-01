package job_visit

import (
	"context"
	"fmt"
	"log/slog"
)

// Create crea una nueva visita de trabajo
func (uc *UseCase) Create(ctx context.Context, visit *JobVisit) (int64, error) {
	if err := visit.ValidateCreate(); err != nil {
		return 0, fmt.Errorf("validation error: %w", err)
	}

	// Verificar que el job existe
	exists, err := uc.jobChecker.JobExists(ctx, visit.JobID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check job existence",
			slog.String("error", err.Error()),
			slog.Int64("jobId", visit.JobID))
		return 0, err
	}
	if !exists {
		return 0, ErrJobNotFound
	}

	// Verificar que el usuario existe
	exists, err = uc.userChecker.UserExists(ctx, visit.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check user existence",
			slog.String("error", err.Error()),
			slog.Int64("userId", visit.UserID))
		return 0, err
	}
	if !exists {
		return 0, ErrUserNotFound
	}

	id, err := uc.repo.Create(ctx, visit)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create job visit",
			slog.String("error", err.Error()))
		return 0, err
	}

	slog.InfoContext(ctx, "Job visit created successfully",
		slog.Int64("id", id),
		slog.Int64("jobId", visit.JobID))

	return id, nil
}
