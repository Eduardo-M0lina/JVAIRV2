package quote_status

import "errors"

var (
	// ErrQuoteStatusNotFound indica que el estado de cotización no fue encontrado
	ErrQuoteStatusNotFound = errors.New("quote status not found")

	// ErrQuoteStatusInUse indica que el estado de cotización tiene cotizaciones asociadas
	ErrQuoteStatusInUse = errors.New("quote status is in use by quotes")
)
