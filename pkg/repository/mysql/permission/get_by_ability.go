package permission

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// GetByAbility obtiene todos los permisos para una ability específica
func (r *Repository) GetByAbility(ctx context.Context, abilityID int64) ([]*permission.Permission, error) {
	query := `
		SELECT id, ability_id, entity_id, entity_type, forbidden, scope
		FROM permissions
		WHERE ability_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, abilityID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var permissions []*permission.Permission
	for rows.Next() {
		var p permission.Permission
		var scope sql.NullInt32

		err := rows.Scan(
			&p.ID, &p.AbilityID, &p.EntityID, &p.EntityType, &p.Forbidden, &scope,
		)

		if err != nil {
			return nil, err
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			p.Scope = &scopeInt
		}

		permissions = append(permissions, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
