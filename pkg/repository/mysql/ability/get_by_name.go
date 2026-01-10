package ability

import (
	"context"
	"database/sql"

	domainAbility "github.com/your-org/jvairv2/pkg/domain/ability"
)

// GetByName obtiene una ability por su nombre
func (r *Repository) GetByName(ctx context.Context, name string) (*domainAbility.Ability, error) {
	query := `
		SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at
		FROM abilities
		WHERE name = ?
	`

	var ability domainAbility.Ability
	var title, entityType, options sql.NullString
	var entityID sql.NullInt64
	var scope sql.NullInt32
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&ability.ID, &ability.Name, &title, &entityID, &entityType,
		&ability.OnlyOwned, &options, &scope, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAbilityNotFound
		}
		return nil, err
	}

	if title.Valid {
		ability.Title = &title.String
	}

	if entityID.Valid {
		ability.EntityID = &entityID.Int64
	}

	if entityType.Valid {
		ability.EntityType = &entityType.String
	}

	if options.Valid {
		ability.Options = &options.String
	}

	if scope.Valid {
		scopeInt := int(scope.Int32)
		ability.Scope = &scopeInt
	}

	if createdAt.Valid {
		ability.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		ability.UpdatedAt = &updatedAt.Time
	}

	return &ability, nil
}
