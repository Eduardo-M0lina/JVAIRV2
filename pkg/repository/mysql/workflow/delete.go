package workflow

import (
	"context"
)

// Delete elimina un workflow y sus relaciones con job_statuses
func (r *Repository) Delete(ctx context.Context, id int64) error {
	// Iniciar transacción
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Eliminar las relaciones con job_statuses
	_, err = tx.ExecContext(ctx, "DELETE FROM job_status_workflow WHERE workflow_id = ?", id)
	if err != nil {
		return err
	}

	// Eliminar el workflow
	_, err = tx.ExecContext(ctx, "DELETE FROM workflows WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
