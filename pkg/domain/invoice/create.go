package invoice

import (
	"context"
	"log/slog"
)

// Create crea una nueva factura
func (uc *UseCase) Create(ctx context.Context, inv *Invoice) error {
	if err := inv.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que el job existe
	if _, err := uc.jobCheck.GetByID(ctx, inv.JobID); err != nil {
		slog.ErrorContext(ctx, "Invalid job",
			slog.Int64("jobId", inv.JobID),
			slog.String("error", err.Error()))
		return ErrInvalidJob
	}

	if err := uc.repo.Create(ctx, inv); err != nil {
		slog.ErrorContext(ctx, "Failed to create invoice",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice created successfully",
		slog.Int64("id", inv.ID))

	return nil
}
