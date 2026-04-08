package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*WarrantyEquipment, error) {
	equipment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty equipment by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_equipment_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty equipment retrieved successfully",
		slog.Int64("warranty_equipment_id", id))

	return equipment, nil
}
