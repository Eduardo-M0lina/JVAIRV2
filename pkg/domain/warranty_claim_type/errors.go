package warranty_claim_type

import "errors"

var (
	// ErrWarrantyClaimTypeNotFound indica que el tipo de reclamo no fue encontrado
	ErrWarrantyClaimTypeNotFound = errors.New("warranty claim type not found")

	// ErrWarrantyClaimTypeInUse indica que el tipo de reclamo tiene reclamos asociados
	ErrWarrantyClaimTypeInUse = errors.New("warranty claim type is in use by warranty claims")
)
