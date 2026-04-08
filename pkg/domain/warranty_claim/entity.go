package warranty_claim

import (
	"fmt"
	"time"
)

// WarrantyClaim representa un reclamo de garantía en el sistema
type WarrantyClaim struct {
	ID                    int64      `json:"id"`
	InternalClaimNumber   string     `json:"internalClaimNumber"`
	WarrantyClaimTypeID   int64      `json:"warrantyClaimTypeId"`
	WarrantyClaimStatusID int64      `json:"warrantyClaimStatusId"`
	JobID                 int64      `json:"jobId"`
	InvoiceNumber         *string    `json:"invoiceNumber,omitempty"`
	WorkDone              bool       `json:"workDone"`
	WarrantyPart          *string    `json:"warrantyPart,omitempty"`
	Manufacturer          *string    `json:"manufacturer,omitempty"`
	ModelNumber           *string    `json:"modelNumber,omitempty"`
	PartNumber            *string    `json:"partNumber,omitempty"`
	ReplacementPartNumber *string    `json:"replacementPartNumber,omitempty"`
	PartDistributor       *string    `json:"partDistributor,omitempty"`
	PartInvoiceNumber     *string    `json:"partInvoiceNumber,omitempty"`
	OldPartSerialNumber   *string    `json:"oldPartSerialNumber,omitempty"`
	NewPartSerialNumber   *string    `json:"newPartSerialNumber,omitempty"`
	EsaNumber             *string    `json:"esaNumber,omitempty"`
	Serial                *string    `json:"serial,omitempty"`
	ClaimNumber           *string    `json:"claimNumber,omitempty"`
	Approved              bool       `json:"approved"`
	PartsCreditReceived   bool       `json:"partsCreditReceived"`
	LaborPaymentReceived  bool       `json:"laborPaymentReceived"`
	Notes                 *string    `json:"notes,omitempty"`
	CreatedAt             *time.Time `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear un reclamo
func (wc *WarrantyClaim) ValidateCreate() error {
	if wc.InternalClaimNumber == "" {
		return fmt.Errorf("internal_claim_number is required")
	}
	if wc.WarrantyClaimTypeID == 0 {
		return fmt.Errorf("warranty_claim_type_id is required")
	}
	if wc.WarrantyClaimStatusID == 0 {
		return fmt.Errorf("warranty_claim_status_id is required")
	}
	if wc.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	return nil
}

// ValidateUpdate valida los campos para actualizar un reclamo
func (wc *WarrantyClaim) ValidateUpdate() error {
	if wc.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if wc.InternalClaimNumber == "" {
		return fmt.Errorf("internal_claim_number is required")
	}
	if wc.WarrantyClaimTypeID == 0 {
		return fmt.Errorf("warranty_claim_type_id is required")
	}
	if wc.WarrantyClaimStatusID == 0 {
		return fmt.Errorf("warranty_claim_status_id is required")
	}
	return nil
}

// IsDeleted verifica si el reclamo está eliminado
func (wc *WarrantyClaim) IsDeleted() bool {
	return wc.DeletedAt != nil
}
