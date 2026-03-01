package warranty_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*WarrantyStatus, error) {
	ws, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty status by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_status_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty status retrieved successfully",
		slog.Int64("warranty_status_id", id))

	return ws, nil
}
