package quote_status

import (
	"fmt"
	"strings"
	"time"
)

// QuoteStatus representa la entidad de dominio para un estado de cotización
type QuoteStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	Order     int        `json:"order"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de cotización
func (qs *QuoteStatus) Validate() error {
	if strings.TrimSpace(qs.Label) == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}
