package warranty_type

import (
	"fmt"
	"strings"
	"time"
)

// WarrantyType representa la entidad de dominio para un tipo de garantía
type WarrantyType struct {
	ID          int64      `json:"id"`
	Label       string     `json:"label"`
	LabelPlural string     `json:"labelPlural"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del tipo de garantía
func (wt *WarrantyType) Validate() error {
	if strings.TrimSpace(wt.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if strings.TrimSpace(wt.LabelPlural) == "" {
		return fmt.Errorf("label_plural is required")
	}
	return nil
}
