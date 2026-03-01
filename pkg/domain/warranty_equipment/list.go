package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) ListByWarrantyID(ctx context.Context, warrantyID int64) ([]*WarrantyEquipment, error) {
	// Verificar que la garantía existe
	if _, err := uc.warrantyCheck.GetByID(ctx, warrantyID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty",
			slog.Int64("warrantyId", warrantyID),
			slog.String("error", err.Error()))
		return nil, ErrInvalidWarranty
	}

	equipment, err := uc.repo.ListByWarrantyID(ctx, warrantyID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list warranty equipment",
			slog.String("error", err.Error()),
			slog.Int64("warrantyId", warrantyID))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty equipment listed successfully",
		slog.Int64("warrantyId", warrantyID),
		slog.Int("count", len(equipment)))

	return equipment, nil
}
