package invoice_payment

import (
	"context"
	"log/slog"
)

// ListByInvoiceID obtiene una lista paginada de pagos de una factura
func (uc *UseCase) ListByInvoiceID(ctx context.Context, invoiceID int64, filters map[string]interface{}, page, pageSize int) ([]*InvoicePayment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Verificar que la factura existe
	if _, err := uc.invoiceCheck.GetByID(ctx, invoiceID); err != nil {
		slog.ErrorContext(ctx, "Invalid invoice for listing payments",
			slog.Int64("invoiceId", invoiceID),
			slog.String("error", err.Error()))
		return nil, 0, ErrInvalidInvoice
	}

	payments, total, err := uc.repo.ListByInvoiceID(ctx, invoiceID, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list invoice payments",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Invoice payments listed successfully",
		slog.Int64("invoiceId", invoiceID),
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return payments, total, nil
}
