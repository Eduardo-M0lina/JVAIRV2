package job_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64, jobID int64) error {
	// Validar que el equipo existe
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job equipment for deletion",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return err
	}

	// Validar que el equipo pertenece al job
	if existing.JobID != jobID {
		slog.WarnContext(ctx, "Equipment does not belong to job",
			slog.Int64("equipment_id", id),
			slog.Int64("job_id", jobID))
		return errors.New("equipment does not belong to this job")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete job equipment",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return err
	}

	slog.InfoContext(ctx, "Job equipment deleted successfully",
		slog.Int64("equipment_id", id))

	return nil
}
