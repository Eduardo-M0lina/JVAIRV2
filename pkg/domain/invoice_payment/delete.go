package invoice_payment

import (
	"context"
	"log/slog"
)

// Delete elimina un pago (soft delete), verificando que pertenece a la factura indicada
func (uc *UseCase) Delete(ctx context.Context, invoiceID, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Payment not found for delete",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return ErrPaymentNotFound
	}

	if existing.IsDeleted() {
		return ErrPaymentNotFound
	}

	// Verificar que el pago pertenece a la factura indicada
	if existing.InvoiceID != invoiceID {
		return ErrPaymentNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete invoice payment",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice payment deleted successfully",
		slog.Int64("id", id))

	return nil
}
