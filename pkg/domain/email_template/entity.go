package email_template

import (
	"fmt"
	"strings"
	"time"
)

// EmailTemplate representa una plantilla de correo.
type EmailTemplate struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (e *EmailTemplate) Validate() error {
	if strings.TrimSpace(e.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if strings.TrimSpace(e.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if strings.TrimSpace(e.Body) == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}
