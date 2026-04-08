package invoice_payment

import (
	"fmt"
	"time"
)

// InvoicePayment representa un pago asociado a una factura
type InvoicePayment struct {
	ID               int64      `json:"id"`
	InvoiceID        int64      `json:"invoiceId"`
	PaymentProcessor string     `json:"paymentProcessor"`
	PaymentID        string     `json:"paymentId"`
	Amount           float64    `json:"amount"`
	Notes            string     `json:"notes"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear un pago
func (p *InvoicePayment) ValidateCreate() error {
	if p.InvoiceID == 0 {
		return fmt.Errorf("invoice_id is required")
	}
	if p.PaymentProcessor == "" {
		return fmt.Errorf("payment_processor is required")
	}
	if p.PaymentID == "" {
		return fmt.Errorf("payment_id is required")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// ValidateUpdate valida los campos para actualizar un pago
func (p *InvoicePayment) ValidateUpdate() error {
	if p.ID == 0 {
		return fmt.Errorf("id is required")
	}
	return nil
}

// IsDeleted verifica si el pago está eliminado
func (p *InvoicePayment) IsDeleted() bool {
	return p.DeletedAt != nil
}
