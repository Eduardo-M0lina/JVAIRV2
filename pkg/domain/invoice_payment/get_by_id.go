package invoice_payment

import (
	"context"
	"log/slog"
)

// GetByID obtiene un pago por su ID, verificando que pertenece a la factura indicada
func (uc *UseCase) GetByID(ctx context.Context, invoiceID, id int64) (*InvoicePayment, error) {
	payment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get invoice payment",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	if payment.IsDeleted() {
		return nil, ErrPaymentNotFound
	}

	// Verificar que el pago pertenece a la factura indicada
	if payment.InvoiceID != invoiceID {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}
