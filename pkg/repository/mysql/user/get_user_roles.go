package user

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/your-org/jvairv2/pkg/domain/role"
)

// GetUserRoles obtiene los roles de un usuario
func (r *Repository) GetUserRoles(ctx context.Context, userID string) ([]*role.Role, error) {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	query := `
		SELECT r.id, r.name, r.title, r.scope, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN assigned_roles ar ON r.id = ar.role_id
		WHERE ar.entity_id = ? AND ar.entity_type = 'App\\Models\\User'
	`

	rows, err := r.db.QueryContext(ctx, query, idInt)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	roles := []*role.Role{}
	for rows.Next() {
		var r role.Role
		var createdAt, updatedAt sql.NullTime

		var title sql.NullString
		var scope sql.NullInt32

		err := rows.Scan(
			&r.ID, &r.Name, &title, &scope, &createdAt, &updatedAt,
		)

		if title.Valid {
			r.Title = &title.String
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			r.Scope = &scopeInt
		}

		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			r.CreatedAt = &createdAt.Time
		}

		if updatedAt.Valid {
			r.UpdatedAt = &updatedAt.Time
		}

		roles = append(roles, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
