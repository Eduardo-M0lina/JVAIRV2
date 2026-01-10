package user

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/your-org/jvairv2/pkg/domain/ability"
)

// GetUserAbilities obtiene las habilidades de un usuario
func (r *Repository) GetUserAbilities(ctx context.Context, userID string) ([]*ability.Ability, error) {
	// Convertir ID de string a int64
	idInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	query := `
		SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at
		FROM abilities a
		INNER JOIN permissions p ON a.id = p.ability_id
		WHERE p.entity_id = ? AND p.entity_type = 'App\\Models\\User'
		UNION
		SELECT a.id, a.name, a.title, a.entity_id, a.entity_type, a.only_owned, a.options, a.scope, a.created_at, a.updated_at
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

	abilities := []*ability.Ability{}
	for rows.Next() {
		var a ability.Ability
		var createdAt, updatedAt sql.NullTime
		var title, entityType, options sql.NullString
		var entityID sql.NullInt64
		var scope sql.NullInt32

		err := rows.Scan(
			&a.ID, &a.Name, &title, &entityID, &entityType, &a.OnlyOwned, &options, &scope, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, err
		}

		if title.Valid {
			a.Title = &title.String
		}

		if entityID.Valid {
			a.EntityID = &entityID.Int64
		}

		if entityType.Valid {
			a.EntityType = &entityType.String
		}

		if options.Valid {
			a.Options = &options.String
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			a.Scope = &scopeInt
		}

		if createdAt.Valid {
			a.CreatedAt = &createdAt.Time
		}

		if updatedAt.Valid {
			a.UpdatedAt = &updatedAt.Time
		}

		abilities = append(abilities, &a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return abilities, nil
}
