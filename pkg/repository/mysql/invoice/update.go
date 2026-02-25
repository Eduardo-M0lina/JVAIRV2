package invoice

import (
	"context"
	"log/slog"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
)

// Update actualiza una factura existente
func (r *Repository) Update(ctx context.Context, inv *domainInvoice.Invoice) error {
	query := `
		UPDATE invoices SET
			job_id = ?, invoice_number = ?, total = ?, description = ?,
			allow_online_payments = ?, notes = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		inv.JobID, inv.InvoiceNumber, inv.Total, inv.Description,
		inv.AllowOnlinePayments, inv.Notes,
		inv.ID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update invoice",
			slog.Int64("id", inv.ID),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
