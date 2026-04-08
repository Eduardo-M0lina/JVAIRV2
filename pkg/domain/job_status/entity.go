package job_status

import (
	"fmt"
	"strings"
	"time"
)

// ValidClasses contiene los valores válidos para el campo class (colores Bootstrap)
var ValidClasses = []string{"blue", "indigo", "purple", "pink", "red", "orange", "yellow", "green", "teal", "cyan", "dark", "light"}

// JobStatus representa la entidad de dominio para un estado de trabajo
type JobStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de trabajo
func (js *JobStatus) Validate() error {
	if strings.TrimSpace(js.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if js.Class != nil && *js.Class != "" {
		if !isValidClass(*js.Class) {
			return fmt.Errorf("invalid class value: %s", *js.Class)
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
