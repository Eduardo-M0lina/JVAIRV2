package warranty_claim_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, wcs *WarrantyClaimStatus) error {
	if err := wcs.Validate(); err != nil {
		return err
	}

	_, err := uc.repo.GetByID(ctx, wcs.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim status for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", wcs.ID))
		return err
	}

	if err := uc.repo.Update(ctx, wcs); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim status",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", wcs.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim status updated successfully",
		slog.Int64("warranty_claim_status_id", wcs.ID),
		slog.String("label", wcs.Label))

	return nil
}
