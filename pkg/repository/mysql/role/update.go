package role

import (
	"context"
	"time"

	domainRole "github.com/your-org/jvairv2/pkg/domain/role"
)

// Update actualiza un rol existente
func (r *Repository) Update(ctx context.Context, role *domainRole.Role) error {
	// Verificar si el rol existe
	existingRole, err := r.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}

	// Verificar si el nombre ya está en uso por otro rol
	if role.Name != existingRole.Name {
		otherRole, err := r.GetByName(ctx, role.Name)
		if err == nil && otherRole != nil && otherRole.ID != role.ID {
			return ErrDuplicateName
		} else if err != ErrRoleNotFound && err != nil {
			return err
		}
	}

	// Actualizar timestamp
	now := time.Now()
	role.UpdatedAt = &now

	// Preparar la consulta
	query := `
		UPDATE roles
		SET name = ?, title = ?, scope = ?, updated_at = ?
		WHERE id = ?
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

	args = append(args, now, role.ID)

	// Ejecutar la consulta
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
