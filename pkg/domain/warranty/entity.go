package warranty

import (
	"fmt"
	"time"
)

// Warranty representa una garantía en el sistema
type Warranty struct {
	ID               int64      `json:"id"`
	WarrantyNumber   string     `json:"warrantyNumber"`
	JobID            int64      `json:"jobId"`
	WarrantyTypeID   int64      `json:"warrantyTypeId"`
	WarrantyStatusID int64      `json:"warrantyStatusId"`
	DateSubmitted    *time.Time `json:"dateSubmitted,omitempty"`
	AgreementNumber  *string    `json:"agreementNumber,omitempty"`
	AuditDone        bool       `json:"auditDone"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear una garantía
func (w *Warranty) ValidateCreate() error {
	if w.WarrantyNumber == "" {
		return fmt.Errorf("warranty_number is required")
	}
	if w.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if w.WarrantyTypeID == 0 {
		return fmt.Errorf("warranty_type_id is required")
	}
	if w.WarrantyStatusID == 0 {
		return fmt.Errorf("warranty_status_id is required")
	}
	return nil
}

// ValidateUpdate valida los campos para actualizar una garantía
func (w *Warranty) ValidateUpdate() error {
	if w.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if w.WarrantyNumber == "" {
		return fmt.Errorf("warranty_number is required")
	}
	if w.WarrantyTypeID == 0 {
		return fmt.Errorf("warranty_type_id is required")
	}
	if w.WarrantyStatusID == 0 {
		return fmt.Errorf("warranty_status_id is required")
	}
	return nil
}

// IsDeleted verifica si la garantía está eliminada
func (w *Warranty) IsDeleted() bool {
	return w.DeletedAt != nil
}
