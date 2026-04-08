package property_equipment

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) List(ctx context.Context, propertyID int64) ([]*PropertyEquipment, error) {
	// Validar que la propiedad existe y no está eliminada
	prop, err := uc.propertyRepo.GetByID(ctx, propertyID)
	if err != nil {
		slog.WarnContext(ctx, "Invalid property_id for listing equipment",
			slog.Int64("property_id", propertyID),
			slog.String("error", err.Error()))
		return nil, errors.New("invalid property_id")
	}
	if prop.DeletedAt != nil {
		slog.WarnContext(ctx, "Property is deleted",
			slog.Int64("property_id", propertyID))
		return nil, errors.New("property is deleted")
	}

	equipment, err := uc.repo.List(ctx, propertyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list property equipment",
			slog.String("error", err.Error()),
			slog.Int64("property_id", propertyID))
		return nil, err
	}

	slog.InfoContext(ctx, "Property equipment listed successfully",
		slog.Int("total", len(equipment)),
		slog.Int64("property_id", propertyID))

	return equipment, nil
}
