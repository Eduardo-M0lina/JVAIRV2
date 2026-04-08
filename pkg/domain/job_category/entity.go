package job_category

import (
	"fmt"
	"strings"
	"time"
)

// JobCategory representa la entidad de dominio para una categoría de trabajo
type JobCategory struct {
	ID          int64      `json:"id"`
	Label       string     `json:"label"`
	LabelPlural string     `json:"labelPlural"`
	Type        string     `json:"type"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos de la categoría de trabajo
func (jc *JobCategory) Validate() error {
	if strings.TrimSpace(jc.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if strings.TrimSpace(jc.LabelPlural) == "" {
		return fmt.Errorf("label_plural is required")
	}
	if strings.TrimSpace(jc.Type) == "" {
		return fmt.Errorf("type is required")
	}
	return nil
}
