package assigned_role

import (
	"context"
)

// Revoke revoca un rol de una entidad
func (r *Repository) Revoke(ctx context.Context, roleID, entityID int64, entityType string) error {
	// Preparar la consulta
	query := `
		DELETE FROM assigned_roles
		WHERE role_id = ? AND entity_id = ? AND entity_type = ?
	`

	// Ejecutar la consulta
	result, err := r.db.ExecContext(ctx, query, roleID, entityID, entityType)
	if err != nil {
		return err
	}

	// Verificar que se haya eliminado al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAssignedRoleNotFound
	}

	return nil
}
