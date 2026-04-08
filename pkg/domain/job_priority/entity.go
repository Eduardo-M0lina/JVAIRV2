package job_priority

import (
	"fmt"
	"strings"
	"time"
)

// ValidClasses contiene los valores válidos para el campo class (colores Bootstrap)
var ValidClasses = []string{"blue", "indigo", "purple", "pink", "red", "orange", "yellow", "green", "teal", "cyan", "dark", "light"}

// JobPriority representa la entidad de dominio para una prioridad de trabajo
type JobPriority struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Order     int        `json:"order"`
	Class     *string    `json:"class,omitempty"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos de la prioridad de trabajo
func (jp *JobPriority) Validate() error {
	if strings.TrimSpace(jp.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if jp.Class != nil && *jp.Class != "" {
		if !isValidClass(*jp.Class) {
			return fmt.Errorf("invalid class value: %s", *jp.Class)
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
