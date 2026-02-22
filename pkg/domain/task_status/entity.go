package task_status

import (
	"fmt"
	"strings"
	"time"
)

// ValidClasses contiene los valores válidos para el campo class (colores Bootstrap)
var ValidClasses = []string{"blue", "indigo", "purple", "pink", "red", "orange", "yellow", "green", "teal", "cyan", "dark", "light"}

// TaskStatus representa la entidad de dominio para un estado de tarea
type TaskStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	Order     int        `json:"order"`
	IsActive  bool       `json:"isActive"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Validate valida los campos requeridos del estado de tarea
func (ts *TaskStatus) Validate() error {
	if strings.TrimSpace(ts.Label) == "" {
		return fmt.Errorf("label is required")
	}
	if ts.Class != nil && *ts.Class != "" {
		if !isValidClass(*ts.Class) {
			return fmt.Errorf("invalid class value: %s", *ts.Class)
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
