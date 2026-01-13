package workflow

import (
	"context"
	"database/sql"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// GetByID obtiene un workflow por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainWorkflow.Workflow, error) {
	query := `
		SELECT id, name, notes, is_active, created_at, updated_at
		FROM workflows
		WHERE id = ?
	`

	var workflow domainWorkflow.Workflow
	var notes sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workflow.ID,
		&workflow.Name,
		&notes,
		&workflow.IsActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainWorkflow.ErrWorkflowNotFound
		}
		return nil, err
	}

	if notes.Valid {
		workflow.Notes = &notes.String
	}
	if createdAt.Valid {
		workflow.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		workflow.UpdatedAt = &updatedAt.Time
	}

	return &workflow, nil
}
