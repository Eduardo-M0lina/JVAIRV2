package warranty_claim_status

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim status for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", id))
		return err
	}

	hasClaims, err := uc.repo.HasClaims(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check warranty claim status claims",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", id))
		return err
	}

	if hasClaims {
		slog.WarnContext(ctx, "Cannot delete warranty claim status with claims",
			slog.Int64("warranty_claim_status_id", id))
		return ErrWarrantyClaimStatusInUse
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty claim status",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_status_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim status deleted successfully",
		slog.Int64("warranty_claim_status_id", id))

	return nil
}
