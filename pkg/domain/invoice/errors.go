package invoice

import "errors"

var (
	// ErrInvoiceNotFound indica que la factura no fue encontrada
	ErrInvoiceNotFound = errors.New("invoice not found")

	// ErrInvoiceDeleted indica que la factura está eliminada
	ErrInvoiceDeleted = errors.New("invoice is deleted")

	// ErrInvalidJob indica que el job no es válido
	ErrInvalidJob = errors.New("invalid job")
)
