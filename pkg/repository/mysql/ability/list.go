package ability

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	domainAbility "github.com/your-org/jvairv2/pkg/domain/ability"
)

// List obtiene una lista paginada de abilities con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainAbility.Ability, int, error) {
	// Construir la consulta base
	countQuery := "SELECT COUNT(*) FROM abilities"
	selectQuery := `
		SELECT id, name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at
		FROM abilities
	`

	// Procesar filtros
	var whereConditions []string
	var args []interface{}

	for key, value := range filters {
		switch key {
		case "name":
			whereConditions = append(whereConditions, "name LIKE ?")
			args = append(args, fmt.Sprintf("%%%s%%", value))
		case "entity_type":
			whereConditions = append(whereConditions, "entity_type = ?")
			args = append(args, value)
		case "scope":
			whereConditions = append(whereConditions, "scope = ?")
			args = append(args, value)
		}
	}

	// Agregar condiciones WHERE si existen
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	// Agregar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	selectQuery += " ORDER BY name ASC LIMIT ? OFFSET ?"

	// Ejecutar consulta de conteo
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Si no hay resultados, devolver un slice vacío
	if total == 0 {
		return []*domainAbility.Ability{}, 0, nil
	}

	// Agregar argumentos de paginación
	args = append(args, pageSize, offset)

	// Ejecutar consulta de selección
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	var abilities []*domainAbility.Ability
	for rows.Next() {
		var ability domainAbility.Ability
		var title, entityType, options sql.NullString
		var entityID sql.NullInt64
		var scope sql.NullInt32
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ability.ID, &ability.Name, &title, &entityID, &entityType,
			&ability.OnlyOwned, &options, &scope, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, 0, err
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

		abilities = append(abilities, &ability)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return abilities, total, nil
}
