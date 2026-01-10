package assigned_role

import (
	"context"
	"database/sql"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

// Assign asigna un rol a una entidad
func (r *Repository) Assign(ctx context.Context, assignedRole *assigned_role.AssignedRole) error {
	// Verificar si ya existe esta asignación
	query := `
		SELECT id FROM assigned_roles
		WHERE role_id = ? AND entity_id = ? AND entity_type = ?
	`
	var existingID int64
	err := r.db.QueryRowContext(ctx, query, assignedRole.RoleID, assignedRole.EntityID, assignedRole.EntityType).Scan(&existingID)
	if err == nil {
		// Ya existe esta asignación
		return ErrDuplicateAssignment
	} else if err != sql.ErrNoRows {
		// Error de base de datos
		return err
	}

	// Establecer timestamps
	now := time.Now()
	assignedRole.CreatedAt = &now
	assignedRole.UpdatedAt = &now

	// Preparar la consulta
	insertQuery := `
		INSERT INTO assigned_roles (role_id, entity_id, entity_type, restricted, scope, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Preparar los argumentos
	var args []interface{}
	args = append(args, assignedRole.RoleID, assignedRole.EntityID, assignedRole.EntityType, assignedRole.Restricted)

	// Manejar campos opcionales
	if assignedRole.Scope != nil {
		args = append(args, *assignedRole.Scope)
	} else {
		args = append(args, nil)
	}

	args = append(args, now, now)

	// Ejecutar la consulta
	result, err := r.db.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		return err
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	assignedRole.ID = id
	return nil
}
