package warranty_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty status for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", id))
		return err
	}

	hasWarranties, err := uc.repo.HasWarranties(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check warranty status warranties",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", id))
		return err
	}

	if hasWarranties {
		slog.WarnContext(ctx, "Cannot delete warranty status with warranties",
			slog.Int64("warranty_status_id", id))
		return ErrWarrantyStatusInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty status",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty status deleted successfully",
		slog.Int64("warranty_status_id", id))

	return nil
}
