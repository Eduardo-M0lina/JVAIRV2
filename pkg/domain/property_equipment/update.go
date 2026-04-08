package property_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, equipment *PropertyEquipment) error {
	// Validar que el equipo existe
	existing, err := uc.repo.GetByID(ctx, equipment.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property equipment for update",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", equipment.ID))
		return err
	}

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

	// Validar que el equipo pertenece a la propiedad
	if existing.PropertyID != equipment.PropertyID {
		slog.WarnContext(ctx, "Equipment does not belong to property",
			slog.Int64("equipment_id", equipment.ID),
			slog.Int64("property_id", equipment.PropertyID))
		return errors.New("equipment does not belong to this property")
	}

	if err := uc.repo.Update(ctx, equipment); err != nil {
		slog.ErrorContext(ctx, "Failed to update property equipment",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", equipment.ID))
		return err
	}

	slog.InfoContext(ctx, "Property equipment updated successfully",
		slog.Int64("equipment_id", equipment.ID))

	return nil
}
