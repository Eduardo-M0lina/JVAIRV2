package property_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*PropertyEquipment, error) {
	equipment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property equipment by ID",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Property equipment retrieved successfully",
		slog.Int64("equipment_id", id))

	return equipment, nil
}
