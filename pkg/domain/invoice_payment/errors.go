package invoice_payment

import "errors"

var (
	// ErrPaymentNotFound indica que el pago no fue encontrado
	ErrPaymentNotFound = errors.New("invoice payment not found")

	// ErrPaymentDeleted indica que el pago está eliminado
	ErrPaymentDeleted = errors.New("invoice payment is deleted")

	// ErrInvalidInvoice indica que la factura no es válida
	ErrInvalidInvoice = errors.New("invalid invoice")
)
