package property_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, equipment *PropertyEquipment) error {
	// Validar que la propiedad existe y no está eliminada
	prop, err := uc.propertyRepo.GetByID(ctx, equipment.PropertyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate property",
			slog.String("error", err.Error()),
			slog.Int64("property_id", equipment.PropertyID))
		return errors.New("invalid property_id")
	}

	if prop.DeletedAt != nil {
		slog.WarnContext(ctx, "Property is deleted",
			slog.Int64("property_id", equipment.PropertyID))
		return errors.New("property is deleted")
	}

	if err := uc.repo.Create(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to create property equipment",
			slog.String("error", err.Error()),
			slog.Int64("property_id", equipment.PropertyID))
		return err
	}

	slog.InfoContext(ctx, "Property equipment created successfully",
		slog.Int64("equipment_id", equipment.ID),
		slog.Int64("property_id", equipment.PropertyID))

	return nil
}
