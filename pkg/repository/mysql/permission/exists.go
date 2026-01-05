package permission

import (
	"context"
	"database/sql"
)

// Exists verifica si existe un permiso específico
func (r *Repository) Exists(ctx context.Context, abilityID, entityID int64, entityType string) (bool, error) {
	// Preparar la consulta
	query := `
		SELECT 1 FROM permissions
		WHERE ability_id = ? AND entity_id = ? AND entity_type = ?
		LIMIT 1
	`

	// Ejecutar la consulta
	var exists int
	err := r.db.QueryRowContext(ctx, query, abilityID, entityID, entityType).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
