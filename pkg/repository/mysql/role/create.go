package role

import (
	"context"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/role"
)

// Create crea un nuevo rol
func (r *Repository) Create(ctx context.Context, role *role.Role) error {
	// Verificar si ya existe un rol con el mismo nombre
	existingRole, err := r.GetByName(ctx, role.Name)
	if err == nil && existingRole != nil {
		return ErrDuplicateName
	} else if err != ErrRoleNotFound {
		return err
	}

	// Establecer timestamps
	now := time.Now()
	role.CreatedAt = &now
	role.UpdatedAt = &now

	// Preparar la consulta
	query := `
		INSERT INTO roles (name, title, scope, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// Preparar los argumentos
	var args []interface{}
	args = append(args, role.Name)

	// Manejar campos opcionales
	if role.Title != nil {
		args = append(args, *role.Title)
	} else {
		args = append(args, nil)
	}

	if role.Scope != nil {
		args = append(args, *role.Scope)
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

	role.ID = id
	return nil
}
