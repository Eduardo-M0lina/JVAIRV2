package alert

import (
	"fmt"
	"time"
)

type Alert struct {
	ID           int64      `json:"id"`
	UserID       *int64     `json:"userId,omitempty"`
	AlertType    string     `json:"alertType"`
	EntityID     int64      `json:"entityId"`
	EntityType   string     `json:"entityType"`
	MessageLevel string     `json:"messageLevel"`
	Message      string     `json:"message"`
	IsRead       bool       `json:"isRead"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

func (a *Alert) ValidateCreate() error {
	if a.AlertType == "" {
		return fmt.Errorf("alert_type is required")
	}
	if a.EntityID == 0 {
		return fmt.Errorf("entity_id is required")
	}
	if a.EntityType == "" {
		return fmt.Errorf("entity_type is required")
	}
	if a.MessageLevel == "" {
		return fmt.Errorf("message_level is required")
	}
	if a.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}
