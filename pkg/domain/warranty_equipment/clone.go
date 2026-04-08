package warranty_equipment

import (
	"context"
	"log/slog"
)

func (uc *UseCase) CloneFromJobEquipment(ctx context.Context, warrantyID int64, jobID int64) error {
	// Verificar que la garantía existe
	if _, err := uc.warrantyCheck.GetByID(ctx, warrantyID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty for cloning",
			slog.Int64("warrantyId", warrantyID),
			slog.String("error", err.Error()))
		return ErrInvalidWarranty
	}

	if err := uc.repo.CloneFromJobEquipment(ctx, warrantyID, jobID); err != nil {
		slog.ErrorContext(ctx, "Failed to clone job equipment to warranty",
			slog.String("error", err.Error()),
			slog.Int64("warrantyId", warrantyID),
			slog.Int64("jobId", jobID))
		return err
	}

	slog.InfoContext(ctx, "Job equipment cloned to warranty successfully",
		slog.Int64("warrantyId", warrantyID),
		slog.Int64("jobId", jobID))

	return nil
}
