package warranty_status

import "errors"

var (
	// ErrWarrantyStatusNotFound indica que el estado de garantía no fue encontrado
	ErrWarrantyStatusNotFound = errors.New("warranty status not found")

	// ErrWarrantyStatusInUse indica que el estado de garantía tiene garantías asociadas
	ErrWarrantyStatusInUse = errors.New("warranty status is in use by warranties")
)
