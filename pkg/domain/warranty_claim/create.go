package warranty_claim

import (
	"context"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, wc *WarrantyClaim) error {
	if err := wc.ValidateCreate(); err != nil {
		return err
	}

	if _, err := uc.jobCheck.GetByID(ctx, wc.JobID); err != nil {
		slog.ErrorContext(ctx, "Invalid job",
			slog.Int64("jobId", wc.JobID),
			slog.String("error", err.Error()))
		return ErrInvalidJob
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

	if err := uc.repo.Create(ctx, wc); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty claim",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Warranty claim created successfully",
		slog.Int64("id", wc.ID))

	return nil
}
