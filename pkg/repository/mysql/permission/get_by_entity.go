package permission

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// GetByEntity obtiene todos los permisos para una entidad específica
func (r *Repository) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]*permission.Permission, error) {
	query := `
		SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at
		FROM permissions
		WHERE entity_type = ? AND entity_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var permissions []*permission.Permission
	for rows.Next() {
		var p permission.Permission
		var createdAt, updatedAt sql.NullTime
		var conditions sql.NullString

		err := rows.Scan(
			&p.ID, &p.AbilityID, &p.EntityID, &p.EntityType, &p.Forbidden, &conditions, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, err
		}

		if conditions.Valid {
			p.Conditions = &conditions.String
		}

		if createdAt.Valid {
			p.CreatedAt = &createdAt.Time
		}

		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
		}

		permissions = append(permissions, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
