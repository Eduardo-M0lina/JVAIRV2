package invoice

import (
	"context"
	"database/sql"
	"log/slog"

	domainInvoice "github.com/angumol/jvairv2/pkg/domain/invoice"
)

// GetByInvoiceNumber obtiene una factura por su número, incluyendo el balance calculado
func (r *Repository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domainInvoice.Invoice, error) {
	query := `
		SELECT
			i.id, i.job_id, i.invoice_number, i.total, i.description,
			i.allow_online_payments, i.notes,
			i.created_at, i.updated_at, i.deleted_at,
			IFNULL(i.total - SUM(ip.amount), i.total) as balance
		FROM invoices i
		LEFT JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.deleted_at IS NULL
		WHERE i.invoice_number = ? AND i.deleted_at IS NULL
		GROUP BY i.id
	`

	inv := &domainInvoice.Invoice{}
	var balance float64
	err := r.db.QueryRowContext(ctx, query, invoiceNumber).Scan(
		&inv.ID, &inv.JobID, &inv.InvoiceNumber, &inv.Total, &inv.Description,
		&inv.AllowOnlinePayments, &inv.Notes,
		&inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		&balance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainInvoice.ErrInvoiceNotFound
		}
		slog.ErrorContext(ctx, "Failed to get invoice by invoice_number",
			slog.String("invoiceNumber", invoiceNumber),
			slog.String("error", err.Error()))
		return nil, err
	}

	inv.Balance = &balance

	return inv, nil
}
