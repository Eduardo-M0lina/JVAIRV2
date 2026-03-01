package warranty_claim_status

import "errors"

var (
	// ErrWarrantyClaimStatusNotFound indica que el estado de reclamo no fue encontrado
	ErrWarrantyClaimStatusNotFound = errors.New("warranty claim status not found")

	// ErrWarrantyClaimStatusInUse indica que el estado de reclamo tiene reclamos asociados
	ErrWarrantyClaimStatusInUse = errors.New("warranty claim status is in use by warranty claims")
)
