package sms_template

import (
	"fmt"
	"strings"
	"time"
)

// SMSTemplate representa una plantilla de SMS.
type SMSTemplate struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Message   string     `json:"message"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (s *SMSTemplate) Validate() error {
	if strings.TrimSpace(s.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if strings.TrimSpace(s.Message) == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}
