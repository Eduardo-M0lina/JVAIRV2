package warranty_claim_type

import (
	"fmt"
	"strings"
	"time"
)

// WarrantyClaimType representa la entidad de dominio para un tipo de reclamo de garantía
type WarrantyClaimType struct {
	ID          int64      `json:"id"`
	Label       string     `json:"label"`
	LabelPlural string     `json:"labelPlural"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del tipo de reclamo
func (wct *WarrantyClaimType) Validate() error {
	if strings.TrimSpace(wct.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if strings.TrimSpace(wct.LabelPlural) == "" {
		return fmt.Errorf("label_plural is required")
	}
	return nil
}
