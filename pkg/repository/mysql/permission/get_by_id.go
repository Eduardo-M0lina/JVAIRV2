package permission

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// GetByID obtiene un permiso por su ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*permission.Permission, error) {
	query := `
		SELECT id, ability_id, entity_id, entity_type, forbidden, conditions, created_at, updated_at
		FROM permissions
		WHERE id = ?
	`

	var p permission.Permission
	var conditions sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.AbilityID, &p.EntityID, &p.EntityType, &p.Forbidden, &conditions, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPermissionNotFound
		}
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

	return &p, nil
}
