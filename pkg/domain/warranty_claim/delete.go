package warranty_claim

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get warranty claim for deletion",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_id", id))
		return err
	}

	if existing.IsDeleted() {
		return ErrWarrantyClaimDeleted
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete warranty claim",
			slog.String("error", err.Error()),
			slog.Int64("warranty_claim_id", id))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim deleted successfully",
		slog.Int64("warranty_claim_id", id))

	return nil
}
