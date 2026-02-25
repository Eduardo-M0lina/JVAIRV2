package invoice

import (
	"context"
	"log/slog"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
)

// Create crea una nueva factura en la base de datos
func (r *Repository) Create(ctx context.Context, inv *domainInvoice.Invoice) error {
	query := `
		INSERT INTO invoices (
			job_id, invoice_number, total, description,
			allow_online_payments, notes,
			created_at, updated_at
		) VALUES (
			?, ?, ?, ?,
			?, ?,
			NOW(), NOW()
		)
	`

	result, err := r.db.ExecContext(ctx, query,
		inv.JobID, inv.InvoiceNumber, inv.Total, inv.Description,
		inv.AllowOnlinePayments, inv.Notes,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create invoice",
			slog.String("error", err.Error()))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last insert ID",
			slog.String("error", err.Error()))
		return err
	}

	inv.ID = id

	return nil
}
