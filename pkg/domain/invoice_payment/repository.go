package invoice_payment

import "context"

// Repository define los métodos para interactuar con el almacenamiento de invoice payments
type Repository interface {
	// Create crea un nuevo pago de factura
	Create(ctx context.Context, payment *InvoicePayment) error

	// GetByID obtiene un pago por su ID
	GetByID(ctx context.Context, id int64) (*InvoicePayment, error)

	// ListByInvoiceID obtiene los pagos de una factura con paginación
	ListByInvoiceID(ctx context.Context, invoiceID int64, filters map[string]interface{}, page, pageSize int) ([]*InvoicePayment, int, error)

	// Update actualiza un pago existente
	Update(ctx context.Context, payment *InvoicePayment) error

	// Delete elimina un pago (soft delete)
	Delete(ctx context.Context, id int64) error
}
