package workflow

import "time"

// Workflow representa un flujo de trabajo en el sistema
type Workflow struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Notes     *string          `json:"notes,omitempty"`
	IsActive  bool             `json:"isActive"`
	CreatedAt *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
	Statuses  []WorkflowStatus `json:"statuses,omitempty"`
}

// WorkflowStatus representa un estado de trabajo asociado a un workflow
type WorkflowStatus struct {
	JobStatusID int64  `json:"jobStatusId"`
	WorkflowID  int64  `json:"workflowId"`
	Order       int    `json:"order"`
	StatusName  string `json:"statusName,omitempty"` // Campo calculado para incluir el nombre del estado
}

// Filters representa los filtros disponibles para listar workflows
type Filters struct {
	Name     string
	IsActive *bool
	Search   string
}
