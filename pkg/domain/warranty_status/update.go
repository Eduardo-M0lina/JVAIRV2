package warranty_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, ws *WarrantyStatus) error {
	if err := ws.Validate(); err != nil {
		return err
	}

	_, err := uc.repo.GetByID(ctx, ws.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty status for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", ws.ID))
		return err
	}

	if err := uc.repo.Update(ctx, ws); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty status",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", ws.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty status updated successfully",
		slog.Int64("warranty_status_id", ws.ID),
		slog.String("label", ws.Label))

	return nil
}
