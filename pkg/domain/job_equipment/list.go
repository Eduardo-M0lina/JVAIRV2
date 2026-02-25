package job_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, jobID int64, equipmentType string) ([]*JobEquipment, error) {
	// Validar que el job existe
	exists, err := uc.jobChecker.JobExists(jobID)
	if err != nil {
		slog.WarnContext(ctx, "Invalid job_id for listing equipment",
			slog.Int64("job_id", jobID),
			slog.String("error", err.Error()))
		return nil, errors.New("invalid job_id")
	}
	if !exists {
		slog.WarnContext(ctx, "Job not found",
			slog.Int64("job_id", jobID))
		return nil, errors.New("invalid job_id")
	}

	// Validar type si se proporciona
	if equipmentType != "" && !isValidType(equipmentType) {
		return nil, errors.New("type must be one of: current, new")
	}

	equipment, err := uc.repo.List(ctx, jobID, equipmentType)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list job equipment",
			slog.String("error", err.Error()),
			slog.Int64("job_id", jobID))
		return nil, err
	}

	slog.InfoContext(ctx, "Job equipment listed successfully",
		slog.Int("total", len(equipment)),
		slog.Int64("job_id", jobID))

	return equipment, nil
}
