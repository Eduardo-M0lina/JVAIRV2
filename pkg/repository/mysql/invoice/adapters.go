package invoice

import (
	"context"
	"database/sql"

	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	domainInvoicePayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// JobCheckerAdapter adapta la verificación de jobs para el use case de invoices
type JobCheckerAdapter struct {
	db *sql.DB
}

func NewJobCheckerAdapter(db *sql.DB) domainInvoice.JobChecker {
	return &JobCheckerAdapter{db: db}
}

func (a *JobCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainInvoice.ErrInvalidJob
	}
	return true, nil
}

// InvoiceCheckerAdapter adapta la verificación de facturas para el use case de invoice payments
type InvoiceCheckerAdapter struct {
	db *sql.DB
}

func NewInvoiceCheckerAdapter(db *sql.DB) domainInvoicePayment.InvoiceChecker {
	return &InvoiceCheckerAdapter{db: db}
}

func (a *InvoiceCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM invoices WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainInvoicePayment.ErrInvalidInvoice
	}
	return true, nil
}
