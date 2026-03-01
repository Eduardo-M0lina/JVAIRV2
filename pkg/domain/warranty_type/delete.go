package warranty_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que el type existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty type for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", id))
		return err
	}

	// Verificar que no tenga warranties asociadas
	hasWarranties, err := uc.repo.HasWarranties(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check warranty type warranties",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", id))
		return err
	}

	if hasWarranties {
		slog.WarnContext(ctx, "Cannot delete warranty type with warranties",
			slog.Int64("warranty_type_id", id))
		return ErrWarrantyTypeInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty type",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty type deleted successfully",
		slog.Int64("warranty_type_id", id))

	return nil
}
