package warranty_claim

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, wc *WarrantyClaim) error {
	if err := wc.ValidateUpdate(); err != nil {
		return err
	}

	existing, err := uc.repo.GetByID(ctx, wc.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim for update",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_id", wc.ID))
		return err
	}

	if existing.IsDeleted() {
		return ErrWarrantyClaimDeleted
	}

	if _, err := uc.typeCheck.GetByID(ctx, wc.WarrantyClaimTypeID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty claim type",
			slog.Int64("warrantyClaimTypeId", wc.WarrantyClaimTypeID),
			slog.String("error", err.Error()))
		return ErrInvalidClaimType
	}

	if _, err := uc.statusCheck.GetByID(ctx, wc.WarrantyClaimStatusID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty claim status",
			slog.Int64("warrantyClaimStatusId", wc.WarrantyClaimStatusID),
			slog.String("error", err.Error()))
		return ErrInvalidClaimStatus
	}

	if err := uc.repo.Update(ctx, wc); err != nil {
		slog.ErrorContext(ctx, "Failed to update warranty claim",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_id", wc.ID))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim updated successfully",
		slog.Int64("warranty_claim_id", wc.ID))

	return nil
}
