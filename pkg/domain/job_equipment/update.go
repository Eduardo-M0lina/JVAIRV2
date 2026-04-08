package job_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, equipment *JobEquipment) error {
	// Validar que el equipo existe
	existing, err := uc.repo.GetByID(ctx, equipment.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job equipment for update",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", equipment.ID))
		return err
	}

	// Validar que el equipo pertenece al job
	if existing.JobID != equipment.JobID {
		slog.WarnContext(ctx, "Equipment does not belong to job",
			slog.Int64("equipment_id", equipment.ID),
			slog.Int64("job_id", equipment.JobID))
		return errors.New("equipment does not belong to this job")
	}

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

	if err := uc.repo.Update(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to update job equipment",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", equipment.ID))
		return err
	}

	slog.InfoContext(ctx, "Job equipment updated successfully",
		slog.Int64("equipment_id", equipment.ID))

	return nil
}
