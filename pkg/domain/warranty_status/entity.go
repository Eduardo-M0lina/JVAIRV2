package warranty_status

import (
	"fmt"
	"strings"
	"time"
)

// WarrantyStatus representa la entidad de dominio para un estado de garantía
type WarrantyStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	Order     int        `json:"order"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de garantía
func (ws *WarrantyStatus) Validate() error {
	if strings.TrimSpace(ws.Label) == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}
