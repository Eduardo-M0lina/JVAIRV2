package user

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetUserRoles obtiene los roles de un usuario
func (r *Repository) GetUserRoles(ctx context.Context, userID string) ([]*user.Role, error) {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	query := `
		SELECT r.id, r.label, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN assigned_roles ar ON r.id = ar.role_id
		WHERE ar.entity_id = ? AND ar.entity_type = 'App\\Models\\User'
	`

	rows, err := r.db.QueryContext(ctx, query, idInt)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	roles := []*user.Role{}
	for rows.Next() {
		var r user.Role
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&r.ID, &r.Label, &createdAt, &updatedAt,
		)

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
