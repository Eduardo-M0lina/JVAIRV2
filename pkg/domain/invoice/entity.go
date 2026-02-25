package invoice

import (
	"fmt"
	"time"
)

// Invoice representa una factura en el sistema
type Invoice struct {
	ID                  int64      `json:"id"`
	JobID               int64      `json:"jobId"`
	InvoiceNumber       string     `json:"invoiceNumber"`
	Total               float64    `json:"total"`
	Description         *string    `json:"description,omitempty"`
	AllowOnlinePayments bool       `json:"allowOnlinePayments"`
	Notes               *string    `json:"notes,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty"`
	// Campo calculado: total - SUM(payments.amount)
	Balance *float64 `json:"balance,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear una factura
func (i *Invoice) ValidateCreate() error {
	if i.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if i.InvoiceNumber == "" {
		return fmt.Errorf("invoice_number is required")
	}
	return nil
}

// ValidateUpdate valida los campos para actualizar una factura
func (i *Invoice) ValidateUpdate() error {
	if i.ID == 0 {
		return fmt.Errorf("id is required")
	}
	return nil
}

// IsDeleted verifica si la factura está eliminada
func (i *Invoice) IsDeleted() bool {
	return i.DeletedAt != nil
}

// IsPaid verifica si la factura está pagada (balance == 0)
func (i *Invoice) IsPaid() bool {
	if i.Balance == nil {
		return false
	}
	return *i.Balance <= 0
}
