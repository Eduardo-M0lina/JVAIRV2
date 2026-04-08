package warranty_claim_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) GetByID(ctx context.Context, id int64) (*WarrantyClaimStatus, error) {
	wcs, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim status by ID",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", id))
		return nil, err
	}

	slog.InfoContext(ctx, "Warranty claim status retrieved successfully",
		slog.Int64("warranty_claim_status_id", id))

	return wcs, nil
}
