package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, equipment *WarrantyEquipment) error {
	if err := equipment.Validate(); err != nil {
		return err
	}

	// Verificar que la garantía existe
	if _, err := uc.warrantyCheck.GetByID(ctx, equipment.WarrantyID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty",
			slog.Int64("warrantyId", equipment.WarrantyID),
			slog.String("error", err.Error()))
		return ErrInvalidWarranty
	}

	if err := uc.repo.Create(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty equipment",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Warranty equipment created successfully",
		slog.Int64("id", equipment.ID))

	return nil
}
