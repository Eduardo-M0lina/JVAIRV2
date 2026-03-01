package warranty_claim_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, wct *WarrantyClaimType) error {
	if err := wct.Validate(); err != nil {
		return err
	}

	_, err := uc.repo.GetByID(ctx, wct.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim type for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_type_id", wct.ID))
		return err
	}

	if err := uc.repo.Update(ctx, wct); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim type",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_type_id", wct.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim type updated successfully",
		slog.Int64("warranty_claim_type_id", wct.ID),
		slog.String("label", wct.Label))

	return nil
}
