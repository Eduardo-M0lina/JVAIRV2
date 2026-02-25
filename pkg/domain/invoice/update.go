package invoice

import (
	"context"
	"log/slog"
)

// Update actualiza una factura existente
func (uc *UseCase) Update(ctx context.Context, inv *Invoice) error {
	if err := inv.ValidateUpdate(); err != nil {
		return err
	}

	// Verificar que la factura existe
	existing, err := uc.repo.GetByID(ctx, inv.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Invoice not found for update",
			slog.Int64("id", inv.ID),
			slog.String("error", err.Error()))
		return ErrInvoiceNotFound
	}

	if existing.IsDeleted() {
		return ErrInvoiceNotFound
	}

	// Verificar job si cambió
	if inv.JobID > 0 && inv.JobID != existing.JobID {
		if _, err := uc.jobCheck.GetByID(ctx, inv.JobID); err != nil {
			return ErrInvalidJob
		}
	}

	if err := uc.repo.Update(ctx, inv); err != nil {
		slog.ErrorContext(ctx, "Failed to update invoice",
			slog.Int64("id", inv.ID),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice updated successfully",
		slog.Int64("id", inv.ID))

	return nil
}
