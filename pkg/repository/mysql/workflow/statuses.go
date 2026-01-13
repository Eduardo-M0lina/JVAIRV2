package workflow

import (
	"context"
	"database/sql"

	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

// GetWorkflowStatuses obtiene los job_statuses asociados a un workflow
func (r *Repository) GetWorkflowStatuses(ctx context.Context, workflowID int64) ([]domainWorkflow.WorkflowStatus, error) {
	query := `
		SELECT jsw.job_status_id, jsw.workflow_id, jsw.order, js.label
		FROM job_status_workflow jsw
		INNER JOIN job_statuses js ON jsw.job_status_id = js.id
		WHERE jsw.workflow_id = ?
		ORDER BY jsw.order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, workflowID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var statuses []domainWorkflow.WorkflowStatus
	for rows.Next() {
		var status domainWorkflow.WorkflowStatus
		var statusName sql.NullString

		err := rows.Scan(
			&status.JobStatusID,
			&status.WorkflowID,
			&status.Order,
			&statusName,
		)
		if err != nil {
			return nil, err
		}

		if statusName.Valid {
			status.StatusName = statusName.String
		}

		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}

// SetWorkflowStatuses establece los job_statuses asociados a un workflow
func (r *Repository) SetWorkflowStatuses(ctx context.Context, workflowID int64, statuses []domainWorkflow.WorkflowStatus) error {
	// Iniciar transacción
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Eliminar las relaciones existentes
	_, err = tx.ExecContext(ctx, "DELETE FROM job_status_workflow WHERE workflow_id = ?", workflowID)
	if err != nil {
		return err
	}

	// Insertar las nuevas relaciones
	if len(statuses) > 0 {
		query := "INSERT INTO job_status_workflow (job_status_id, workflow_id, `order`) VALUES (?, ?, ?)"
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer func() {
			_ = stmt.Close()
		}()

		for _, status := range statuses {
			_, err = stmt.ExecContext(ctx, status.JobStatusID, workflowID, status.Order)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
