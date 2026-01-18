package permission

import (
	"context"
	"database/sql"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/permission"
)

// List obtiene una lista paginada de permisos con filtros opcionales
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*permission.Permission, int, error) {
	// Construir la consulta base
	countQuery := "SELECT COUNT(*) FROM permissions"
	selectQuery := `
		SELECT id, ability_id, entity_id, entity_type, forbidden, scope
		FROM permissions
	`

	// Procesar filtros
	whereConditions := []string{}
	args := []interface{}{}

	if abilityID, ok := filters["ability_id"].(int64); ok && abilityID > 0 {
		whereConditions = append(whereConditions, "ability_id = ?")
		args = append(args, abilityID)
	}

	if entityID, ok := filters["entity_id"].(int64); ok && entityID > 0 {
		whereConditions = append(whereConditions, "entity_id = ?")
		args = append(args, entityID)
	}

	if entityType, ok := filters["entity_type"].(string); ok && entityType != "" {
		whereConditions = append(whereConditions, "entity_type = ?")
		args = append(args, entityType)
	}

	if forbidden, ok := filters["forbidden"].(bool); ok {
		whereConditions = append(whereConditions, "forbidden = ?")
		args = append(args, forbidden)
	}

	// Construir la cláusula WHERE si hay condiciones
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Aplicar la cláusula WHERE a las consultas
	countQueryWithFilters := countQuery + whereClause
	selectQueryWithFilters := selectQuery + whereClause

	// Obtener el total de registros
	var total int
	err := r.db.QueryRowContext(ctx, countQueryWithFilters, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Si no hay resultados, devolver una lista vacía
	if total == 0 {
		return []*permission.Permission{}, 0, nil
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	selectQueryWithPagination := selectQueryWithFilters + " ORDER BY id LIMIT ? OFFSET ?"
	queryArgs := append(args, pageSize, offset)

	// Ejecutar la consulta paginada
	rows, err := r.db.QueryContext(ctx, selectQueryWithPagination, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	// Procesar resultados
	permissions := []*permission.Permission{}
	for rows.Next() {
		var p permission.Permission
		var scope sql.NullInt32

		err := rows.Scan(
			&p.ID, &p.AbilityID, &p.EntityID, &p.EntityType, &p.Forbidden, &scope,
		)

		if err != nil {
			return nil, 0, err
		}

		if scope.Valid {
			scopeInt := int(scope.Int32)
			p.Scope = &scopeInt
		}

		permissions = append(permissions, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}
