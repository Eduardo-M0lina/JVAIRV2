package job_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*JobEquipment, error) {
	equipment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get job equipment by ID",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Job equipment retrieved successfully",
		slog.Int64("equipment_id", id))

	return equipment, nil
}
