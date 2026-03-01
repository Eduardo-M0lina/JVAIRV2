package warranty_claim_status

import (
	"fmt"
	"strings"
	"time"
)

// WarrantyClaimStatus representa la entidad de dominio para un estado de reclamo de garantía
type WarrantyClaimStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	Order     int        `json:"order"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de reclamo
func (wcs *WarrantyClaimStatus) Validate() error {
	if strings.TrimSpace(wcs.Label) == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}
