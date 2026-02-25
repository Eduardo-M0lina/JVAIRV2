package invoice

import (
	"context"
	"database/sql"
	"log/slog"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
)

// GetByID obtiene una factura por su ID, incluyendo el balance calculado
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainInvoice.Invoice, error) {
	query := `
		SELECT
			i.id, i.job_id, i.invoice_number, i.total, i.description,
			i.allow_online_payments, i.notes,
			i.created_at, i.updated_at, i.deleted_at,
			IFNULL(i.total - SUM(ip.amount), i.total) as balance
		FROM invoices i
		LEFT JOIN invoice_payments ip ON ip.invoice_id = i.id AND ip.deleted_at IS NULL
		WHERE i.id = ?
		GROUP BY i.id
	`

	inv := &domainInvoice.Invoice{}
	var balance float64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.JobID, &inv.InvoiceNumber, &inv.Total, &inv.Description,
		&inv.AllowOnlinePayments, &inv.Notes,
		&inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		&balance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainInvoice.ErrInvoiceNotFound
		}
		slog.ErrorContext(ctx, "Failed to get invoice by ID",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	inv.Balance = &balance

	return inv, nil
}
