package supervisor

import (
	"fmt"
	"strings"
	"time"
)

type Supervisor struct {
	ID         int64      `json:"id"`
	CustomerID int64      `json:"customerId"`
	Name       string     `json:"name"`
	Phone      *string    `json:"phone,omitempty"`
	Email      *string    `json:"email,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

// Validate valida los campos requeridos del supervisor
func (s *Supervisor) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if s.CustomerID == 0 {
		return fmt.Errorf("customer_id is required")
	}

	return nil
}
