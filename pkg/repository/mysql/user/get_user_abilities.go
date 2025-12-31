package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// GetUserAbilities obtiene las habilidades de un usuario
func (r *Repository) GetUserAbilities(ctx context.Context, userID string) ([]*user.Ability, error) {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	query := `
		SELECT a.id, a.name, a.title, a.description, a.created_at, a.updated_at
		FROM abilities a
		INNER JOIN permissions p ON a.id = p.ability_id
		WHERE p.entity_id = ? AND p.entity_type = 'App\\Models\\User'
		UNION
		SELECT a.id, a.name, a.title, a.description, a.created_at, a.updated_at
		FROM abilities a
		INNER JOIN permissions p ON a.id = p.ability_id
		INNER JOIN assigned_roles ar ON p.entity_id = ar.role_id
		WHERE ar.entity_id = ? AND p.entity_type = 'App\\Models\\Role' AND ar.entity_type = 'App\\Models\\User'
	`

	rows, err := r.db.QueryContext(ctx, query, idInt, idInt)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	abilities := []*user.Ability{}
	for rows.Next() {
		var a user.Ability
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&a.ID, &a.Name, &a.Title, &a.Description, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, err
		}

		a.CreatedAt = createdAt
		a.UpdatedAt = updatedAt

		abilities = append(abilities, &a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return abilities, nil
}
