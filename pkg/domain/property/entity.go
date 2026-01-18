package property

import (
	"fmt"
	"strings"
	"time"
)

// Property representa una propiedad de un cliente
type Property struct {
	ID           int64      `json:"id"`
	CustomerID   int64      `json:"customerId"`
	PropertyCode *string    `json:"propertyCode,omitempty"`
	Street       string     `json:"street"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Zip          string     `json:"zip"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

// GetAddress retorna la dirección completa formateada
func (p *Property) GetAddress() string {
	return fmt.Sprintf("%s, %s, %s %s", p.Street, p.City, p.State, p.Zip)
}

// GetName retorna el nombre de la propiedad
// Si tiene property_code: "Property {property_code}"
// Si no: "{street}"
func (p *Property) GetName() string {
	if p.PropertyCode != nil && *p.PropertyCode != "" {
		return fmt.Sprintf("Property %s", *p.PropertyCode)
	}
	return p.Street
}

// IsDeleted verifica si la propiedad está eliminada
func (p *Property) IsDeleted() bool {
	return p.DeletedAt != nil
}

// Validate valida los campos requeridos de la propiedad
func (p *Property) Validate() error {
	if p.CustomerID == 0 {
		return fmt.Errorf("customer_id is required")
	}

	if strings.TrimSpace(p.Street) == "" {
		return fmt.Errorf("street is required")
	}

	if strings.TrimSpace(p.City) == "" {
		return fmt.Errorf("city is required")
	}

	if strings.TrimSpace(p.State) == "" {
		return fmt.Errorf("state is required")
	}

	if strings.TrimSpace(p.Zip) == "" {
		return fmt.Errorf("zip is required")
	}

	return nil
}
