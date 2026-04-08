package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty equipment for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_equipment_id", id))
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty equipment",
			slog.String("error", err.Error()),
			slog.Int64("warranty_equipment_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty equipment deleted successfully",
		slog.Int64("warranty_equipment_id", id))

	return nil
}
