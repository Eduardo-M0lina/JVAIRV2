package workflow

import (
	"context"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// Create crea un nuevo workflow
func (r *Repository) Create(ctx context.Context, workflow *domainWorkflow.Workflow) error {
	query := `
		INSERT INTO workflows (name, notes, is_active)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		workflow.Name,
		workflow.Notes,
		workflow.IsActive,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	workflow.ID = id
	return nil
}
