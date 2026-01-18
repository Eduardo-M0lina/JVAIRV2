package customer

import (
	"fmt"
	"strings"
	"time"
)

type Customer struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Email                *string    `json:"email,omitempty"`
	Phone                *string    `json:"phone,omitempty"`
	Mobile               *string    `json:"mobile,omitempty"`
	Fax                  *string    `json:"fax,omitempty"`
	PhoneOther           *string    `json:"phoneOther,omitempty"`
	Website              *string    `json:"website,omitempty"`
	ContactName          *string    `json:"contactName,omitempty"`
	ContactEmail         *string    `json:"contactEmail,omitempty"`
	ContactPhone         *string    `json:"contactPhone,omitempty"`
	BillingAddressStreet *string    `json:"billingAddressStreet,omitempty"`
	BillingAddressCity   *string    `json:"billingAddressCity,omitempty"`
	BillingAddressState  *string    `json:"billingAddressState,omitempty"`
	BillingAddressZip    *string    `json:"billingAddressZip,omitempty"`
	WorkflowID           int64      `json:"workflowId"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            *time.Time `json:"createdAt,omitempty"`
	UpdatedAt            *time.Time `json:"updatedAt,omitempty"`
	DeletedAt            *time.Time `json:"deletedAt,omitempty"`
}

type BillingAddress struct {
	Street *string `json:"street,omitempty"`
	City   *string `json:"city,omitempty"`
	State  *string `json:"state,omitempty"`
	Zip    *string `json:"zip,omitempty"`
}

func (c *Customer) GetBillingAddress() *BillingAddress {
	if c.BillingAddressStreet == nil && c.BillingAddressCity == nil &&
		c.BillingAddressState == nil && c.BillingAddressZip == nil {
		return nil
	}

	return &BillingAddress{
		Street: c.BillingAddressStreet,
		City:   c.BillingAddressCity,
		State:  c.BillingAddressState,
		Zip:    c.BillingAddressZip,
	}
}

// Validate valida los campos requeridos del customer
func (c *Customer) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if c.WorkflowID == 0 {
		return fmt.Errorf("workflow_id is required")
	}

	return nil
}
