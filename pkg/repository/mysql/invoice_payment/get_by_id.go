package invoice_payment

import (
	"context"
	"database/sql"
	"log/slog"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// GetByID obtiene un pago por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainPayment.InvoicePayment, error) {
	query := `
		SELECT
			id, invoice_id, payment_processor, payment_id, amount, notes,
			created_at, updated_at, deleted_at
		FROM invoice_payments
		WHERE id = ?
	`

	payment := &domainPayment.InvoicePayment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID, &payment.InvoiceID, &payment.PaymentProcessor, &payment.PaymentID, &payment.Amount, &payment.Notes,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainPayment.ErrPaymentNotFound
		}
		slog.ErrorContext(ctx, "Failed to get invoice payment by ID",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	return payment, nil
}
