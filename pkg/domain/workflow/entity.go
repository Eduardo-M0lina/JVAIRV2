package workflow

import "time"

// Workflow representa un flujo de trabajo en el sistema
type Workflow struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Notes     *string          `json:"notes,omitempty"`
	IsActive  bool             `json:"is_active"`
	CreatedAt *time.Time       `json:"created_at,omitempty"`
	UpdatedAt *time.Time       `json:"updated_at,omitempty"`
	Statuses  []WorkflowStatus `json:"statuses,omitempty"`
}

// WorkflowStatus representa un estado de trabajo asociado a un workflow
type WorkflowStatus struct {
	JobStatusID int64  `json:"job_status_id"`
	WorkflowID  int64  `json:"workflow_id"`
	Order       int    `json:"order"`
	StatusName  string `json:"status_name,omitempty"` // Campo calculado para incluir el nombre del estado
}

// Filters representa los filtros disponibles para listar workflows
type Filters struct {
	Name     string
	IsActive *bool
	Search   string
}
