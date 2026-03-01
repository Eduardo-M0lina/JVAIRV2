package warranty_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, ws *WarrantyStatus) error {
	if err := ws.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, ws); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty status",
			slog.String("error", err.Error()),
			slog.String("label", ws.Label))
		return err
	}

	slog.InfoContext(ctx, "Warranty status created successfully",
		slog.Int64("warranty_status_id", ws.ID),
		slog.String("label", ws.Label))

	return nil
}
