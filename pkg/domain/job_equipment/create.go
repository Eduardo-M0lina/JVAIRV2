package job_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, equipment *JobEquipment) error {
	// Validar que el job existe
	exists, err := uc.jobChecker.JobExists(equipment.JobID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate job",
			slog.String("error", err.Error()),
			slog.Int64("job_id", equipment.JobID))
		return errors.New("invalid job_id")
	}

	if !exists {
		slog.WarnContext(ctx, "Job not found",
			slog.Int64("job_id", equipment.JobID))
		return errors.New("invalid job_id")
	}

	if err := uc.repo.Create(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to create job equipment",
			slog.String("error", err.Error()),
			slog.Int64("job_id", equipment.JobID))
		return err
	}

	slog.InfoContext(ctx, "Job equipment created successfully",
		slog.Int64("equipment_id", equipment.ID),
		slog.Int64("job_id", equipment.JobID))

	return nil
}
