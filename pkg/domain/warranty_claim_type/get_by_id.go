package warranty_claim_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*WarrantyClaimType, error) {
	wct, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim type by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_type_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty claim type retrieved successfully",
		slog.Int64("warranty_claim_type_id", id))

	return wct, nil
}
