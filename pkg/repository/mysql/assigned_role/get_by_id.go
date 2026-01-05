package assigned_role

import (
	"context"
	"database/sql"

	domainAssignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

// GetByID obtiene una asignación de rol por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainAssignedRole.AssignedRole, error) {
	query := `
		SELECT id, role_id, entity_id, entity_type, restricted, scope, created_at, updated_at
		FROM assigned_roles
		WHERE id = ?
	`

	var assignedRole domainAssignedRole.AssignedRole
	var scope sql.NullInt32
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&assignedRole.ID, &assignedRole.RoleID, &assignedRole.EntityID, &assignedRole.EntityType,
		&assignedRole.Restricted, &scope, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignedRoleNotFound
		}
		return nil, err
	}

	if scope.Valid {
		scopeInt := int(scope.Int32)
		assignedRole.Scope = &scopeInt
	}

	if createdAt.Valid {
		assignedRole.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		assignedRole.UpdatedAt = &updatedAt.Time
	}

	return &assignedRole, nil
}
