package property_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64, propertyID int64) error {
	// Validar que el equipo existe
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property equipment for deletion",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return err
	}

	// Validar que el equipo pertenece a la propiedad
	if existing.PropertyID != propertyID {
		slog.WarnContext(ctx, "Equipment does not belong to property",
			slog.Int64("equipment_id", id),
			slog.Int64("property_id", propertyID))
		return errors.New("equipment does not belong to this property")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete property equipment",
			slog.String("error", err.Error()),
			slog.Int64("equipment_id", id))
		return err
	}

	slog.InfoContext(ctx, "Property equipment deleted successfully",
		slog.Int64("equipment_id", id))

	return nil
}
