package assigned_role

import (
	"context"
	"database/sql"
)

// HasRole verifica si una entidad tiene un rol específico
func (r *Repository) HasRole(ctx context.Context, roleID, entityID int64, entityType string) (bool, error) {
	// Preparar la consulta
	query := `
		SELECT 1 FROM assigned_roles
		WHERE role_id = ? AND entity_id = ? AND entity_type = ?
		LIMIT 1
	`

	// Ejecutar la consulta
	var exists int
	err := r.db.QueryRowContext(ctx, query, roleID, entityID, entityType).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
