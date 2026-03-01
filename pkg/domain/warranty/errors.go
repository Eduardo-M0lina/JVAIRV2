package warranty

import "errors"

var (
	// ErrWarrantyNotFound indica que la garantía no fue encontrada
	ErrWarrantyNotFound = errors.New("warranty not found")

	// ErrWarrantyDeleted indica que la garantía está eliminada
	ErrWarrantyDeleted = errors.New("warranty is deleted")

	// ErrInvalidJob indica que el job no es válido
	ErrInvalidJob = errors.New("invalid job")

	// ErrInvalidWarrantyType indica que el tipo de garantía no es válido
	ErrInvalidWarrantyType = errors.New("invalid warranty type")

	// ErrInvalidWarrantyStatus indica que el estado de garantía no es válido
	ErrInvalidWarrantyStatus = errors.New("invalid warranty status")
)
