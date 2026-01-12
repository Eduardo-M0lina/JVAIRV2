package assigned_role

import (
	"context"
	"database/sql"

	domainAssignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
)

// GetByEntity obtiene todas las asignaciones de roles para una entidad específica
func (r *Repository) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*domainAssignedRole.AssignedRole, error) {
	query := `
		SELECT id, role_id, entity_id, entity_type, restricted_to_id, restricted_to_type, scope
		FROM assigned_roles
		WHERE entity_type = ? AND entity_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var assignedRoles []*domainAssignedRole.AssignedRole
	for rows.Next() {
		var assignedRole domainAssignedRole.AssignedRole
		var restrictedToID sql.NullInt64
		var restrictedToType sql.NullString
		var scope sql.NullInt32

		err := rows.Scan(
			&assignedRole.ID, &assignedRole.RoleID, &assignedRole.EntityID, &assignedRole.EntityType,
			&restrictedToID, &restrictedToType, &scope,
		)
		if err != nil {
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

		assignedRoles = append(assignedRoles, &assignedRole)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assignedRoles, nil
}
