package invoice

import (
	"context"
	"log/slog"
)

// Delete elimina una factura (soft delete)
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Invoice not found for delete",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return ErrInvoiceNotFound
	}

	if existing.IsDeleted() {
		return ErrInvoiceNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete invoice",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice deleted successfully",
		slog.Int64("id", id))

	return nil
}
