package warranty_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*WarrantyType, error) {
	wt, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty type by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_type_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty type retrieved successfully",
		slog.Int64("warranty_type_id", id))

	return wt, nil
}
