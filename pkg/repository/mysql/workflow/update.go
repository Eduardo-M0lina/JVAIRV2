package workflow

import (
	"context"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// Update actualiza un workflow existente
func (r *Repository) Update(ctx context.Context, workflow *domainWorkflow.Workflow) error {
	query := `
		UPDATE workflows
		SET name = ?, notes = ?, is_active = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		workflow.Name,
		workflow.Notes,
		workflow.IsActive,
		workflow.ID,
	)

	return err
}
