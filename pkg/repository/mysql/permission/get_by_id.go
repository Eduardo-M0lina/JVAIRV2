package permission

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// GetByID obtiene un permiso por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*permission.Permission, error) {
	query := `
		SELECT id, ability_id, entity_id, entity_type, forbidden, scope
		FROM permissions
		WHERE id = ?
	`

	var p permission.Permission
	var scope sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.AbilityID, &p.EntityID, &p.EntityType, &p.Forbidden, &scope,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	if scope.Valid {
		scopeInt := int(scope.Int32)
		p.Scope = &scopeInt
	}

	return &p, nil
}
