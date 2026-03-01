package warranty_claim

import "errors"

var (
	// ErrWarrantyClaimNotFound indica que el reclamo no fue encontrado
	ErrWarrantyClaimNotFound = errors.New("warranty claim not found")

	// ErrWarrantyClaimDeleted indica que el reclamo está eliminado
	ErrWarrantyClaimDeleted = errors.New("warranty claim is deleted")

	// ErrInvalidJob indica que el job no es válido
	ErrInvalidJob = errors.New("invalid job")

	// ErrInvalidClaimType indica que el tipo de reclamo no es válido
	ErrInvalidClaimType = errors.New("invalid warranty claim type")

	// ErrInvalidClaimStatus indica que el estado de reclamo no es válido
	ErrInvalidClaimStatus = errors.New("invalid warranty claim status")
)
