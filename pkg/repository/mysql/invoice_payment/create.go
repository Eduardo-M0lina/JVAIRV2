package invoice_payment

import (
	"context"
	"log/slog"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// Create crea un nuevo pago de factura en la base de datos
func (r *Repository) Create(ctx context.Context, payment *domainPayment.InvoicePayment) error {
	query := `
		INSERT INTO invoice_payments (
			invoice_id, payment_processor, payment_id, amount, notes,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.InvoiceID, payment.PaymentProcessor, payment.PaymentID, payment.Amount, payment.Notes,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create invoice payment",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	payment.ID = id

	return nil
}
