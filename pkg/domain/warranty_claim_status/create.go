package warranty_claim_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, wcs *WarrantyClaimStatus) error {
	if err := wcs.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, wcs); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty claim status",
			slog.String("error", err.Error()),
			slog.String("label", wcs.Label))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim status created successfully",
		slog.Int64("warranty_claim_status_id", wcs.ID),
		slog.String("label", wcs.Label))

	return nil
}
