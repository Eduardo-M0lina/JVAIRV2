package technician_job_status

import (
	"fmt"
	"strings"
	"time"
)

// ValidClasses contiene los valores válidos para el campo class (colores Bootstrap)
var ValidClasses = []string{"blue", "indigo", "purple", "pink", "red", "orange", "yellow", "green", "teal", "cyan", "dark", "light"}

// TechnicianJobStatus representa la entidad de dominio para un estado de técnico de trabajo
type TechnicianJobStatus struct {
	ID          int64      `json:"id"`
	Label       string     `json:"label"`
	Class       *string    `json:"class,omitempty"`
	JobStatusID *int64     `json:"jobStatusId,omitempty"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de técnico de trabajo
func (tjs *TechnicianJobStatus) Validate() error {
	if strings.TrimSpace(tjs.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if tjs.Class != nil && *tjs.Class != "" {
		if !isValidClass(*tjs.Class) {
			return fmt.Errorf("invalid class value: %s", *tjs.Class)
		}
	}
	return nil
}

func isValidClass(class string) bool {
	for _, c := range ValidClasses {
		if c == class {
			return true
		}
	}
	return false
}
