package workflow

import "context"

// Repository define los métodos para interactuar con la base de datos de workflows
type Repository interface {
	List(ctx context.Context, filters Filters, page, pageSize int) ([]Workflow, int64, error)
	GetByID(ctx context.Context, id int64) (*Workflow, error)
	Create(ctx context.Context, workflow *Workflow) error
	Update(ctx context.Context, workflow *Workflow) error
	Delete(ctx context.Context, id int64) error
	Duplicate(ctx context.Context, id int64) (*Workflow, error)

	// Métodos para gestionar la relación con job_statuses
	GetWorkflowStatuses(ctx context.Context, workflowID int64) ([]WorkflowStatus, error)
	SetWorkflowStatuses(ctx context.Context, workflowID int64, statuses []WorkflowStatus) error
}
