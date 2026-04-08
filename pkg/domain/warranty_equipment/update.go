package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, equipment *WarrantyEquipment) error {
	if err := equipment.Validate(); err != nil {
		return err
	}

	// Verificar que el equipo existe
	_, err := uc.repo.GetByID(ctx, equipment.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty equipment for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_equipment_id", equipment.ID))
		return err
	}

	if err := uc.repo.Update(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty equipment",
			slog.String("error", err.Error()),
			slog.Int64("warranty_equipment_id", equipment.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty equipment updated successfully",
		slog.Int64("warranty_equipment_id", equipment.ID))

	return nil
}
