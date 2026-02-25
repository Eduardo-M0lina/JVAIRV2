package invoice_payment

import (
	"context"
	"log/slog"
)

// Update actualiza un pago existente
func (uc *UseCase) Update(ctx context.Context, payment *InvoicePayment) error {
	if err := payment.ValidateUpdate(); err != nil {
		return err
	}

	// Verificar que el pago existe
	existing, err := uc.repo.GetByID(ctx, payment.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Payment not found for update",
			slog.Int64("id", payment.ID),
			slog.String("error", err.Error()))
		return ErrPaymentNotFound
	}

	if existing.IsDeleted() {
		return ErrPaymentNotFound
	}

	// Mantener el invoice_id original
	payment.InvoiceID = existing.InvoiceID

	if err := uc.repo.Update(ctx, payment); err != nil {
		slog.ErrorContext(ctx, "Failed to update invoice payment",
			slog.Int64("id", payment.ID),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice payment updated successfully",
		slog.Int64("id", payment.ID))

	return nil
}
