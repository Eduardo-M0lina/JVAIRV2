package invoice_payment

import (
	"context"
	"log/slog"
)

// Create crea un nuevo pago de factura
func (uc *UseCase) Create(ctx context.Context, payment *InvoicePayment) error {
	if err := payment.ValidateCreate(); err != nil {
		return err
	}

	// Verificar que la factura existe
	if _, err := uc.invoiceCheck.GetByID(ctx, payment.InvoiceID); err != nil {
		slog.ErrorContext(ctx, "Invalid invoice",
			slog.Int64("invoiceId", payment.InvoiceID),
			slog.String("error", err.Error()))
		return ErrInvalidInvoice
	}

	if err := uc.repo.Create(ctx, payment); err != nil {
		slog.ErrorContext(ctx, "Failed to create invoice payment",
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Invoice payment created successfully",
		slog.Int64("id", payment.ID))

	return nil
}
