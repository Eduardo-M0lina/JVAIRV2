package warranty_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, wt *WarrantyType) error {
	if err := wt.Validate(); err != nil {
		return err
	}

	// Validar que el type existe
	_, err := uc.repo.GetByID(ctx, wt.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty type for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", wt.ID))
		return err
	}

	if err := uc.repo.Update(ctx, wt); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty type",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", wt.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty type updated successfully",
		slog.Int64("warranty_type_id", wt.ID),
		slog.String("label", wt.Label))

	return nil
}
