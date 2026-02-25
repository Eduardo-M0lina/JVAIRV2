package invoice_payment

import (
	"context"
	"log/slog"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// Update actualiza un pago existente
func (r *Repository) Update(ctx context.Context, payment *domainPayment.InvoicePayment) error {
	query := `
		UPDATE invoice_payments SET
			payment_processor = ?, payment_id = ?, amount = ?, notes = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		payment.PaymentProcessor, payment.PaymentID, payment.Amount, payment.Notes,
		payment.ID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update invoice payment",
			slog.Int64("id", payment.ID),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
