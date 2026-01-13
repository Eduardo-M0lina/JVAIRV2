package workflow

import (
	"context"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// Duplicate duplica un workflow existente
func (r *Repository) Duplicate(ctx context.Context, id int64) (*domainWorkflow.Workflow, error) {
	// Obtener el workflow original
	original, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Crear una copia del workflow
	duplicated := &domainWorkflow.Workflow{
		Name:     original.Name,
		Notes:    original.Notes,
		IsActive: original.IsActive,
	}

	// Crear el nuevo workflow
	if err := r.Create(ctx, duplicated); err != nil {
		return nil, err
	}

	return duplicated, nil
}
