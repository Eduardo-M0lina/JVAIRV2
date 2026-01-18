package permission

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// Update actualiza un permiso existente
func (r *Repository) Update(ctx context.Context, permission *permission.Permission) error {
	// Verificar si el permiso existe
	existingPermission, err := r.GetByID(ctx, permission.ID)
	if err != nil {
		return err
	}

	// Verificar si el permiso ya existe para otra combinación de ability y entidad
	if permission.AbilityID != existingPermission.AbilityID ||
		permission.EntityID != existingPermission.EntityID ||
		permission.EntityType != existingPermission.EntityType {
		exists, err := r.Exists(ctx, permission.AbilityID, permission.EntityID, permission.EntityType)
		if err != nil {
			return err
		}
		if exists {
			return ErrDuplicatePermission
		}
	}

	// Preparar la consulta
	query := `
		UPDATE permissions
		SET ability_id = ?, entity_id = ?, entity_type = ?, forbidden = ?, scope = ?
		WHERE id = ?
	`

	// Preparar los argumentos
	var args []interface{}
	args = append(args, permission.AbilityID, permission.EntityID, permission.EntityType, permission.Forbidden)

	// Manejar campos opcionales
	if permission.Scope != nil {
		args = append(args, *permission.Scope)
	} else {
		args = append(args, nil)
	}

	args = append(args, permission.ID)

	// Ejecutar la consulta
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
