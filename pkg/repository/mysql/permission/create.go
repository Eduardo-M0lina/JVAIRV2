package permission

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// Create crea un nuevo permiso
func (r *Repository) Create(ctx context.Context, permission *permission.Permission) error {
	// Verificar si ya existe este permiso
	exists, err := r.Exists(ctx, permission.AbilityID, permission.EntityID, permission.EntityType)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicatePermission
	}

	// Preparar la consulta
	query := `
		INSERT INTO permissions (ability_id, entity_id, entity_type, forbidden, scope)
		VALUES (?, ?, ?, ?, ?)
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

	// Ejecutar la consulta
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	permission.ID = id
	return nil
}
