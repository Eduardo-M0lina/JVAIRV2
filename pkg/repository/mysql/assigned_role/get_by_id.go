package assigned_role

import (
	"context"
	"database/sql"

	domainAssignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

// GetByID obtiene una asignación de rol por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainAssignedRole.AssignedRole, error) {
	query := `
		SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope
		FROM assigned_roles
		WHERE id = ?
	`

	var assignedRole domainAssignedRole.AssignedRole
	var restrictedToID sql.NullInt64
	var restrictedToType sql.NullString
	var scope sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&assignedRole.ID, &assignedRole.RoleID, &assignedRole.EntityID, &assignedRole.EntityType,
		&restrictedToID, &restrictedToType, &scope,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAssignedRoleNotFound
		}
		return nil, err
	}

	if restrictedToID.Valid {
		assignedRole.RestrictedToID = &restrictedToID.Int64
		assignedRole.Restricted = true
	}

	if restrictedToType.Valid {
		assignedRole.RestrictedToType = &restrictedToType.String
	}

	if scope.Valid {
		scopeInt := int(scope.Int32)
		assignedRole.Scope = &scopeInt
	}

	return &assignedRole, nil
}
