package permission

import (
	"context"
	"time"

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

	// Establecer timestamps
	now := time.Now()
	permission.CreatedAt = &now
	permission.UpdatedAt = &now

	// Preparar la consulta
	query := `
		INSERT INTO permissions (ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Preparar los argumentos
	var args []interface{}
	args = append(args, permission.AbilityID, permission.EntityID, permission.EntityType, permission.Forbidden)

	// Manejar campos opcionales
	if permission.Conditions != nil {
		args = append(args, *permission.Conditions)
	} else {
		args = append(args, nil)
	}

	args = append(args, now, now)

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
