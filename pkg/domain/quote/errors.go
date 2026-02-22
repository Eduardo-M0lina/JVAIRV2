package quote

import "errors"

var (
	// ErrQuoteNotFound indica que la cotización no fue encontrada
	ErrQuoteNotFound = errors.New("quote not found")

	// ErrQuoteDeleted indica que la cotización está eliminada
	ErrQuoteDeleted = errors.New("quote is deleted")

	// ErrInvalidJob indica que el trabajo no es válido
	ErrInvalidJob = errors.New("invalid job")

	// ErrInvalidQuoteStatus indica que el estado de cotización no es válido
	ErrInvalidQuoteStatus = errors.New("invalid quote status")
)
