package quote

import (
	"fmt"
	"strings"
	"time"
)

// Quote representa la entidad de dominio para una cotización
type Quote struct {
	ID            int64      `json:"id"`
	JobID         int64      `json:"jobId"`
	QuoteNumber   string     `json:"quoteNumber"`
	QuoteStatusID int64      `json:"quoteStatusId"`
	Amount        float64    `json:"amount"`
	Description   *string    `json:"description,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear una cotización
func (q *Quote) ValidateCreate() error {
	if strings.TrimSpace(q.QuoteNumber) == "" {
		return fmt.Errorf("quote_number is required")
	}
	if q.JobID <= 0 {
		return fmt.Errorf("job_id is required")
	}
	if q.QuoteStatusID <= 0 {
		return fmt.Errorf("quote_status_id is required")
	}
	return nil
}
