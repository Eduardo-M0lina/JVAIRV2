package warranty

import (
	"context"
	"log/slog"
)

// Create crea una nueva garantía
func (uc *UseCase) Create(ctx context.Context, w *Warranty) error {
	if err := w.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que el job existe
	if _, err := uc.jobCheck.GetByID(ctx, w.JobID); err != nil {
		slog.ErrorContext(ctx, "Invalid job",
			slog.Int64("jobId", w.JobID),
			slog.String("error", err.Error()))
		return ErrInvalidJob
	}

	// Verificar que el tipo de garantía existe
	if _, err := uc.typeCheck.GetByID(ctx, w.WarrantyTypeID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty type",
			slog.Int64("warrantyTypeId", w.WarrantyTypeID),
			slog.String("error", err.Error()))
		return ErrInvalidWarrantyType
	}

	// Verificar que el estado de garantía existe
	if _, err := uc.statusCheck.GetByID(ctx, w.WarrantyStatusID); err != nil {
		slog.ErrorContext(ctx, "Invalid warranty status",
			slog.Int64("warrantyStatusId", w.WarrantyStatusID),
			slog.String("error", err.Error()))
		return ErrInvalidWarrantyStatus
	}

	if err := uc.repo.Create(ctx, w); err != nil {
		slog.ErrorContext(ctx, "Failed to create warranty",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Warranty created successfully",
		slog.Int64("id", w.ID))

	return nil
}
