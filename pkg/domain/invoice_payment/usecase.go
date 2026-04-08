package invoice_payment

import "context"

// Service define la interfaz del servicio de invoice payments
type Service interface {
	Create(ctx context.Context, payment *InvoicePayment) error
	GetByID(ctx context.Context, invoiceID, id int64) (*InvoicePayment, error)
	ListByInvoiceID(ctx context.Context, invoiceID int64, filters map[string]interface{}, page, pageSize int) ([]*InvoicePayment, int, error)
	Update(ctx context.Context, payment *InvoicePayment) error
	Delete(ctx context.Context, invoiceID, id int64) error
}

// InvoiceChecker verifica existencia de facturas
type InvoiceChecker interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
}

// UseCase implementa la lógica de negocio de invoice payments
type UseCase struct {
	repo         Repository
	invoiceCheck InvoiceChecker
}

// NewUseCase crea una nueva instancia del caso de uso de invoice payments
func NewUseCase(repo Repository, invoiceCheck InvoiceChecker) *UseCase {
	return &UseCase{
		repo:         repo,
		invoiceCheck: invoiceCheck,
	}
}
