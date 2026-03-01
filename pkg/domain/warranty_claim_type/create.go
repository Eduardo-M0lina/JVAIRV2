package warranty_claim_type

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, wct *WarrantyClaimType) error {
	if err := wct.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, wct); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty claim type",
			slog.String("error", err.Error()),
			slog.String("label", wct.Label))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim type created successfully",
		slog.Int64("warranty_claim_type_id", wct.ID),
		slog.String("label", wct.Label))

	return nil
}
