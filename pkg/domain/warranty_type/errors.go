package warranty_type

import "errors"

var (
	// ErrWarrantyTypeNotFound indica que el tipo de garantía no fue encontrado
	ErrWarrantyTypeNotFound = errors.New("warranty type not found")

	// ErrWarrantyTypeInUse indica que el tipo de garantía tiene garantías asociadas
	ErrWarrantyTypeInUse = errors.New("warranty type is in use by warranties")
)
