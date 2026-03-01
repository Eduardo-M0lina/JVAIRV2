package warranty

import (
	"context"
	"log/slog"
)

// Update actualiza una garantía existente
func (uc *UseCase) Update(ctx context.Context, w *Warranty) error {
	if err := w.ValidateUpdate(); err != nil {
		return err
	}

	// Verificar que la garantía existe
	existing, err := uc.repo.GetByID(ctx, w.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_id", w.ID))
		return err
	}

	if existing.IsDeleted() {
		return ErrWarrantyDeleted
	}

	// Verificar que el tipo de garantía existe
	if _, err := uc.typeCheck.GetByID(ctx, w.WarrantyTypeID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty type",
			slog.Int64("warrantyTypeId", w.WarrantyTypeID),
			slog.String("error", err.Error()))
		return ErrInvalidWarrantyType
	}

	// Verificar que el estado de garantía existe
	if _, err := uc.statusCheck.GetByID(ctx, w.WarrantyStatusID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty status",
			slog.Int64("warrantyStatusId", w.WarrantyStatusID),
			slog.String("error", err.Error()))
		return ErrInvalidWarrantyStatus
	}

	if err := uc.repo.Update(ctx, w); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty",
			slog.String("error", err.Error()),
			slog.Int64("warranty_id", w.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty updated successfully",
		slog.Int64("warranty_id", w.ID))

	return nil
}
